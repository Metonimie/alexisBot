package bot

import (
    "testing"
)

var data string = `{
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

func TestParseConverseResponse(t *testing.T) {
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