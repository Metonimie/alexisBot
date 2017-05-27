//Copyright (c) 2017 Denis <denis.nutiu@gmail.com>
package main

import (
    "net/http"
    "log"
    "golang.org/x/crypto/acme/autocert"
    "crypto/tls"
    "github.com/paked/messenger"
    "os"
    "bot"
)

const host = "denisnutiu.me"
const ip = "" // empty means it binds to everything.
const port = "443" // Facebook wants us to be "secure". FaCeeBoK WanTs Us To bE sEcUrE.

const verify = false // If the app should verify itself.
var verifyToken string = os.Getenv("FB_VERIFY_TOKEN") // The facebook verify token.
var pageToken string = os.Getenv("FB_PAGE_TOKEN") // The facebook page token.

var client *messenger.Messenger

// MessageHandler is triggered when a message is received
func MessageHandler(m messenger.Message, r *messenger.Response) {
    // Get the Profile
    p, err := client.ProfileByID(m.Sender.ID)
    if err != nil {
        log.Println("Something went wrong!", err)
    }

    err = bot.HandleMessage(&m, r, &p)

    if err != nil {
        log.Println("Error by " + p.FirstName + p.LastName)
        log.Println(err.Error())
        r.Text("Hold on a minute!")
        r.Text("I'm experiencing something called a buffer overflow.")
        r.Text("Jk, I'm written in Go.")
        r.Text("www.golang.com :>")
    }
}

func main() {
    // Let's Encrypt Certificate Manager.
    certManager := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(host), //your domain here
        Cache:      autocert.DirCache("certs"), //folder for storing certificates
    }

    // Messenger Settings
    client = messenger.New(messenger.Options{
        Verify:      verify,
        VerifyToken: verifyToken,
        Token:       pageToken,
        WebhookURL: "/webhook",
    })

    // Server Settings
    server := &http.Server{
        Addr: ip + ":" + port,
        TLSConfig: &tls.Config{
            GetCertificate: certManager.GetCertificate,
        },
        Handler: client.Handler(),
    }

    // Setup a handler to be triggered when a message is received
    client.HandleMessage(MessageHandler)

    log.Println("Hello world, I think, therefore I am.")
    // Start the server and log.
    log.Fatalln(server.ListenAndServeTLS( "", ""))
}
