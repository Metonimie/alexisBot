package main

import (
    "net/http"
    "log"
    "golang.org/x/crypto/acme/autocert"
    "crypto/tls"
    "github.com/paked/messenger"
    "fmt"
    "time"
)

const host = "denisnutiu.me"
const ip = "" // empty means it binds to everything.
const port = "443" // Facebook wants us to be "secure". FaCeeBoK WanTs Us To bE sEcUrE.

const verify = false // If the app should verify itself.
const verifyToken = "" // The facebook verify token.
const pageToken = "" // The facebook page token.

var client *messenger.Messenger

// Handler to be triggered when a message is received
func MessageHandler(m messenger.Message, r *messenger.Response) {
    fmt.Printf("%v (Sent, %v)\n", m.Text, m.Time.Format(time.UnixDate))

    p, err := client.ProfileByID(m.Sender.ID)
    if err != nil {
        fmt.Println("Something went wrong!", err)
    }

    r.Text(fmt.Sprintf("Hello, %v!", p.FirstName))
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

    // Server settings
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
