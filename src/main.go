//Copyright (c) 2017 Denis <denis.nutiu@gmail.com>
package main

import (
    "net/http"
    "log"
    "golang.org/x/crypto/acme/autocert"
    "crypto/tls"
    "github.com/paked/messenger"
    "fmt"
    "time"
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

// Handler to be triggered when a message is received
func MessageHandler(m messenger.Message, r *messenger.Response) {
    log.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

    // p is delcared here.
    p, err := client.ProfileByID(m.Sender.ID)
    if err != nil {
        log.Println("Something went wrong!", err)
    }

    err = bot.HandleMessage(&m, r, &p)

    if err != nil {
        r.Text("Aaaahhhahahhh! I'm experiecingg.....")
        r.Text("A bufferrr, oveerrfllowooooow.")
        r.Text("Jk, I'm written in Go.")
        r.Text("www.golang.com :>")
    }
}

// Handler to be triggered when a message is delivered
func DeliveryHandler(d messenger.Delivery, r *messenger.Response) {
    fmt.Println("Delivered at:", d.Watermark().Format(time.UnixDate))
}

// Handler to be triggered when a message is read
func ReadHandler(m messenger.Read, r *messenger.Response) {
    fmt.Println("Read at:", m.Watermark().Format(time.UnixDate))
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
    // Setup a handler to be triggered when a message is read
    client.HandleRead(ReadHandler)
    // Setup a handler to be triggered when a message is delivered
    client.HandleDelivery(DeliveryHandler)

    log.Println("Hello world, I think, therefore I am.")
    // Start the server and log.
    log.Fatalln(server.ListenAndServeTLS( "", ""))
}
