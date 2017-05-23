package main

import (
    "net/http"
    "log"
    "fmt"
)

const ip = "" // empty means it binds to everything.
const port = "8080"

// The webhook functions receives challenges from Facebook.
func webhook(w http.ResponseWriter, r * http.Request) {
    if r.Method == "POST" {
        fmt.Println(r.Body)
    } else {
        // If we don't have POST then we don't POST!!
        http.Error(w, "Invalid request method.", 405)
    }
}

func main() {
    // Handle functions
    http.HandleFunc("/webhook", webhook)



    log.Println("Hello world, I think, therefore I am.")
    log.Println("Running on", "http://"+ ip + ":" + port)
    // Start the server and log.
    log.Fatalln(http.ListenAndServe(ip + ":" + port, nil))
}
