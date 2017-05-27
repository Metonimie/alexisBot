package bot

import (
    "os"
    "github.com/paked/messenger"
)

var token = os.Getenv("WIT_API_TOKEN")

// Try to make a single wit.ai Client
var witAiClient *Client

// Initialize the wit.ai client
func initClient()  {
    if witAiClient == nil {
        witAiClient = NewWitClient(token)
    }
}

// Handle client requests by sending them to Wit.ai and process the request.
func HandleMessage(m *messenger.Message, r *messenger.Response, p *messenger.Profile) error {
    // Reuse the client.
    if witAiClient == nil {
        witAiClient = NewWitClient(token)
    }

    // Send text to wit.ai
    response, err := witAiClient.Message(m.Text, nil, nil, nil)
    if err != nil {
        return err
    }

    // Parse the response into Go structs
    parsedResponse, err := ParseMessageResponse(response)
    if err != nil {
        return err
    }

    r.Text("DEBUG: " + parsedResponse.Text)
    for _, v := range parsedResponse.Entities {
        r.Text(v.Name)
    }

    return nil
}