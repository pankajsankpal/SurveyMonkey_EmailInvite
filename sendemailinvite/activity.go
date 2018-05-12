package sendemailinvite

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pankajsankpal/SurveyMonkey_EmailInvite/busslogic"
)

//logger
var log = logger.GetLogger("sendemailinvite_activity.go logger")

const (
	ivAccessToken     = "authToken"
	ivSurveyName      = "surveyName"
	ivSenderEmail     = "senderEmail"
	ivRecipientList   = "recipientList"
	ivTypeofEmail     = "type"
	ivRecipientStatus = "recipientStatus"
	ivSubject         = "subject"
	ivBody            = "body"
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
	accessToken := context.GetInput(ivAccessToken).(string)
	surveyName := context.GetInput(ivSurveyName).(string)
	senderEmail := context.GetInput(ivSenderEmail).(string)
	recipientList := context.GetInput(ivRecipientList).(string)
	typeofEmail := context.GetInput(ivTypeofEmail).(string)
	recipientStatus := context.GetInput(ivRecipientStatus).(string)
	subject := context.GetInput(ivSubject).(string)
	body := context.GetInput(ivBody).(string)

	//containError := ""
	_, errResp := busslogic.SendEmail(accessToken, surveyName, senderEmail, recipientList, typeofEmail, recipientStatus, subject, body)
	if errResp != nil {
		log.Errorf("error: [%s]", errResp.Error())
		return false, errResp
	} else {
		log.Infof("Email sent successfully")
	}
	//fmt.Printf(status)
	return true, nil
}
