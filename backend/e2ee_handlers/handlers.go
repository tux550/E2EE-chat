package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	api "tux.tech/e2ee/api"
)

type APIHandler struct {
	server *Server
}

func NewAPIHandler(server *Server) *APIHandler {
	return &APIHandler{
		server: server,
	}
}

// WS request forwarded by API Gateway
type HTTPNormalizedRequest struct {
	ConnectionID string `json:"connectionId"`
	// Raw json
	Body json.RawMessage `json:"body"`
}

// RESPONSES
func (h *APIHandler) SetSuccessResponse(w http.ResponseWriter) {
	// Log
	log.Println("Success")

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Success"))
}

func (h *APIHandler) SetSuccessResponseWithMessage(w http.ResponseWriter, message []byte) {
	// Log
	log.Println("Success:", message)

	w.WriteHeader(http.StatusOK)
	w.Write(message)
}

func (h *APIHandler) SetErrorResponse(w http.ResponseWriter, message string) {
	// Log
	log.Println("Error:", message)

	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(message))
}

// API
func (h *APIHandler) HandleRequests(w http.ResponseWriter, r *http.Request) {
	// If method is GET return 200 OK
	if r.Method == "GET" {
		log.Println("Health check")
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
		return
	}
	// Parse request
	request := HTTPNormalizedRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Error decoding request:", err)
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}
	// Parse post body
	log.Println("Request Body:", string(request.Body))
	message := &api.InboundMessage{}
	err = json.Unmarshal(request.Body, message)
	if err != nil {
		log.Println("Error parsing message:", err)
		http.Error(w, "Error parsing message", http.StatusBadRequest)
		return
	}
	// Handle message
	switch message.Method {
	case "echo":
		// ECHO
		h.handleEcho(w, message.Params, request.ConnectionID)
		return
	case "get_bundle":
		h.HandleGetBundle(w, message.Params, request.ConnectionID)
		return
	case "upload_bundle":
		h.HandleUploadBundle(w, message.Params, request.ConnectionID)
		return
	case "send_message":
		h.HandleSendMessage(w, message.Params, request.ConnectionID)
		return
	case "receive_message":
		h.HandleReceiveMessage(w, message.Params, request.ConnectionID)
		return
	case "status":
		h.HandleStatus(w, message.Params, request.ConnectionID)
		return
	case "upload_new_otps":
		h.HandleUploadNewOTPs(w, message.Params, request.ConnectionID)
		return
	default:
		// Invalid method
		h.SetErrorResponse(w, fmt.Sprintf("Invalid method: %s", message.Method))
		return

	}
}

// UTILS
func buildOutboundMessage(params interface{}, method string) (*api.OutboundMessage, error) {
	// Marshal
	marshalledParams, err := json.Marshal(params)
	if err != nil {
		return nil, err
	}
	// Build API call
	api_call := &api.OutboundMessage{
		Method: method,
		Params: marshalledParams,
	}
	return api_call, nil
}

func buildOutboundMessageBytes(params interface{}, method string) ([]byte, error) {
	// Build OutboundMessage from params
	response, err := buildOutboundMessage(params, method)
	if err != nil {
		return nil, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return nil, err
	}
	return responseBytes, nil
}

// HANDLER
func (h *APIHandler) handleEcho(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Build OutboundMessage from params
	responseBytes, err := buildOutboundMessageBytes(params, "echo")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send to connection
	h.server.WebSocketSendConnection(connectionID, responseBytes)
	// Return response
	h.SetSuccessResponseWithMessage(w, responseBytes)
}

