package busslogic

import (
	"bytes"
	"errors"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//var declaration
var surveyID string

const (
	methodGet  = "GET"
	methodPost = "POST"
)

// func main() {
// 	_, issue := SendEmail("z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ", "DemoServey", "sankpal22pankaj@gmail.com", "", "reminder", "has_not_responded", "", "")
// 	if issue != nil {
// 		fmt.Printf(issue.Error())
// 	}
// }

func callURL(method string, url string, bodyContent *bytes.Buffer, accessToken string) (succ string, err error) {
	request, _ := http.NewRequest(method, url, bodyContent)
	request.Header.Set("Authorization", "bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{
		Timeout: time.Second * 30,
	}
	succResp, errorResp := client.Do(request)
	if errorResp != nil {
		return "", errorResp
	}
	defer succResp.Body.Close()
	surveyResponse, _ := ioutil.ReadAll(succResp.Body)
	hasError := gjson.Get(string(surveyResponse), "error.http_status_code").String()
	if hasError != "" {
		outResult := gjson.Get(string(surveyResponse), "error.message").String()
		return "", errors.New(outResult)
	}
	return string(surveyResponse), nil
}

// SendEmail func will make a set of API Calls to send email
func SendEmail(accessToken string, surveyName string, senderEmail string, recipientList string, typeofEmail string, recipientStatus string, subject string, body string) (bool, error) {

	var isInvite = false
	var surveyIdurl = "https://api.surveymonkey.com/v3/surveys?title="
	var collectorURL = "https://api.surveymonkey.com/v3/surveys/"
	var messageURL = "https://api.surveymonkey.com/v3/collectors/"
	var recipientURL = "https://api.surveymonkey.com/v3/collectors/"
	var sendURL = "https://api.surveymonkey.com/v3/collectors/"

	if typeofEmail == "invite" {
		isInvite = true
	}

	//get surveyID , API call #1
	surveyIdurl = surveyIdurl + surveyName
	var jsonBody = []byte("")
	reqSurveyID, err := callURL(methodGet, surveyIdurl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		return false, err
	}
	surveyID = gjson.Get(reqSurveyID, "data.0.id").String()

	//set email invite and get Collector id , API Call #2
	collectorID := ""
	reqCollectorID := ""
	collectorURL = collectorURL + surveyID + "/collectors"
	jsonBody = []byte("")
	reqCollectorID, err = callURL(methodGet, collectorURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
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
		jsonBody = []byte(`{"type":"email","sender_email":"` + senderEmail + `"}`)
		reqCollectorID, err = callURL(methodPost, collectorURL, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			return false, err
		}
		collectorID = gjson.Get(reqCollectorID, "id").String()
	}
	//get message ID , APICall #3
	messageID := ""
	messageURL = messageURL + collectorID + "/messages"
	if body != "" {
		surveyLink := "[SurveyLink]"
		optLink := "[OptOutLink]"
		footerLink := "[FooterLink]"
		mandatoryContent := "\\n" + "Take survey: " + surveyLink + "\\n" + "Remove me from Mailing List: " + optLink + "\\n" + "Footer: " + footerLink
		if isInvite {
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `","body_text":"` + body + mandatoryContent + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","recipient_status":"` + recipientStatus + `","body_text":"` + body + mandatoryContent + `"}`)
		}
	} else {
		if isInvite {
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","recipient_status":"` + recipientStatus + `"}`)
		}
	}
	reqMessageID, err := callURL(methodPost, messageURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		return false, err
	}
	messageID = gjson.Get(reqMessageID, "id").String()

	//add multiple email ids , API Call #4
	var emailParentJSON string
	recipientURL = recipientURL + collectorID + "/messages/" + messageID + "/recipients/bulk"
	emails := strings.Split(recipientList, ",")
	emailParentJSON = `{ "contacts": [`
	count := 0
	for i := 0; i < len(emails); i++ {
		if emails[i] == "" {
			count++
			continue
		}

		innerJSONContent, _ := sjson.Set("", "email", emails[i])
		emailParentJSON = emailParentJSON + innerJSONContent
		if count < len(emails)-1 {
			count = count + 1
			emailParentJSON = emailParentJSON + ","
		}
	}
	emailParentJSON = emailParentJSON + "]}"
	jsonBody = []byte(emailParentJSON)
	reqRecipientBulk, err := callURL(methodPost, recipientURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		return false, err
	}
	succStatus1 := gjson.Get(reqRecipientBulk, "succeeded.#").String()
	succStatus2 := gjson.Get(reqRecipientBulk, "existing.#").String()
	if !(succStatus1 != "" || succStatus2 != "") {
		return false, errors.New("emails not added")
	}

	//add schedule email date , API Call #5
	sendURL = sendURL + collectorID + "/messages/" + messageID + "/send"
	currDate := time.Now().Format("2006-01-02T15:04:05+00:00")
	jsonBody = []byte(`{ "scheduled_date": "` + currDate + `"}`)
	_, err = callURL(methodPost, sendURL, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		return false, err
	}
	return true, nil
}
