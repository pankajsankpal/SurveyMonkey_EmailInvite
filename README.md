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
      "name": "Survey Name",
      "type": "string",
	  "required": true
    },
	{
      "name": "AuthToken",
      "type": "string",
	  "required": true
    },
	{
      "name": "Sender's Email",
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
        "name": "recipient_status",
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
      "name": "Subject",
      "type": "string"
    },
    {
      "name": "Body",
      "type": "string"
    }
  ]
}
```

### Activity Input


| Name | Required | Type | Description |
| ---- | -------- | ---- |------------ |
| Survey Name | True | String | Name of the survey |
| AuthToken | True | String | Authentication Token for user |
| Sender's Email  | True | String | Sender email for email collectors|
| Type  | True | String | Type of email(invite,reminder) |
| Recipient Status  | False | String | If type is 'reminder’, acceptable values are: 'has_not_responded’ or 'partially_responded’, with the default being 'has_not_responded’.|
| Recipient List  | False | String | comma(,) separated list of user |
| Subject  | False | String | Subject of the email message to be sent to recipients |
| Body  | False | String |  Body for the email|



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
              "name": "Survey Name",
              "value": "DemoServey",
              "required": true,
              "type": "string"
            },
            {
              "name": "AuthToken",
              "value": "z8UFEI9i5ua1WWhI40S1xo8yLlFJFsOPMdwtsB83YYAJy.1fr.zPLQ9mfrh7a2qTZHqdCwwnMHHn9.U0OvXcyx5SjYLRjcMUsE-YE6mcZAB0fg4lP2zoDNg-sL8fxDoQ",
              "required": true,
              "type": "string"
            },
            {
              "name": "Sender's Email",
              "value": "sankpal22pankaj@gmail.com",
              "required": true,
              "type": "string"
            },
            {
              "name": "Type",
              "value": "invite",
              "required": true,
              "type": "string"
            },
            {
              "name": "Recipient Status"
              "value": "",
              "required": false,
              "type": "string"
            },
            {
              "name": "Recipient List",
              "value": "psankpal@tibco.com,anansing@tibco.com",
              "required": false,
              "type": "string"
            },
            {
              "name": "Subject",
              "value": "testEmail",
              "required": false,
              "type": "string"
            },
            {
              "name": "Body",
              "value": "Body Content",
              "required": false,
              "type": "string"
            }
          ]
       }
```
