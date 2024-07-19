module tux.tech/e2ee/client

go 1.22.2

replace tux.tech/x3dh/core => ../x3dh_core

replace tux.tech/x3dh/client => ../x3dh_client

replace tux.tech/e2ee/api => ../e2ee_api

require (
	github.com/asaskevich/govalidator v0.0.0-20230301143203-a9d515a09cc2 // indirect
	github.com/go-openapi/errors v0.22.0 // indirect
	github.com/go-openapi/strfmt v0.23.0 // indirect
	github.com/google/uuid v1.6.0 // indirect
	github.com/mattn/go-runewidth v0.0.15 // indirect
	github.com/mitchellh/mapstructure v1.5.0 // indirect
	github.com/oklog/ulid v1.3.1 // indirect
	github.com/rivo/uniseg v0.2.0 // indirect
	go.mongodb.org/mongo-driver v1.14.0 // indirect
	golang.org/x/sys v0.21.0 // indirect
)

require tux.tech/x3dh/client v0.0.0-00010101000000-000000000000

require (
	github.com/jedib0t/go-pretty v4.3.0+incompatible
	tux.tech/e2ee/api v0.0.0-00010101000000-000000000000
	tux.tech/x3dh/core v0.0.0-00010101000000-000000000000
)

require (
	filippo.io/edwards25519 v1.1.0 // indirect
	github.com/gorilla/websocket v1.5.3
	go.step.sm/crypto v0.47.1 // indirect
	golang.org/x/crypto v0.24.0 // indirect
)
