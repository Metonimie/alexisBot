package bot

import (
    "testing"
    "os"
)

func TestNewWitClient(t *testing.T) {
    witClient := NewWitClient(os.Getenv("WIT_API_TOKEN"))
    if witClient == nil {
        t.Fail()
    }
}

func TestClient_Message(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipped: TestClient_Message")
    }
    witClient := NewWitClient(os.Getenv("WIT_API_TOKEN"))
    resp, err := witClient.Message("Hello", nil, nil, nil)
    if err != nil {
        t.Fail()
    }
    if resp == nil {
        t.Fail()
    }
}

func TestClient_Converse(t *testing.T) {
    if testing.Short() {
        t.Skip("Skipped: TestClient_Converse")
    }
    witClient := NewWitClient(os.Getenv("WIT_API_TOKEN"))
    resp, err := witClient.Message("Hello", nil, nil, nil)
    if err != nil {
        t.Fail()
    }
    if resp == nil {
        t.Fail()
    }
}

