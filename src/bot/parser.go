// Copyright 2017 Denis Nutiu
// This file contains logic for parsing json received from wit.ai
package bot

import (
    "encoding/json"
)

// ConverseResponse struct with coded values from the wit.ai api
type ConverseResponse struct {
    Type       string // merge, action, msg, stop
    Data       string // Will be filled up based on the Type
    Confidence float64
    Entities   []byte // Contains raw json for later processing
}

// Will translate a json from wit.ai to a go struct
func parseConverseResponse(data []byte) (*ConverseResponse, error) {
    // Empty interface to hold all the Data
    var container interface{}
    response := new(ConverseResponse)

    err := json.Unmarshal(data, &container)
    if err != nil {
        return nil, err
    }

    // Construct a map from the empty interface
    m := container.(map[string]interface{})

    if m["Confidence"] != nil {
        response.Confidence = m["Confidence"].(float64) // type assetion convert
    }
    if m["type"] != nil {
        response.Type = m["type"].(string)
    }
    if m["msg"] != nil {
        response.Data = m["msg"].(string)
    }
    if m["action"] != nil {
        response.Data = m["action"].(string)
    }
    if m["Entities"] != nil {
        b, e := json.Marshal(m["Entities"])
        // If we manage to construct the json successfully
        if e == nil {
            response.Entities = b
        }
    }

    return response, nil
}