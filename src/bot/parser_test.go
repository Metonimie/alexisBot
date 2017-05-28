package bot

import (
	"fmt"
	"testing"
)

var data = `{
  "Confidence" : 0.08051390588273054,
  "type" : "msg",
  "msg" : "Hello there",
  "Entities" : {
    "greetings" : [ {
      "Confidence" : 1.0,
      "value" : "true"
    } ]
  }
}`

var messageData = `{
  "msg_id" : "a833cf6f-83c7-4702-9db1-40c6c30a0b87",
  "_text" : "Hello",
  "entities" : {
    "reminder" : [ {
      "confidence" : 0.7916673811727122,
      "entities" : {
        "contact" : [ {
          "confidence" : 0.5748402256063738,
          "type" : "value",
          "value" : "Hello",
          "suggested" : true
        } ]
      },
      "type" : "value",
      "value" : "Hello",
      "suggested" : true
    } ],
    "greetings" : [ {
      "confidence" : 1.0,
      "value" : "true"
    } ]
  }
}
`

func TestParseConverseResponse(t *testing.T) {
	t.Parallel()
	pdata, err := ParseConverseResponse([]byte(data))

	if err != nil {
		t.Fail()
	}
	if pdata.Data != "Hello there" {
		t.Fail()
	}
	if pdata.Type != "msg" {
		t.Fail()
	}
}

func TestParseMessageResponse(t *testing.T) {
	t.Parallel()
	response, err := ParseMessageResponse([]byte(messageData))

	if err != nil {
		t.Fail()
	}

	if response.Text != "Hello" {
		t.Fail()
	}
}

func ExampleParseMessageResponse() {
	response, _ := ParseMessageResponse([]byte(messageData))

	fmt.Println(response.Entities[0].Name)
	fmt.Println(response.Entities[1].Name)
	// Output:
	// reminder
	// greetings
}

func TestMessageResponse_ContainsEntity(t *testing.T) {
	t.Parallel()

	messageResponse := MessageResponse{}
	messageResponse.Entities = append(messageResponse.Entities, &Entity{Name: "calme"})

	if messageResponse.ContainsEntity("calme") != true {
		t.Fail()
	}
	if messageResponse.ContainsEntity("calmee") != false {
		t.Fail()
	}
}
