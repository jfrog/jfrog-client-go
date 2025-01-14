package io

import "io"

// You may implement this interface to display progress indication of files transfer (upload / download)
type ProgressMgr interface {
	// Initializes a new reader progress indicator for a new file transfer.
	// Input: 'total' - file size, 'label' - the title of the operation, 'path' - the path of the file being processed.
	// Output: progress indicator id
	NewProgressReader(total int64, label, path string) (progress Progress)
	// Changes progress indicator state to merging.
	SetMergingState(id int, useSpinner bool) (bar Progress)
	// Returns the requested progress indicator.
	GetProgress(id int) (progress Progress)
	// Aborts a progress indicator. Called on both successful and unsuccessful operations.
	RemoveProgress(id int)
	// Increments the general progress by 1
	IncrementGeneralProgress()
	// Quits the whole progress mechanism.
	Quit() (err error)
	// Increments the general progress total count by given n.
	IncGeneralProgressTotalBy(n int64)
	// Replace the headline progress indicator message with new one.
	SetHeadlineMsg(msg string)
	// Terminate the headline progress indicator.
	ClearHeadlineMsg()
	// Specific initialization of reader progress indicators.
	// Should be called before the first call to NewProgressReader.
	InitProgressReaders()
	// Clear tasks progress bar
	ClearProgress()
}

type Progress interface {
	// Used for updating the progress indicator progress.
	ActionWithProgress(reader io.Reader) (results io.Reader)
	// Used for setting the progress indicator progress.
	SetProgress(progress int64)
	// Aborts a progress indicator. Called on both successful and unsuccessful operations
	Abort()
	// Returns the Progress ID
	GetId() (Id int)
}
