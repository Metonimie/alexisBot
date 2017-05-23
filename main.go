package main

import (
    "net/http"
    "log"
    "fmt"
    "golang.org/x/crypto/acme/autocert"
    "crypto/tls"
)

const host = "denisnutiu.me"
const ip = "" // empty means it binds to everything.
const port = "443" // Facebook wants us to be "secure". FaCeeBoK WanTs Us To bE sEcUrE.

// The Webhook functions receives challenges from Facebook.
func Webhook(w http.ResponseWriter, r * http.Request) {
    //if r.Method == "POST" {
        fmt.Println(r.Body)
    //} else {
    //    // If we don't have POST then we don't POST!!
    //    http.Error(w, "Invalid request method.", 405)
    //}
}

func main() {
    // Let's Encrypt Certificate Manager.
    certManager := autocert.Manager{
        Prompt:     autocert.AcceptTOS,
        HostPolicy: autocert.HostWhitelist(host), //your domain here
        Cache:      autocert.DirCache("certs"), //folder for storing certificates
    }

    // Handle functions
    http.HandleFunc("/webhook", Webhook)


    server := &http.Server{
        Addr: ip + ":" + port,
        TLSConfig: &tls.Config{
            GetCertificate: certManager.GetCertificate,
        },
    }

    log.Println("Hello world, I think, therefore I am.")
    // Start the server and log.
    //log.Fatalln(http.ListenAndServe(ip + ":" + port, nil))
    log.Fatalln(server.ListenAndServeTLS( "", ""))
}
