# SurveyMonkey_EmailInvite

This activity allows the user to send Email invitations for the surveys to a user group, reminding them and sending thank you emails upon successfully completing the surveys.It makes set of API calls to achieve this functionality.

## Installation

### Flogo CLI

```
flogo install github.com/pankajsankpal/SurveyMonkey_EmailInvite/sendEmailInvite
```

### Third-party libraries used
- #### GJSON :
GJSON is a Go package that provides a fast and simple way to get values from a json document. It has features such as one line retrieval, dot notation paths, iteration, and parsing json lines.
- #### SJSON :
SJSON is a Go package that provides a very fast and simple way to set a value in a json document. The purpose for this library is to provide efficient json updating in the SurveyMonkey_GetResponse activity.

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
| type  | True | String | type of email(invite,reminder,thank_you) |
| recipient_status  | False | String | Used In case of reminder and thank_you message (has_not_completed,completed,responded) |
| recipientList  | False | String | comma(,) separated list of user |
| Subject  | False | String | Subject of the email message to be sent to recipients |
| Body  | False | String |  |



### Example :
This activity willaccepts the inputs in following way,

```
