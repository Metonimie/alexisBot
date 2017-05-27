package bot

import (
    "testing"
    "github.com/paked/messenger"
)

func TestHandleMessage(t *testing.T) {
    if testing.Short() {
        t.Skip()
    }
    profile := messenger.Profile{LastName: "Test", FirstName: "Test"}
    message := messenger.Message{Text: "What is the next course?"}
    response := messenger.Response{}

    err := HandleMessage(&message, &response, &profile)
    if err != nil {
        t.Fail()
    }
}
