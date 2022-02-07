module github.com/pflege-de/cc-connection

go 1.17

require github.com/nats-io/nats.go v1.13.1-0.20220121202836-972a071d373d

require (
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/nats-io/nats-server/v2 v2.7.2 // indirect
	google.golang.org/protobuf v1.27.1 // indirect
)

require (
	github.com/golang-jwt/jwt/v4 v4.2.0
	github.com/google/uuid v1.3.0
	github.com/nats-io/nkeys v0.3.0 // indirect
	github.com/nats-io/nuid v1.0.1 // indirect
	github.com/nimajalali/go-force v0.0.0-20200831220737-454890ee2b7c
	github.com/pkg/errors v0.9.1
	golang.org/x/crypto v0.0.0-20220112180741-5e0467b6c7ce // indirect
)

replace github.com/nimajalali/go-force => ./go-force
