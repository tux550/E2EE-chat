module tux.tech/e2ee/server

go 1.22.2

replace tux.tech/e2ee/api => ../e2ee_api

replace tux.tech/x3dh/core => ../x3dh_core

replace tux.tech/x3dh/server => ../x3dh_server

require tux.tech/x3dh/core v0.0.0-00010101000000-000000000000 // indirect

require tux.tech/x3dh/server v0.0.0-00010101000000-000000000000

require tux.tech/e2ee/api v0.0.0-00010101000000-000000000000

require github.com/gorilla/websocket v1.5.3

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.step.sm/crypto v0.47.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
)
