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
	// Parse request
	request := HTTPNormalizedRequest{}
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Error decoding request:", err)
		http.Error(w, "Error decoding request", http.StatusBadRequest)
		return
	}
	// Parse post body
	message := &api.InboundMessage{}
	err = json.Unmarshal([]byte(request.Body), message)
	if err != nil {
		log.Println("Error parsing message:", err)
		http.Error(w, "Error parsing message", http.StatusBadRequest)
		return
	}
	// Handle message
	switch message.Method {
	default:
		// ECHO
		h.handleEcho(w, message.Params) // Handlers must also recieve connectionID to be able to reply
		return
	case "get_bundle":
		h.HandleGetBundle(w, message.Params, request.ConnectionID)
		return
	case "upload_bundle":
		h.HandleUploadBundle(w, message.Params, request.ConnectionID)
		return
		/*
				case "send_message":
					return h.HandleSendMessage(request, message.Params)
				case "receive_message":
					return h.HandleReceiveMessage(request, message.Params)
				case "status":
					return h.HandleStatus(request, message.Params)
				case "upload_new_otps":
					return h.HandleUploadNewOTPs(request, message.Params)

			default:
				return events.APIGatewayProxyResponse{}, nil
		*/
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
func (h *APIHandler) handleEcho(w http.ResponseWriter, params json.RawMessage) {
	// Build OutboundMessage from params
	responseBytes, err := buildOutboundMessageBytes(params, "echo")
	if err != nil {
		h.SetErrorResponse(w, fmt.Sprintf("Error building response: %v", err))
		return
	}
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

/*
func (h *APIHandler) HandleSendMessage(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestSendMsg{}
	err := json.Unmarshal(params, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Get connection entry
	// FIXME: Dunno if this is the right way to get the connection entry
	entry, err := h.getConnectionEntry(req)
	if err != nil {
		log.Println("Error getting connection entry:", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// Send the message
	ok := h.server.SendMessage(request.RecipientID, entry.Username, request.MessageData)

	// Build response
	response, err := buildOutboundMessage(&api.ResponseSendMsg{
		Success: ok,
	}, "send_message")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error building response",
		}, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling response",
		}, err
	}

	notification, err := buildOutboundMessage(&api.NotifyNewMessage{
		SenderID: entry.Username,
	}, "notify_new_message")
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error building notification",
		}, err
	}

	notificationBytes, err := json.Marshal(notification)
	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusInternalServerError,
			Body:       "Error marshalling notification",
		}, err
	}

	h.server.SendNotificationToUser(request.RecipientID, notificationBytes)

	// Return response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBytes),
	}, nil
}

func (h *APIHandler) HandleReceiveMessage(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestReceiveMsg{}
	err := json.Unmarshal(params, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Get connection entry
	entry, err := h.getConnectionEntry(req)
	if err != nil {
		log.Println("Error getting connection entry:", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// Get the message
	messageData, ok, err := h.server.GetMessage(entry.Username)
	if err != nil {
		log.Println("Error getting message for user:", entry.Username)
		return events.APIGatewayProxyResponse{}, err
	}

	if !ok {
		fmt.Println("User", entry.Username, "requested message but none available")
		// Send response
		response, err := buildOutboundMessage(&api.ResponseReceiveMsg{
			Success: false,
		}, "receive_message")
		if err != nil {
			fmt.Println("Error marshalling response to receive_message")
			return events.APIGatewayProxyResponse{}, err
		}
		// Marshal response
		responseBytes, err := json.Marshal(response)
		if err != nil {
			fmt.Println("Error marshalling fail response to receive_message")
			return events.APIGatewayProxyResponse{}, err
		}
		// Return response
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(responseBytes),
		}, nil
	}

	fmt.Println("User", entry.Username, "received message from user")

	// Send response
	response, err := buildOutboundMessage(&api.ResponseReceiveMsg{
		Success:     true,
		SenderID:    messageData.SenderID,
		MessageData: messageData.Message,
	}, "receive_message")
	if err != nil {
		fmt.Println("Error marshalling response to receive_message")
		return events.APIGatewayProxyResponse{}, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling response to receive_message")
		return events.APIGatewayProxyResponse{}, err
	}
	// Return response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,

		Body: string(responseBytes),
	}, nil
}

func (h *APIHandler) HandleStatus(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestUserStatus{}
	err := json.Unmarshal(params, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Get connection entry
	// FIXME: call the appropriate function
	entry, err := h.getConnectionEntry(req)
	if err != nil {
		log.Println("Error getting connection entry:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	registered, err := h.server.IsClientRegistered(entry.Username)

	if err != nil {
		log.Println("Error checking if user", entry.Username, "is registered")
		return events.APIGatewayProxyResponse{}, err
	}

	fmt.Println("User", entry.Username, "checked if self is registered")

	// Build response
	response, err := buildOutboundMessage(&api.ResponseUserStatus{
		Success: registered,
	}, "status")
	if err != nil {
		fmt.Println("Error marshalling response to status", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		fmt.Println("Error marshalling success response to status", err)
		return events.APIGatewayProxyResponse{}, err
	}

	if !registered {
		fmt.Println("User", entry.Username, "is not registered")
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusUnauthorized,
			Body:       string(responseBytes),
		}, nil
	}

	count, err := h.server.GetRemainingOTPCount(entry.Username)
	if err != nil {
		fmt.Println("Error getting remaining OTP count for user", entry.Username)

		return events.APIGatewayProxyResponse{}, err
	}

	if count < 3 {
		// Notify user
		fmt.Println("Notifying user", entry.Username, "that OTP is running low")
		// Send notification
		notificationBytes := getLowOTPNotification()
		return events.APIGatewayProxyResponse{
			StatusCode: http.StatusOK,
			Body:       string(notificationBytes),
		}, nil
	}
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       "Count is not less than 3",
	}, nil
}

func (h *APIHandler) HandleUploadNewOTPs(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestUploadOTPs{}
	err := json.Unmarshal(params, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}

	// Get connection entry
	entry, err := h.getConnectionEntry(req)
	if err != nil {
		log.Println("Error getting connection entry:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	h.server.ExpandOTPSet(entry.Username, request.OTPs)
	fmt.Println("User", entry.Username, "uploaded #", len(request.OTPs), "new OTPs")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       "Success",
	}, nil
}
*/
