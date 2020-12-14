package services

import (
	"testing"
	"time"

	"github.com/jfrog/jfrog-client-go/artifactory/buildinfo"
)

func TestCalculateMinimumBuildDate(t *testing.T) {

	time1, _ := time.Parse(buildinfo.TimeFormat, "2018-05-07T17:34:49.729+0300")
	time2, _ := time.Parse(buildinfo.TimeFormat, "2018-05-07T17:34:49.729+0300")
	time3, _ := time.Parse(buildinfo.TimeFormat, "2018-05-07T17:34:49.729+0300")

	tests := []struct {
		testName      string
		startingDate  time.Time
		maxDaysString string
		expectedTime  string
	}{
		{"test_max_days=3", time1, "3", "2018-05-04T17:34:49.729+0300"},
		{"test_max_days=0", time2, "0", "2018-05-07T17:34:49.729+0300"},
		{"test_max_days=-1", time3, "-3", "2018-05-10T17:34:49.729+0300"},
	}

	for _, test := range tests {
		t.Run(test.testName, func(t *testing.T) {
			actual, _ := calculateMinimumBuildDate(test.startingDate, test.maxDaysString)
			if test.expectedTime != actual {
				t.Errorf("Test name: %s: Expected: %s, Got: %s", test.testName, test.expectedTime, actual)
			}
		})
	}
}
