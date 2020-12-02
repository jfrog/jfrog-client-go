package io

import "io"

// You may implement this interface to display progress indication of files transfer (upload / download)
type ProgressMgr interface {
	// Initializes a new reader progress indicator for a new file transfer.
	// Input: 'total' - file size, 'prefix' - optional description, 'extraInformation' -extra information for disply.
	// Output: progress indicator id
	NewProgressReader(total int64, prefix, extraInformation string) (progress Progress)
	// Changes progress indicator state
	SetProgressState(id int, state string)
	// Returns the requested progress indicator.
	GetProgress(id int) (progress Progress)
	// Aborts a progress indicator. Called on both successful and unsuccessful operations
	RemoveProgress(id int)
	// Quits the whole progress mechanism
	Quit()
	// Increments the general progress total count by given n.
	IncGeneralProgressTotalBy(n int64)
}

type Progress interface {
	// Used for updating the progress indicator progress.
	ActionWithProgress(reader io.Reader) (results io.Reader)
	// Aborts a progress indicator. Called on both successful and unsuccessful operations
	Abort()
	// Returns the Progress ID
	GetId() (Id int)
}
