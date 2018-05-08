package sendEmailInvite

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
	tc.SetInput("AuthToken", "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ")
	tc.SetInput("Survey Name", "FLG_2_QA_Variety")
	tc.SetInput("type", "invite")
	tc.SetInput("Sender's Email", "sankpal22pankaj@gmail.com")
	tc.SetInput("Body", "")
	tc.SetInput("Subject", "Gathering inputs")
	tc.SetInput("recipientList", "psankpal@tibco.com")
	tc.SetInput("recipient_status", "responded")
	act.Eval(tc)

	// testcase #2
	tc.SetInput("AuthToken", "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ")
	tc.SetInput("Survey Name", "FLG_2_QA_Variety")
	tc.SetInput("type", "reminder")
	tc.SetInput("Sender's Email", "sankpal22pankaj@gmail.com")
	tc.SetInput("Body", "")
	tc.SetInput("Subject", "Gathering inputs")
	tc.SetInput("recipientList", "")
	tc.SetInput("recipient_status", "has_not_responded")
	act.Eval(tc)
	//testcase #3
	tc.SetInput("AuthToken", "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ")
	tc.SetInput("Survey Name", "FLG_2_QA_Variety")
	tc.SetInput("type", "thank_you")
	tc.SetInput("Sender's Email", "sankpal22pankaj@gmail.com")
	tc.SetInput("Body", "we need this survey response as soon as possible..")
	tc.SetInput("Subject", "")
	tc.SetInput("recipientList", "")
	tc.SetInput("recipient_status", "responded")
	act.Eval(tc)
	//check result attr
}
