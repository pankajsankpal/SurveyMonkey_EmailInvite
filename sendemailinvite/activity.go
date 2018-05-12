package sendemailinvite

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pankajsankpal/SurveyMonkey_EmailInvite/busslogic"
)

//logger
var log = logger.GetLogger("sendemailinvite_activity.go logger")

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
	accessToken := context.GetInput("authToken").(string)
	surveyName := context.GetInput("surveyName").(string)
	senderEmail := context.GetInput("senderEmail").(string)
	recipientList := context.GetInput("recipientList").(string)
	typeofEmail := context.GetInput("type").(string)
	recipientStatus := context.GetInput("recipientStatus").(string)
	subject := context.GetInput("subject").(string)
	body := context.GetInput("body").(string)

	//containError := ""
	status, errResp := busslogic.SendEmail(accessToken, surveyName, senderEmail, recipientList, typeofEmail, recipientStatus, subject, body)
	if errResp != nil && status {
		log.Errorf("error: [%s]", err.Error())
		return false, errResp
	}
	return true, nil
}
