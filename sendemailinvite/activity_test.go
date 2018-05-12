package sendemailinvite

import (
	"io/ioutil"
	"testing"

	"github.com/TIBCOSoftware/flogo-contrib/action/flow/test"
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
)

var activityMetadata *activity.Metadata

func getActivityMetadata() *activity.Metadata {

	if activityMetadata == nil {
		jsonMetadataBytes, err := ioutil.ReadFile("activity.json")
		if err != nil {
			panic("No Json Metadata found for activity.json path")
		}

		activityMetadata = activity.NewMetadata(string(jsonMetadataBytes))
	}

	return activityMetadata
}

func TestCreate(t *testing.T) {

	act := NewActivity(getActivityMetadata())

	if act == nil {
		t.Error("Activity Not Created")
		t.Fail()
		return
	}
}

func TestEval(t *testing.T) {

	defer func() {
		if r := recover(); r != nil {
			t.Failed()
			t.Errorf("panic during execution: %v", r)
		}
	}()

	act := NewActivity(getActivityMetadata())
	tc := test.NewTestActivityContext(getActivityMetadata())

	//setup attrs testcase#1
	tc.SetInput("authToken", "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ")
	tc.SetInput("surveyName", "FLG_2_QA_Variety")
	tc.SetInput("type", "invite")
	tc.SetInput("senderEmail", "sankpal22pankaj@gmail.com")
	tc.SetInput("body", "")
	tc.SetInput("subject", "Gathering inputs")
	tc.SetInput("recipientList", "psankpal@tibco.com")
	tc.SetInput("recipientStatus", "responded")
	act.Eval(tc)

	// testcase #2
	tc.SetInput("authToken", "UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ")
	tc.SetInput("surveyName", "FLG_2_QA_Variety")
	tc.SetInput("type", "reminder")
	tc.SetInput("senderEmail", "sankpal22pankaj@gmail.com")
	tc.SetInput("body", "")
	tc.SetInput("subject", "Gathering inputs")
	tc.SetInput("recipientList", "psankpal@gmail.com")
	tc.SetInput("recipientStatus", "has_not_responded")
	act.Eval(tc)

	//check result attr
}
