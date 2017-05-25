// Copyright 2017 Denis Nutiu
// The client file contains the logic to communicate with the Wit.ai API
//
package bot

import (
    "net/http"
    "bytes"
    "io/ioutil"
    "net/url"
)

const apiUrl = "https://api.wit.ai"
const apiVersion  = "20170307"

type Client struct {
    apiUrl string
    apiKey string
    apiVersion string
    userAgent string
}

type queryValues map[string]string

// Create a new Wit Client.
//
//    client := bot.Client("WIT TOKEN")
func NewWitClient(witToken string) *Client {
    client := new(Client)
    client.apiKey = witToken
    client.apiUrl = apiUrl
    client.apiVersion = apiVersion
    client.userAgent = "AlexisBot"
    return client
}

// Constructs an url ready to be used for the request.
//
//    query := make(queryValues)
//    query["session_id"] = "abc123"
//    theUrl := client.makeUrl("/converse", query)
//
//    theUrl will have the following value: https://api.wit.ai/converse?v=123&session_id=abc123
func (client *Client) makeUrl(path string, args queryValues) string {
    var buffer bytes.Buffer
    buffer.WriteString(client.apiUrl)
    buffer.WriteString(path)
    buffer.WriteString("?v=")
    buffer.WriteString(client.apiVersion)

    // For each key in queryValues add it to the url.
    for k, v := range args {
        buffer.WriteString("&" + k + "=")
        buffer.WriteString(v)
    }

    return buffer.String()
}

// Call the wit.ai /converse endpoint
// See more: https://wit.ai/docs/http/20170307#post--converse-link
func (client * Client) Converse(sessionId string, q string, reset bool) (string, error) {
    // Make a query map
    query := make(queryValues)
    query["session_id"] = sessionId

    // Properly encode the query before sending it.
    cleanQ, err := url.Parse(q)
    if err != nil {
        return "", err
    }
    query["q"] = cleanQ.String()

    // Add reset to request if provided.
    if reset == true {
        query["reset"] = "true"
    }
    // Get the url for the post request.
    theUrl := client.makeUrl("/converse", query)

    // Create the request
    request, err := http.NewRequest("POST", theUrl, nil)
    if err != nil {
        return "", err
    }
    // Set the headers.
    request.Header.Add("Content-Type", "application/json")
    request.Header.Add("Authorization", "Bearer " + client.apiKey)

    // Execute the request
    httpClient := http.Client{}
    resp, err := httpClient.Do(request)
    if err != nil {
        return "", err
    }

    // Read the response
    body, err := ioutil.ReadAll(resp.Body)
    if err != nil {
        return "", err
    }

    return string(body), nil
}

