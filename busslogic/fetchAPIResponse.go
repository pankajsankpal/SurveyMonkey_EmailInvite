package busslogic

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//logger and other var initialization
var log = logger.GetLogger("fetchAPIResponse logger")
var isInvite bool
var surveyID string
var method string

// func main() {
// 	SendEmail("z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ", "DemoServey", "sankpal22pankaj@gmail.com", "psankpal@tibco.com", "reminder", "has_not_responded", "TestInvite", "body string")
// }

func callURL(method string, url string, bodyContent *bytes.Buffer, accessToken string) (succ string, err error) {
	request, _ := http.NewRequest(method, url, bodyContent)
	request.Header.Set("Authorization", "bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	succResp, errorResp := client.Do(request)
	if errorResp != nil {
		return "", errorResp
	}
	surveyResponse, _ := ioutil.ReadAll(succResp.Body)
	hasError := gjson.Get(string(surveyResponse), "error.http_status_code").String()
	if hasError != "" {
		outResult := `{ "Error" : { "message" : ` + gjson.Get(string(surveyResponse), "error.message").String() + ` } }`
		return "", errors.New(outResult)
	}
	return string(surveyResponse), nil
}

// SendEmail func will make a set of API Calls to send email
func SendEmail(accessToken string, surveyName string, senderEmail string, recipientList string, typeofEmail string, recipientStatus string, subject string, body string) (resp bool, err error) {

	var surveyIdurl = "https://api.surveymonkey.com/v3/surveys?title="
	var collectorURL = "https://api.surveymonkey.com/v3/surveys/"
	var messageURL = "https://api.surveymonkey.com/v3/collectors/"
	var recipientURL = "https://api.surveymonkey.com/v3/collectors/"
	var sendURL = "https://api.surveymonkey.com/v3/collectors/"

	fmt.Printf("typeofEmail:" + typeofEmail)
	if typeofEmail == "invite" {
		isInvite = true
	}

	//get surveyID , API call #1
	method = "GET"
	surveyIdurl = surveyIdurl + surveyName
	var jsonBody = []byte("")
	reqSurveyID, err := callURL(method, surveyIdurl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching surveyID [%s]", err.Error())
		return false, err
	}
	surveyID = gjson.Get(reqSurveyID, "data.0.id").String()
	log.Infof("surveyId: [%s]", surveyID)

	//set email invite and get Collector id , API Call #2
	collectorID := ""
	reqCollectorID := ""
	collectorURL = collectorURL + surveyID + "/collectors"
	method = "GET"
	jsonBody = []byte("")
	reqCollectorID, err = callURL(method, collectorURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching collectorID: [%s]", err.Error())
		return false, err
	}
	collectorList := gjson.Get(reqCollectorID, "data")
	for _, item := range collectorList.Array() {
		if strings.Contains(gjson.Get(item.String(), "name").String(), "Email") {
			collectorID = gjson.Get(item.String(), "id").String()
			break
		}
	}
	if collectorID == "" {
		method = "POST"
		jsonBody = []byte(`{"type":"email","sender_email":"` + senderEmail + `"}`)
		reqCollectorID, err = callURL(method, collectorURL, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Errorf("error while fetching collectorID [%s] ", err.Error())
			return false, err
		}
		collectorID = gjson.Get(reqCollectorID, "id").String()
	}

	//get message ID , APICall #3
	method = "POST"
	messageID := ""
	messageURL = messageURL + collectorID + "/messages"
	if body != "" {
		surveyLink := "[SurveyLink]"
		optLink := "[OptOutLink]"
		footerLink := "[FooterLink]"
		mandatoryContent := "\\n" + "Take survey: " + surveyLink + "\\n" + "Remove me from Mailing List: " + optLink + "\\n" + "Footer: " + footerLink
		if isInvite {
			fmt.Printf("inside invite")
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `","body_text":"` + body + mandatoryContent + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","recipient_status":"` + recipientStatus + `","body_text":"` + body + mandatoryContent + `"}`)
		}
	} else {
		if isInvite {
			fmt.Printf("inside invite")
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","subject":"` + subject + `","recipient_status":"` + recipientStatus + `"}`)
		}
	}
	reqMessageID, err := callURL(method, messageURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching messageID: [%s]", err.Error())
		return false, err
	}
	messageID = gjson.Get(reqMessageID, "id").String()

	//add multiple email ids , API Call #4
	if isInvite {
		method = "POST"
		recipientURL = recipientURL + collectorID + "/messages/" + messageID + "/recipients/bulk"
		emails := strings.Split(recipientList, ",")
		emailParentJSON := `{ "contacts": [`
		count := 0
		for i := 0; i < len(emails); i++ {
			innerJSONContent, _ := sjson.Set("", "email", emails[i])
			emailParentJSON = emailParentJSON + innerJSONContent
			if count < len(emails)-1 {
				count = count + 1
				emailParentJSON = emailParentJSON + ","
			}
		}
		emailParentJSON = emailParentJSON + "]}"
		jsonBody = []byte(emailParentJSON)
		reqRecipientBulk, err := callURL(method, recipientURL, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Errorf("error while attatching a message body :[%s]", err.Error())
			return false, err
		}
		succStatus1 := gjson.Get(reqRecipientBulk, "succeeded.#").String()
		succStatus2 := gjson.Get(reqRecipientBulk, "existing.#").String()
		if succStatus1 != "" || succStatus2 != "" {
			log.Infof("emails added successfully")
		} else {
			return false, errors.New("error while attaching bulk emails")
		}
	}

	//add schedule email date , API Call #5
	method = "POST"
	sendURL = sendURL + collectorID + "/messages/" + messageID + "/send"
	currDate := time.Now().Format("2006-01-02T15:04:05+00:00")
	jsonBody = []byte(`{ "scheduled_date": "` + currDate + `"}`)
	_, err = callURL(method, sendURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while sending emails: [%s]", err.Error())
		return false, err
	}
	return true, nil
}
