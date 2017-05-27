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

// Entity struct will hold infromation about wit.ai's entities
type Entity struct {
    Name string // Name will be used to determine what the next action will be.
    Value string `json:"value"`
    Type string `json:"type"`
    Suggested bool `json:"suggested"`
    Confidence float64 `json:"confidence"`
}

// MessageResponse struct will hold only important data from wit.ai needed by our logic.
type MessageResponse struct {
    Text string `json:"_text"`
    Entities []*Entity `json:"entities"`
}

// Check if the message response contains an entity with name x
func (ms *MessageResponse) ContainsEntity(name string) bool {
    for _, entity := range ms.Entities {
        if entity.Name == name {
            return true
        }
    }
    return false
}

// ParseMessageResponse will translate the wit's ai json response into a Go MessageResponse struct.
func ParseMessageResponse(data []byte) (*MessageResponse, error)  {
    // Container to hold all the json's data.
    var container interface{}
    response := new(MessageResponse)

    err := json.Unmarshal(data, &container)
    if err != nil {
        return nil, err
    }

    // Translate container to a map type.
    m := container.(map[string]interface{})

    if m["_text"] != nil {
        response.Text = m["_text"].(string)
    }
    if m["entities"] != nil {
        // Construct a new map of entities
        entities := m["entities"].(map[string]interface{})
        // For each entity object that is present in entities
        for k, v := range entities {
            entity := new(Entity)
            entity.Name = k

            // Construct a new map of entities values
            entityValues := v.([]interface{})[0].(map[string]interface{})
            for ke, ve := range entityValues {
                if ke == "type" {
                    entity.Type = ve.(string)
                } else if ke == "value" {
                    entity.Value = ve.(string)
                } else if ke == "suggested" {
                    entity.Suggested = ve.(bool)
                } else if ke == "confidence" {
                    entity.Confidence = ve.(float64)
                }
            }
            // Add the entity to the response's entities array.
            response.Entities = append(response.Entities, entity)
        }
    }
    return response, err
}

// Will translate a json from wit.ai to a go struct
func ParseConverseResponse(data []byte) (*ConverseResponse, error) {
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