package content

import (
	"encoding/json"
	"io"
	"os"
	"path/filepath"
	"sync"
	"testing"
)

// TestPeekResetRaceReproduction reproduces the production panic:
//
//	panic: send on closed channel
//	  contentreader.go:216 (cr.dataChannel <- ResultItem)
//	  contentreader.go:178 (run -> readSingleFile)
//	  contentreader.go:78  (NextRecord.func1.1 -> run)
//	  contentreader.go:75  (NextRecord.func1   -> go func)
//
// The pattern that triggers it (used in
// artifactory/services/utils/searchutil.go:filterBuildArtifactsAndDependencies):
//
//  1. NextRecord(item)        // peeks one record; spawns producer G1 via sync.Once
//  2. reader.Reset()           // swaps cr.dataChannel and cr.once *without* waiting for G1
//  3. NextRecord(item) ...     // spawns producer G2 (new sync.Once); G1 and G2 now share C2,
//                              // and whichever finishes first closes C2 while the other sends.
//
// Run with: go test -race -run TestPeekResetRaceReproduction ./utils/io/content/
func TestPeekResetRaceReproduction(t *testing.T) {
	tmpDir := t.TempDir()
	filePath := filepath.Join(tmpDir, "records.json")

	// Write enough records that the producer goroutine is still actively
	// pumping when the test calls Reset()/NextRecord() again. We stay under
	// MaxBufferSize so the producer never blocks (blocking the producer would
	// just leak the goroutine instead of triggering the panic).
	const records = 20000
	if err := writeRecordsFile(filePath, records); err != nil {
		t.Fatalf("setup: %v", err)
	}

	// Trigger the race repeatedly. The exact ordering of G1/G2 finishing is
	// timing-dependent, so one attempt isn't always enough on a fast machine;
	// looping makes the panic deterministic in practice.
	for attempt := 0; attempt < 50; attempt++ {
		reader := NewContentReader(filePath, DefaultKey)

		// Step 1: peek (mimics filterBuildArtifactsAndDependencies).
		var peeked map[string]interface{}
		if err := reader.NextRecord(&peeked); err != nil {
			t.Fatalf("attempt %d: peek: %v", attempt, err)
		}

		// Step 2: Reset() while producer G1 is almost certainly still running.
		reader.Reset()

		// Step 3: drive a second consumer pass; this spawns G2 because Reset()
		// installed a fresh sync.Once. Now G1 and G2 race over the same channel.
		var wg sync.WaitGroup
		wg.Add(1)
		go func() {
			defer wg.Done()
			for item := new(map[string]interface{}); ; item = new(map[string]interface{}) {
				if err := reader.NextRecord(item); err != nil {
					if err == io.EOF {
						return
					}
					return
				}
			}
		}()
		wg.Wait()
	}
}

func writeRecordsFile(path string, n int) error {
	f, err := os.Create(path)
	if err != nil {
		return err
	}
	defer f.Close()

	if _, err := f.WriteString(`{"results":[`); err != nil {
		return err
	}
	enc := json.NewEncoder(f)
	for i := 0; i < n; i++ {
		if i > 0 {
			if _, err := f.WriteString(","); err != nil {
				return err
			}
		}
		if err := enc.Encode(map[string]any{
			"repo":        "some-virtual-repo",
			"path":        "org/example",
			"name":        "artifact.jar",
			"actual_sha1": "0123456789abcdef0123456789abcdef01234567",
			"size":        1234,
			"index":       i,
		}); err != nil {
			return err
		}
	}
	if _, err := f.WriteString(`]}`); err != nil {
		return err
	}
	return nil
}
