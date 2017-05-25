//MIT License
//
//Copyright (c) 2017 Denis <denis.nutiu@gmail.com>
//
//Permission is hereby granted, free of charge, to any person obtaining a copy
//of this software and associated documentation files (the "Software"), to deal
//in the Software without restriction, including without limitation the rights
//to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
//copies of the Software, and to permit persons to whom the Software is
//furnished to do so, subject to the following conditions:
//
//The above copyright notice and this permission notice shall be included in all
//copies or substantial portions of the Software.
//
//THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
//IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
//FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
//AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
//LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
//OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
//SOFTWARE.

package src

import (
    "net/http"
    "log"
    "golang.org/x/crypto/acme/autocert"
    "crypto/tls"
    "github.com/paked/messenger"
    "fmt"
    "time"
    "strings"
    "os"
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
    _, err := client.ProfileByID(m.Sender.ID)
    if err != nil {
        log.Println("Something went wrong!", err)
    }

    // Spongebob's algorithm
    var newS string
    counter := 0
    for _, c := range m.Text {
         if counter % 2 != 0 {
              newS += strings.ToLower(string(c))
         } else {
              newS += strings.ToUpper(string(c))
         }
         counter++
    }

    r.Text(newS)
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
