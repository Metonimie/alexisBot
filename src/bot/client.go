// Package bot
// Copyright 2017 Denis Nutiu
// The client file contains the logic to communicate with the Wit.ai API
//
package bot

import (
	"bytes"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

// TODO: Research timeout for the requests to wit.

const apiURL = "https://api.wit.ai"
const apiVersion = "20170307"

// Client struct that holds all the necessary data for a successful communication with Wit.ai
type Client struct {
	apiUrl     string
	apiKey     string
	apiVersion string
	userAgent  string
}

type queryValues map[string]string

// Create a new Wit Client.
//
//    client := bot.Client("WIT TOKEN")
func NewWitClient(witToken string) *Client {
	client := new(Client)
	client.apiKey = witToken
	client.apiUrl = apiURL
	client.apiVersion = apiVersion
	client.userAgent = "AlexisBot"
	return client
}

// Constructs an url ready to be used for the request.
//
//    query := make(queryValues)
//    query["session_id"] = "abc123"
//    theUrl := client.makeURL("/converse", query)
//
//    theUrl will have the following value: https://api.wit.ai/converse?v=123&session_id=abc123
func (client *Client) makeURL(path string, args queryValues) string {
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

// Execute the request for the give url.
// This functions sets the necessary headers and content-type and returns
// back the response a bytes array.
func (client *Client) executeRequest(url string) ([]byte, error) {
	// Create the request
	request, err := http.NewRequest("POST", url, nil)
	if err != nil {
		return nil, err
	}
	// Set the headers.
	request.Header.Add("Content-Type", "application/json")
	request.Header.Add("Authorization", "Bearer "+client.apiKey)

	// Execute the request
	httpClient := http.Client{}
	resp, err := httpClient.Do(request)
	defer resp.Body.Close() // Close when this function returns.
	if err != nil {
		return nil, err
	}

	// Read the response
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

// Converse calls the wit.ai /converse endpoint
// See more: https://wit.ai/docs/http/20170307#post--converse-link
func (client *Client) Converse(sessionId string, q string, reset bool) ([]byte, error) {
	// Make a query map
	query := make(queryValues)
	query["session_id"] = sessionId

	// Properly encode the query before sending it.
	cleanQ, err := url.Parse(q)
	if err != nil {
		return nil, err
	}
	query["q"] = cleanQ.String()

	// Add reset to request if provided.
	if reset == true {
		query["reset"] = "true"
	}
	// Get the url for the post request.
	theUrl := client.makeURL("/converse", query)

	// Execute request
	response, err := client.executeRequest(theUrl)
	if err != nil {
		log.Println(err.Error())
	}

	return response, nil
}

// Message calls the wit.ai /message endpoint
// Only q is required, others are optional
// See more: https://wit.ai/docs/http/20170307#get--message-link
func (client *Client) Message(q string, msg_id *string, thread_id *string, n *string) ([]byte, error) {
	// Make a query map
	query := make(queryValues)

	// Properly encode the query before sending it.
	cleanQ := url.QueryEscape(q)
	query["q"] = cleanQ

	if msg_id != nil {
		query["msg_id"] = *msg_id
	}
	if thread_id != nil {
		query["thread_id"] = *thread_id
	}
	if n != nil {
		query["n"] = *n
	}

	// Get the url for the post request.
	theUrl := client.makeURL("/message", query)

	// Execute request
	response, err := client.executeRequest(theUrl)
	if err != nil {
		log.Println(err.Error())
	}

	return response, nil
}
