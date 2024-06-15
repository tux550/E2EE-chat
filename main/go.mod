module main

go 1.22.2

replace tux.tech/x3dh/client => ../x3dh_client

replace tux.tech/x3dh/server => ../x3dh_server

replace tux.tech/x3dh/core => ../x3dh_core

require tux.tech/x3dh/client v0.0.0-00010101000000-000000000000

require tux.tech/x3dh/server v0.0.0-00010101000000-000000000000

require tux.tech/x3dh/core v0.0.0-00010101000000-000000000000 // indirect

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.step.sm/crypto v0.47.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)
