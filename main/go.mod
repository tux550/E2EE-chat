module main

go 1.22.2

replace x3dh => ../x3dh

require x3dh v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	go.step.sm/crypto v0.47.0 // indirect
	golang.org/x/crypto v0.23.0 // indirect
)
