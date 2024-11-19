module tux.tech/e2ee/server

go 1.23

replace tux.tech/e2ee/api => ../e2ee_api

replace tux.tech/x3dh/core => ../x3dh_core

replace tux.tech/x3dh/server => ../x3dh_server

require tux.tech/x3dh/core v0.0.0-00010101000000-000000000000 // indirect

require tux.tech/x3dh/server v0.0.0-00010101000000-000000000000

require tux.tech/e2ee/api v0.0.0-00010101000000-000000000000

require github.com/gorilla/websocket v1.5.3

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/aws/aws-lambda-go v1.47.0
	github.com/aws/aws-sdk-go v1.55.5 // indirect
	github.com/golang/snappy v0.0.4 // indirect
	github.com/jmespath/go-jmespath v0.4.0 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20240726163527-a2c0da244d78 // indirect
	go.mongodb.org/mongo-driver v1.17.1 // indirect
	go.step.sm/crypto v0.54.0 // indirect
	golang.org/x/crypto v0.29.0 // indirect
	golang.org/x/sync v0.9.0 // indirect
	golang.org/x/text v0.20.0 // indirect
)
