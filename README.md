# go-imap-proxy

A [go-imap](https://github.com/emersion/go-imap) backend that proxies all commands and responses to another IMAP server.

## Usage

```go
package main

import (
	"log"

	"github.com/emersion/go-imap/server"
	"github.com/emersion/go-imap-proxy"
)

func main() {
	be := proxy.NewTLS("mail.example.org:993", nil)

	// Create a new server
	s := server.New(be)
	s.Addr = ":1143"
	// Since we will use this server for testing only, we can allow plain text
	// authentication over unencrypted connections
	s.AllowInsecureAuth = true

	log.Println("Starting IMAP server at localhost:1143")
	if err := s.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}
```

## License

MIT
