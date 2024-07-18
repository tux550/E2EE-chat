module tux.tech/e2ee/client

go 1.22.2

replace tux.tech/x3dh/core => ../x3dh_core

replace tux.tech/x3dh/client => ../x3dh_client

require tux.tech/x3dh/core v0.0.0-00010101000000-000000000000 // indirect

require tux.tech/x3dh/client v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.3
	go.step.sm/crypto v0.47.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)
