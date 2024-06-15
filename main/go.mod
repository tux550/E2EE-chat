module main

go 1.22.2

replace x3dh_core => ../x3dh_core

replace x3dh_client => ../x3dh_client

replace x3dh_server => ../x3dh_server

require x3dh_core v0.0.0-00010101000000-000000000000

require x3dh_client v0.0.0-00010101000000-000000000000

require x3dh_server v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.step.sm/crypto v0.47.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)
