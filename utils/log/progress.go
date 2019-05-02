package log

import "io"

type ProgressBar interface {
	NewBar(total int64, prefix, description string) (barId int)
	NewBarReplacement(replaceBarId int, prefix, description string) (barId int)
	ReadWithProgress(barId int, reader io.Reader) io.Reader
	Abort(barId int)
	Quit()
}
