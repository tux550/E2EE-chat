package e2ee_api

import (
	"encoding/json"

	x3dh_core "tux.tech/x3dh/core"
)

type InboundMessage struct {
	Method string          `json:"method"`
	Params json.RawMessage `json:"params"`
}

type RequestUserBundle struct {
	UserID string `json:"user_id"`
}

type RequestUploadBundle x3dh_core.X3DHClientBundle

type RequestSendMsg struct {
	RecipientID string                   `json:"recipient_id"`
	MessageData x3dh_core.InitialMessage `json:"message"`
}

type RequestReceiveMsg struct{}

type OutboundMessage struct {
	Method string      `json:"method"`
	Params interface{} `json:"params"`
}

type ResponseUserBundle x3dh_core.X3DHKeyBundle

type ResponseUploadBundle struct {
	Success bool `json:"success"`
}

type ResponseSendMsg struct {
	Success bool `json:"success"`
}

type ResponseReceiveMsg struct {
	//SenderID    string                   `json:"sender_id"`
	MessageData x3dh_core.InitialMessage `json:"message"`
}
