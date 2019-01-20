package spinner

import (
	"time"

	"github.com/briandowns/spinner"
)

var spinnerCharsets = spinner.CharSets[14]
var spinnerDuration = 100 * time.Millisecond
var globalSpinner = spinner.New(spinnerCharsets, spinnerDuration)
var defaultSuffix = ""
var completeMessage = ""

func ChangeStatus(status string) {
	globalSpinner.Suffix = "  " + status + defaultSuffix
}

func SetDefaultSuffix(suffix string) {
	defaultSuffix = suffix
}

func SetCompleteMessage(msg string) {
	completeMessage = msg
}

func DisplayStatus(status string) {
	ChangeStatus(status)
	globalSpinner.Start()
}

func CompleteStatus() {
	globalSpinner.FinalMSG = completeMessage
	globalSpinner.Stop()
}
