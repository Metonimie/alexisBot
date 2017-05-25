// Copyright 2017 Denis Nutiu
// This file contains logic for parsing json received from wit.ai
package bot

import (
    "encoding/json"
)

// Response struct with coded values from the wit.ai api
type Response struct {
    responseType string // merge, action, msg, stop
    data string // Will be filled up based on the responseType
    confidence float64
    entities []byte // Contains raw json for later processing
}

// Will translate a json from wit.ai to a go struct
func parseResponse(data []byte) (*Response, error) {
    // Empty interface to hold all the data
    var container interface{}
    response := new(Response)

    err := json.Unmarshal(data, &container)
    if err != nil {
        return nil, err
    }

    // Construct a map from the empty interface
    m := container.(map[string]interface{})

    if m["confidence"] != nil {
        response.confidence = m["confidence"].(float64) // type assetion convert
    }
    if m["type"] != nil {
        response.responseType = m["type"].(string)
    }
    if m["msg"] != nil {
        response.data = m["msg"].(string)
    }
    if m["action"] != nil {
        response.data = m["action"].(string)
    }
    if m["entities"] != nil {
        b, e := json.Marshal(m["entities"])
        // If we manage to construct the json successfully
        if e == nil {
            response.entities = b
        }
    }

    return response, nil
}