module tux.tech/x3dh/server

go 1.22.2

replace tux.tech/x3dh/core => ../x3dh_core

require tux.tech/x3dh/core v0.0.0-00010101000000-000000000000

require (
	github.com/golang/snappy v0.0.4 // indirect
	github.com/klauspost/compress v1.13.6 // indirect
	github.com/montanaflynn/stats v0.7.1 // indirect
	github.com/xdg-go/pbkdf2 v1.0.0 // indirect
	github.com/xdg-go/scram v1.1.2 // indirect
	github.com/xdg-go/stringprep v1.0.4 // indirect
	github.com/youmark/pkcs8 v0.0.0-20181117223130-1be2e3e5546d // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.15.0 // indirect
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.mongodb.org/mongo-driver v1.16.0
	go.step.sm/crypto v0.47.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
)