func (h *APIHandler) HandleGetBundle(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestUserBundle{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get the bundle
	bundle, ok, err := h.server.GetClientBundle(request.UserID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting client bundle: %v", err))
		return
	}
	// Build bundle response
	response, err := buildOutboundMessageBytes(&api.ResponseUserBundle{
		Success: ok,
		Bundle:  bundle,
	}, "get_bundle")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send response
	h.server.WebSocketSendConnection(connectionID, response)
	// Notify user if OTP is running low
	count, err := h.server.GetRemainingOTPCount(request.UserID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting remaining OTP count for user: %v", err))
		return
	}
	if count < 3 {
		// Notify user
		// Send notification
		notifyLow, err := buildOutboundMessageBytes(&api.NotifyLowOTP{}, "notify_low_otp")
		if err != nil {
			h.SetErrorResponse(w, fmt.Sprintf("Error building low OTP notification: %v", err))
			return
		}
		// Send notification
		h.server.WebSocketSendUser(request.UserID, notifyLow)
	}

	// API response
	h.SetSuccessResponse(w)
}

func (h *APIHandler) HandleUploadBundle(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestUploadBundle{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get connection entry
	entry, err := h.server.getConnectionEntry(connectionID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting connection entry: %v", err))
		return
	}
	// Upload the bundle
	err = h.server.RegisterClient(entry.Username, request.Bundle)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error uploading client bundle: %v", err))
	}
	// Build response
	responseBytes, err := buildOutboundMessageBytes(&api.ResponseUploadBundle{
		Success: true,
	}, "upload_bundle")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send response
	h.server.WebSocketSendConnection(connectionID, responseBytes)

	// API response
	h.SetSuccessResponse(w)
}

func (h *APIHandler) HandleSendMessage(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestSendMsg{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get connection entry
	entry, err := h.server.getConnectionEntry(connectionID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting connection entry: %v", err))
		return
	}
	// Send the message
	ok := h.server.SendMessage(request.RecipientID, entry.Username, request.MessageData)
	// Build response
	responseBytes, err := buildOutboundMessageBytes(&api.ResponseSendMsg{
		Success: ok,
	}, "send_message")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send response
	h.server.WebSocketSendConnection(connectionID, responseBytes)
	// Notify recipient
	notifyRecipient, err := buildOutboundMessageBytes(&api.NotifyNewMessage{
		SenderID: entry.Username,
	}, "notify_new_message")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building notification: %v", err))
		return
	}
	// Send notification
	h.server.WebSocketSendUser(request.RecipientID, notifyRecipient)

	// API response
	h.SetSuccessResponse(w)
}

func (h *APIHandler) HandleReceiveMessage(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestReceiveMsg{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get connection entry
	entry, err := h.server.getConnectionEntry(connectionID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting connection entry: %v", err))
		return
	}
	// Get the message
	messageData, ok, err := h.server.GetMessage(entry.Username)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting message for user: %v", err))
		return
	}
	// Build response
	responseBytes, err := buildOutboundMessageBytes(&api.ResponseReceiveMsg{
		Success:     ok,
		SenderID:    messageData.SenderID,
		MessageData: messageData.Message,
	}, "receive_message")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send response
	h.server.WebSocketSendConnection(connectionID, responseBytes)
	// API response
	h.SetSuccessResponse(w)
}

func (h *APIHandler) HandleStatus(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestUserStatus{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get connection entry
	entry, err := h.server.getConnectionEntry(connectionID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting connection entry: %v", err))
		return
	}
	// Check if user is registered
	registered, err := h.server.IsClientRegistered(entry.Username)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error checking if user is registered: %v", err))
		return
	}
	// Build response
	responseBytes, err := buildOutboundMessageBytes(&api.ResponseUserStatus{
		Success: registered,
	}, "status")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
	// Send response
	h.server.WebSocketSendConnection(connectionID, responseBytes)
	// Notify user if OTP is running low
	count, err := h.server.GetRemainingOTPCount(entry.Username)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting remaining OTP count for user: %v", err))
		return
	}
	if count < 3 {
		// Notify user
		// Send notification
		notifyLow, err := buildOutboundMessageBytes(&api.NotifyLowOTP{}, "notify_low_otp")
		if err != nil {
			h.SetErrorResponse(w, fmt.Sprintf("Error building low OTP notification: %v", err))
			return
		}
		// Send notification
		h.server.WebSocketSendUser(entry.Username, notifyLow)
	}

	// API response
	h.SetSuccessResponse(w)
}

func (h *APIHandler) HandleUploadNewOTPs(w http.ResponseWriter, params json.RawMessage, connectionID string) {
	// Parse request body (JSON)
	request := &api.RequestUploadOTPs{}
	err := json.Unmarshal(params, request)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error parsing request: %v", err))
		return
	}
	// Get connection entry
	entry, err := h.server.getConnectionEntry(connectionID)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error getting connection entry: %v", err))
		return
	}
	// Upload the OTPs
	err = h.server.ExpandOTPSet(entry.Username, request.OTPs)
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error uploading OTPs: %v", err))
		return
	}
	// API response
	h.SetSuccessResponse(w)
}
