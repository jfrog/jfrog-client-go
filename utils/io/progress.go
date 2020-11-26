package io

// You may implement this interface to display progress indication of files transfer (upload / download)
type ProgressMgr interface {
	// Initializes a new reader progress indication for a new file transfer.
	// Input: 'total' - file size, 'prefix' - optional description, 'filePath' - path of the file being transferred (for description purposes only).
	// Output: progress indication id
	NewReaderProgressBar(total int64, prefix, filePath string) (bar ProgressBar)
	// Replaces an indication (with the 'replaceId') when completed. Used when an additional work is done as part of the transfer.
	AddNewReplacementSpinner(replaceId int, prefix, extraInformation string) (id int)
	// Returns the requested progress bar
	GetProgressBar(id int) (bar ProgressBar)
	// Aborts a progress indication. Called on both successful and unsuccessful operations
	Abort(id int)
	// Quits the whole progress mechanism
	Quit()
	// Increses the general progress bar total count by given n.
	IncreaseGeneralProgressTotalBy(n int64)
}
