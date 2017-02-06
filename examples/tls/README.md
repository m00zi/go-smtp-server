# Example using TLS

## Usage

Generate keys in PEM format:

`openssl req -x509 -nodes -days 365 -newkey rsa:2048 -keyout selfsigned.key -out selfsigned.crt`

Run:

`go run examples/tls/main.go -addr=:1025 -cert-file=selfsigned.crt -key-file=selfsigned.key`

Test the server using any client. For example, [swaks](http://www.jetmore.org/john/code/swaks/):

`swaks --to whatever@localhost --server localhost:1025 --tls`
