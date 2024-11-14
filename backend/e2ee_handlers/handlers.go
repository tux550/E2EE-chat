package main

import (
	"encoding/json"
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
		/*
				case "get_bundle":
					return h.HandleGetBundle(request, message.Params)
				case "upload_bundle":
					return h.HandleUploadBundle(request, message.Params)

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

// HANDLER
func (h *APIHandler) handleEcho(w http.ResponseWriter, params json.RawMessage) {
	// Build OutboundMessage from params
	response, err := buildOutboundMessage(params, "echo")
	if err != nil {
		log.Println("Error building response:", err)
		http.Error(w, "Error building response", http.StatusInternalServerError)
		return
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		log.Println("Error marshalling response:", err)
		http.Error(w, "Error marshalling response", http.StatusInternalServerError)
		return
	}
	// Return response
	w.WriteHeader(http.StatusOK)
	w.Write(responseBytes)
}

/*
func (h *APIHandler) HandleGetBundle(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestUserBundle{}
	err := json.Unmarshal(params, request)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Get the bundle
	bundle, ok, err := h.server.GetClientBundle(request.UserID)
	if err != nil {
		log.Println("Error getting client bundle:", err)
		return events.APIGatewayProxyResponse{}, err
	}

	// TODO:
	// NOTIFY RECIPIENT IF OTP IS RUNNING LOW

	// Build response
	response, err := buildOutboundMessage(&api.ResponseUserBundle{
		Success: ok,
		Bundle:  bundle,
	}, "get_bundle")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Return response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBytes),
	}, nil

}

func (h *APIHandler) HandleUploadBundle(req events.APIGatewayWebsocketProxyRequest, params json.RawMessage) (events.APIGatewayProxyResponse, error) {
	// Parse request body (JSON)
	request := &api.RequestUploadBundle{}
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
	// Upload the bundle
	err = h.server.RegisterClient(entry.Username, request.Bundle)
	if err != nil {
		log.Println("Error uploading client bundle:", err)
		return events.APIGatewayProxyResponse{}, err
	}
	// Build response
	response, err := buildOutboundMessage(&api.ResponseUploadBundle{
		Success: true,
	}, "upload_bundle")
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Marshal response
	responseBytes, err := json.Marshal(response)
	if err != nil {
		return events.APIGatewayProxyResponse{}, err
	}
	// Return response
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(responseBytes),
	}, nil
}
*/
