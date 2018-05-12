# SurveyMonkey_EmailInvite

This activity allows the user to send Email invitations for the surveys to a specific user group and reminding them.

## Installation

### Flogo CLI

```
flogo install github.com/pankajsankpal/SurveyMonkey_EmailInvite/sendemailinvite
```

### Third-party libraries used
- #### GJSON :
GJSON is a Go package that provides a fast and simple way to get values from a json document. It has features such as one line retrieval, dot notation paths, iteration, and parsing json lines.
- #### SJSON :
SJSON is a Go package that provides a simple way to set a value in a json document. The purpose for this library is to provide efficient json updating in the SurveyMonkey_EmailInvite activity.

### Schema

```
{
"inputs":[
    {
      "name": "surveyName",
      "type": "string",
	  "required": true
    },
	{
      "name": "authToken",
      "type": "string",
	  "required": true
    },
	{
      "name": "senderEmail",
      "type": "string",
	  "required": true
    },
	{
      "name": "type",
      "type": "string",
	  "allowed": [
        "invite",
        "reminder",
        "thank_you"
      ],
	  "value": "invite",
    "required": true
    },
    {
        "name": "recipientStatus",
        "type": "string",
        "allowed": [
            "has_not_responded",
            "completed",
            "responded"
          ],
        "value": ""
    },
    {
      "name": "recipientList",
      "type": "string"
    },
    {
      "name": "subject",
      "type": "string"
    },
    {
      "name": "body",
      "type": "string"
    }
  ]
}
```

### Activity Input


| Name | Required | Type | Description |
| ---- | -------- | ---- |------------ |
| surveyName | True | String | Name of the survey |
| authToken | True | String | Authentication Token for user |
| senderEmail  | True | String | Sender email for email collectors|
| type  | True | String | Type of email(invite,reminder) |
| recipientStatus  | False | String | If type is 'reminder', acceptable values are: 'has_not_responded’ or 'partially_responded’, with the default being 'has_not_responded’.|
| recipientList  | False | String | comma(,) separated list of user |
| subject  | False | String | subject of the email message to be sent to recipients |
| body  | False | String |  Body for the email|



### Example :
This activity will accept the inputs in following way,

```
{
          "id": "sendemailinvite_2",
          "name": "sendemailinvite",
          "description": "activity description",
          "type": 1,
          "activityType": "sendemailinvite",
          "activityRef": "github.com/pankajsankpal/SurveyMonkey_EmailInvite/sendemailinvite",
          "attributes": [
            {
              "name": "surveyName",
              "value": "DemoServey",
              "required": true,
              "type": "string"
            },
            {
              "name": "authToken",
              "value": "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ",
              "required": true,
              "type": "string"
            },
            {
              "name": "senderEmail",
              "value": "sankpal22pankaj@gmail.com",
              "required": true,
              "type": "string"
            },
            {
              "name": "type",
              "value": "invite",
              "required": true,
              "type": "string"
            },
            {
              "name": "recipientStatus"
              "value": "",
              "required": false,
              "type": "string"
            },
            {
              "name": "recipientList",
              "value": "psankpal@tibco.com,anansing@tibco.com",
              "required": false,
              "type": "string"
            },
            {
              "name": "subject",
              "value": "testEmail",
              "required": false,
              "type": "string"
            },
            {
              "name": "body",
              "value": "Body Content",
              "required": false,
              "type": "string"
            }
          ]
       }
```
