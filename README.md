Go SMTP Server
==============

Package `smtp` provides a simple wrapper around [MailHog](https://github.com/mailhog/MailHog)'s SMTP server
written in a style that's hopefully reminiscent of `net/http`. Most of the
code is adapted from [MailHog-Server's smtp package](https://github.com/mailhog/MailHog-Server/tree/master/smtp).

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
