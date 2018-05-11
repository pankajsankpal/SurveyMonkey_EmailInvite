package sendemailinvite

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/pankajsankpal/SurveyMonkey_EmailInvite/busslogic"
)

// MyActivity is a stub for your Activity implementation
type MyActivity struct {
	metadata *activity.Metadata
}

// NewActivity creates a new activity
func NewActivity(metadata *activity.Metadata) activity.Activity {
	return &MyActivity{metadata: metadata}
}

// Metadata implements activity.Activity.Metadata
func (a *MyActivity) Metadata() *activity.Metadata {
	return a.metadata
}

// Eval implements activity.Activity.Eval
func (a *MyActivity) Eval(context activity.Context) (done bool, err error) {
	// Initialize variable
	accessToken := context.GetInput("AuthToken").(string)
	surveyName := context.GetInput("Survey Name").(string)
	senderEmail := context.GetInput("Sender's Email").(string)
	recipientList := context.GetInput("recipientList").(string)
	typeofEmail := context.GetInput("type").(string)
	recipientStatus := context.GetInput("recipient_status").(string)
	subject := context.GetInput("Subject").(string)
	body := context.GetInput("Body").(string)

	//containError := ""

	status, errResp := busslogic.SendEmail(accessToken, surveyName, senderEmail, recipientList, typeofEmail, recipientStatus, subject, body)
	if errResp != nil && status {
		return false, errResp
	}
	return true, nil

}
