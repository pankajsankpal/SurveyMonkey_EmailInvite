package busslogic

import (
  "bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"

	"github.com/TIBCOSoftware/flogo-lib/core/activity"
	"github.com/TIBCOSoftware/flogo-lib/logger"
	"github.com/tidwall/gjson"
	"github.com/tidwall/sjson"
)

//logger
var log = logger.GetLogger("activity-go logger")
var isInvite bool = false

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
			fmt.Println(outResult)
			return "", errors.New(outResult)
		} else {
			return string(survey_response), nil
		}
	}
}

func sendEmail( accessToken string, surveyName string, senderEmail string, recipientList string, typeofEmail string, recipientStatus string,subject string,body string) (resp bool,err error){

if typeofEmail=="invite"{
  isInvite=true
}

//get surveyID , API call #1
  method = "GET"
	log.Debugf("The Incoming Survey Name [%s]", surveyName)
	surveyIdurl := "https://api.surveymonkey.com/v3/surveys?title=" + surveyName
	var jsonBody = []byte("")
	reqSurveyID, err := callUrl(method, surveyIdurl, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Debugf("Error while fetching surveyID [%s]", err.Error())
		return false, err
	} else {
		surveyID = gjson.Get(reqSurveyID, "data.0.id").String()
		log.Debugf("srveyId: [%s]", surveyID)
		fmt.Println("surveyId: " + surveyID)
	}

	//set email invite and get Collector id , API Call #2
	collectorID := ""
	collector_url := "https://api.surveymonkey.com/v3/surveys/" + surveyID + "/collectors"
	reqCollectorID := ""
	method = "GET"
	jsonBody = []byte("")
	reqCollectorID, err = callUrl(method, collector_url, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Debugf("Error while fetching collectorID: [%s]", err.Error())
		return false, err
	} else {
		collectorList := gjson.Get(reqCollectorID, "data")
		for _, item := range collectorList.Array() {
			if strings.Contains(gjson.Get(item.String(), "name").String(), "Email") {
				collectorID = gjson.Get(item.String(), "id").String()
				fmt.Println("Collector ID: " + collectorID)
				break
			}
		}
	}
	if collectorID == "" {
		method = "POST"
		jsonBody = []byte(`{"type":"email","sender_email":"` + senderEmail + `"}`)
		reqCollectorID, err = callUrl(method, collector_url, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Debugf("Error while fetching collectorID [%s] ", err.Error())
			return false, err
		} else {
			collectorID = gjson.Get(reqCollectorID, "id").String()
			fmt.Println("Collector ID: " + collectorID)
		}
	}

	//get message ID , APICall #3
	method = "POST"
	messageID := ""
	message_url := "https://api.surveymonkey.com/v3/collectors/" + collectorID + "/messages"
	if body != "" {
		surveyLink := "[SurveyLink]"
		optLink := "[OptOutLink]"
		footerLink := "[FooterLink]"
		if isInvite {
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `","body_text":"` + body + "<a href=" + "\\" + "\"" + surveyLink + "\\" + "\"" + " >Take the survey!</a> <a href=" + "\\" + "\"" + optLink + "\\" + "\"" + ">Please remove me from your mailing list.</a> <a href=" + "\\" + "\"" + footerLink + "\\" + "\"" + ">Footer!</a>" + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","recipient_status":"` + recipientStatus + `","body_text":"` + body + "<a href=" + "\\" + "\"" + surveyLink + "\\" + "\"" + " >Take the survey!</a> <a href=" + "\\" + "\"" + optLink + "\\" + "\"" + ">Please remove me from your mailing list.</a> <a href=" + "\\" + "\"" + footerLink + "\\" + "\"" + ">Footer!</a>" + `"}`)
		}
	} else {
		if isInvite {
			jsonBody = []byte(`{"type":"invite","subject":"` + subject + `"}`)
		} else {
			jsonBody = []byte(`{"type":"` + typeofEmail + `","subject":"` + subject + `","recipient_status":"` + recipientStatus + `"}`)
		}
	}
	reqMessageID, err := callUrl(method, message_url, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Debugf("Error while fetching messageID: [%s]", err.Error())
		return false, err
	} else {
		messageID = gjson.Get(reqMessageID, "id").String()
		fmt.Println("Message ID: " + messageID)
	}

	//add multiple email ids , API Call #4
	if isInvite {
		method = "POST"
		recipient_url := "https://api.surveymonkey.com/v3/collectors/" + collectorID + "/messages/" + messageID + "/recipients/bulk"
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
		reqRecipientBulk, err := callUrl(method, recipient_url, bytes.NewBuffer(jsonBody), accessToken)
		if err != nil {
			log.Debugf("Error while attatching a message body :[%s]", err.Error())
			return false, err
		} else {
			succStatus_1 := gjson.Get(reqRecipientBulk, "succeeded.#").String()
			succStatus_2 := gjson.Get(reqRecipientBulk, "existing.#").String()
			if succStatus_1 != "" || succStatus_2 != "" {
				fmt.Println("emails added successfully...")
			} else {
				context.SetOutput("status", `{ "Error" : { "message" : "Error while adding recipients emails." }}`)
				return true, nil
			}
		}
	}

	//add schedule email date , API Call #5
	method = "POST"
	send_url := "https://api.surveymonkey.com/v3/collectors/" + collectorID + "/messages/" + messageID + "/send"
	curr_date := time.Now().Format("2006-01-02T15:04:05+00:00")
	jsonBody = []byte(`{ "scheduled_date": "` + curr_date + `"}`)
	_, err = callUrl(method, send_url, bytes.NewBuffer(jsonBody), accessToken)
	if err != nil {
		log.Debugf("Error while sending emails: [%s]", err.Error())
		return false, err
	} else {
		fmt.Println("Email sent succesfully")
	}
  return true,nil
}
