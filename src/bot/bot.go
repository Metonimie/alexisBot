package bot

import (
    "os"
    "github.com/paked/messenger"
    "math/rand"
    "time"
    "log"
)

var token = os.Getenv("WIT_API_TOKEN")

// Messages arrays
var greetingMessages []string = []string{"Hello there ", "Hey, ", "Hi "}
var goodbyeMessages []string = []string{"Goodbye ", "Bye ", "See you later "}

// Holds all variables needed to respond to avoid long argument lists.
type Response struct {
    profile *messenger.Profile
    message *messenger.Message
    response *messenger.Response
    witResponse *MessageResponse // TODO: Add polymorphism.
}

// Try to make a single wit.ai Client
var witAiClient *Client

// Initialize the wit.ai client
func initClient()  {
    if witAiClient == nil {
        witAiClient = NewWitClient(token)
        rand.Seed(time.Now().Unix())
    }
}

// Get a random string from a slice
func getRandomString(slice []string) string {
    return slice[rand.Intn(len(slice))]
}

// If the response contains both bye and greeting entity then something is fishy.
func invalidEntities(r *MessageResponse) bool {
    if r.ContainsEntity("greetings") && r.ContainsEntity("bye") {
        return true
    }
    return false
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
        log.Println(err)
        return err
    }

    // Parse the response into Go structs
    parsedResponse, err := ParseMessageResponse(response)
    if err != nil {
        return err
    }

    // Construct new Response struct
    responseObject := new(Response)
    responseObject.response = r
    responseObject.message = m
    responseObject.profile = p
    responseObject.witResponse = parsedResponse

    respondToMessage(responseObject)

    return nil
}

// Deals with the responding logic.
func respondToMessage(response *Response) {
    // Check if wit.ai understood the message.
    if invalidEntities(response.witResponse) {
        response.response.Text("Sorry. I didn't understand. :(")
        response.response.Text("Try typing: help.")
        return
    }

    if response.witResponse.ContainsEntity("help") {
        response.response.Text("My only purpose is to tell you the schedule.")
        response.response.Text("Try asking me:")
        response.response.Text("What's my next course and lab?")
        response.response.Text("What's the schedule on Monday?")
        response.response.Text("Something like that. :)")
    } else if response.witResponse.ContainsEntity("greetings") {
        response.response.Text(getRandomString(greetingMessages[:]) + response.profile.FirstName + "!")
    } else if response.witResponse.ContainsEntity("bye") {
        response.response.Text(getRandomString(goodbyeMessages[:]) + response.profile.FirstName + " :)")
    } else {
        response.response.Text("I understood what you've said, but...")
        response.response.Text("I'm currently not able to provide you with the info.")
        response.response.Text("Try again later and stay tuned! :)")
    }
}