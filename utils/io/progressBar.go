package io

type ProgressBar interface {
	// Used to updated the progress bar progress.
	ActionWithProgress(...interface{}) (results interface{})
	// Aborts a progress indication. Called on both successful and unsuccessful operations
	Abort()
	// Returns the ProgressBar ID
	GetId() (Id int)
}
