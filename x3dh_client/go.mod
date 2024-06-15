module x3dh_client

go 1.22.2

require go.step.sm/crypto v0.47.1
replace x3dh_core => ../x3dh_core

require x3dh_core v0.0.0-00010101000000-000000000000

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)
