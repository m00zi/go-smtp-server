Go SMTP Server
==============

Package `smtp` provides a simple wrapper around MailHog's SMTP server
written in a style that's hopefully reminiscent of `net/http`.

Usage
-----

```golang
server := &smtp.Server{
	Addr: ":2500",
	Hostname: "smtp.example.org",
	Handler: func(message smtp.Message) error {
		log.Print(message.Content.Body)
		return nil
	},
}

server.Listen()
```
