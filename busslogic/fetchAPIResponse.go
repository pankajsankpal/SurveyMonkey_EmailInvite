package busslogic

import (
	"bytes"
	"errors"
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
// 	SendEmail("z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ", "DemoServey", "sankpal22pankaj@gmail.com", "psankpal@tibco.com", "invite", "", "TestInvite", "body string")
// }

func callUrl(method string, url string, bodyContent *bytes.Buffer, accessToken string) (succ string, err error) {
	request, _ := http.NewRequest(method, url, bodyContent)
	request.Header.Set("Authorization", "bearer "+accessToken)
	request.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	succ_resp, error_resp := client.Do(request)
	if error_resp != nil {
		return "", error_resp
	} else {
		survey_response, _ := ioutil.ReadAll(succ_resp.Body)
		hasError := gjson.Get(string(survey_response), "error.http_status_code").String()
		if hasError != "" {
			outResult := `{ "Error" : { "message" : ` + gjson.Get(string(survey_response), "error.message").String() + ` } }`
			return "", errors.New(outResult)
		} else {
			return string(survey_response), nil
		}
	}
}

func SendEmail(accessToken string, surveyName string, senderEmail string, recipientList string, typeofEmail string, recipientStatus string, subject string, body string) (resp bool, err error) {

	var surveyIdurl string = "https://api.surveymonkey.com/v3/surveys?title="
	var collectorUrl string = "https://api.surveymonkey.com/v3/surveys/"
	var messageUrl string = "https://api.surveymonkey.com/v3/collectors/"
	var recipientUrl string = "https://api.surveymonkey.com/v3/collectors/"
	var sendUrl string = "https://api.surveymonkey.com/v3/collectors/"

	if typeofEmail == "invite" {
		isInvite = true
	}

	//get surveyID , API call #1
	method = "GET"
	surveyIdurl = surveyIdurl + surveyName
	var jsonBody = []byte("")
	reqSurveyID, err := callUrl(method, surveyIdurl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching surveyID [%s]", err.Error())
		return false, err
	} else {
		surveyID = gjson.Get(reqSurveyID, "data.0.id").String()
		log.Infof("surveyId: [%s]", surveyID)
	}

	//set email invite and get Collector id , API Call #2
	collectorID := ""
	reqCollectorID := ""
	collectorUrl = collectorUrl + surveyID + "/collectors"
	method = "GET"
	jsonBody = []byte("")
	reqCollectorID, err = callUrl(method, collectorUrl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching collectorID: [%s]", err.Error())
		return false, err
	} else {
		collectorList := gjson.Get(reqCollectorID, "data")
		for _, item := range collectorList.Array() {
			if strings.Contains(gjson.Get(item.String(), "name").String(), "Email") {
				collectorID = gjson.Get(item.String(), "id").String()
				break
			}
		}
	}
	if collectorID == "" {
		method = "POST"
		jsonBody = []byte(`{"type":"email","sender_email":"` + senderEmail + `"}`)
		reqCollectorID, err = callUrl(method, collectorUrl, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Errorf("error while fetching collectorID [%s] ", err.Error())
			return false, err
		} else {
			collectorID = gjson.Get(reqCollectorID, "id").String()
		}
	}

	//get message ID , APICall #3
	method = "POST"
	messageID := ""
	messageUrl = messageUrl + collectorID + "/messages"
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
			jsonBody = []byte(`{"type":"` + typeofEmail + `","subject":"` + subject + `","recipient_status":"` + recipientStatus + `"}`)
		}
	}
	reqMessageID, err := callUrl(method, messageUrl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while fetching messageID: [%s]", err.Error())
		return false, err
	} else {
		messageID = gjson.Get(reqMessageID, "id").String()
	}

	//add multiple email ids , API Call #4
	if isInvite {
		method = "POST"
		recipientUrl = recipientUrl + collectorID + "/messages/" + messageID + "/recipients/bulk"
		emails := strings.Split(recipientList, ",")
		emailParent_Json := `{ "contacts": [`
		count := 0
		for i := 0; i < len(emails); i++ {
			innerJsonContent, _ := sjson.Set("", "email", emails[i])
			emailParent_Json = emailParent_Json + innerJsonContent
			if count < len(emails)-1 {
				count = count + 1
				emailParent_Json = emailParent_Json + ","
			}
		}
		emailParent_Json = emailParent_Json + "]}"
		jsonBody = []byte(emailParent_Json)
		reqRecipientBulk, err := callUrl(method, recipientUrl, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Errorf("error while attatching a message body :[%s]", err.Error())
			return false, err
		} else {
			succStatus_1 := gjson.Get(reqRecipientBulk, "succeeded.#").String()
			succStatus_2 := gjson.Get(reqRecipientBulk, "existing.#").String()
			if succStatus_1 != "" || succStatus_2 != "" {
				log.Infof("emails added successfully")
			} else {
				return false, errors.New("error while attaching bulk emails")
			}
		}
	}

	//add schedule email date , API Call #5
	method = "POST"
	sendUrl = sendUrl + collectorID + "/messages/" + messageID + "/send"
	curr_date := time.Now().Format("2006-01-02T15:04:05+00:00")
	jsonBody = []byte(`{ "scheduled_date": "` + curr_date + `"}`)
	_, err = callUrl(method, sendUrl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Errorf("error while sending emails: [%s]", err.Error())
		return false, err
	}
	return true, nil
}
