package sendEmailInvite

import (
	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/pankajsankpal/SurveyMonkey_EmailInvite/busslogic"
)

//logger
var log = logger.GetLogger("activity-go logger")

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

	_, errResp := busslogic.SendEmail(accessToken, surveyName, senderEmail, recipientList, typeofEmail, recipientStatus, subject, body)
	if errResp != nil {
		log.Debugf("error occured " + errResp.Error())
		context.SetOutput("status", errResp.Error())
		return false, errResp
	} else {
		context.SetOutput("status", "Email sent successfully")
	}
	return true, nil

}
