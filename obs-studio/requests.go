package obs_studio

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

// RequestManager handles request/response operations
type RequestManager struct {
	pendingRequests map[string]*PendingRequest
	mutex           sync.RWMutex
	timeout         time.Duration
}

// PendingRequest represents a request waiting for response
type PendingRequest struct {
	RequestID   string
	RequestType string
	Response    chan *RequestResponseData
	Error       chan error
	Timeout     *time.Timer
	Context     context.Context
	Cancel      context.CancelFunc
}

// NewRequestManager creates a new request manager
func NewRequestManager(timeout time.Duration) *RequestManager {
	return &RequestManager{
		pendingRequests: make(map[string]*PendingRequest),
		timeout:         timeout,
	}
}

// CreateRequest creates a new request with auto-generated ID
func (rm *RequestManager) CreateRequest(requestType string, requestData interface{}) *RequestData {
	requestID := uuid.New().String()

	return &RequestData{
		RequestType: requestType,
		RequestId:   requestID,
		RequestData: requestData,
	}
}

// SendRequest sends a request and waits for response
func (rm *RequestManager) SendRequest(ctx context.Context, request *RequestData, sendFunc func(*WebSocketMessage) error) (*RequestResponseData, error) {
	// Create pending request
	pending := &PendingRequest{
		RequestID:   request.RequestId,
		RequestType: request.RequestType,
		Response:    make(chan *RequestResponseData, 1),
		Error:       make(chan error, 1),
	}

	// Set up context with timeout
	pending.Context, pending.Cancel = context.WithTimeout(ctx, rm.timeout)

	// Register pending request
	rm.mutex.Lock()
	rm.pendingRequests[request.RequestId] = pending
	rm.mutex.Unlock()

	// Create WebSocket message
	message := &WebSocketMessage{
		Op: OpCodeRequest,
		D:  request,
	}

	// Send the request
	if err := sendFunc(message); err != nil {
		rm.removePendingRequest(request.RequestId)
		return nil, fmt.Errorf("failed to send request: %w", err)
	}

	// Wait for response or timeout
	select {
	case response := <-pending.Response:
		rm.removePendingRequest(request.RequestId)
		return response, nil
	case err := <-pending.Error:
		rm.removePendingRequest(request.RequestId)
		return nil, err
	case <-pending.Context.Done():
		rm.removePendingRequest(request.RequestId)
		return nil, fmt.Errorf("request timeout: %s", request.RequestType)
	}
}

// SendBatchRequest sends a batch of requests
func (rm *RequestManager) SendBatchRequest(ctx context.Context, requests []RequestData, haltOnFailure bool, sendFunc func(*WebSocketMessage) error) (*BatchResponseData, error) {
	batchID := uuid.New().String()

	batchRequest := &BatchRequestData{
		RequestId:     batchID,
		HaltOnFailure: haltOnFailure,
		Requests:      requests,
	}

	// Create pending request for batch
	pending := &PendingRequest{
		RequestID:   batchID,
		RequestType: "BatchRequest",
		Response:    make(chan *RequestResponseData, 1),
		Error:       make(chan error, 1),
	}

	pending.Context, pending.Cancel = context.WithTimeout(ctx, rm.timeout)

	// Register pending request
	rm.mutex.Lock()
	rm.pendingRequests[batchID] = pending
	rm.mutex.Unlock()

	// Create WebSocket message
	message := &WebSocketMessage{
		Op: OpCodeRequestBatch,
		D:  batchRequest,
	}

	// Send the batch request
	if err := sendFunc(message); err != nil {
		rm.removePendingRequest(batchID)
		return nil, fmt.Errorf("failed to send batch request: %w", err)
	}

	// Wait for response or timeout
	select {
	case response := <-pending.Response:
		rm.removePendingRequest(batchID)

		// Convert response to BatchResponseData
		batchResponse := &BatchResponseData{}
		responseBytes, err := json.Marshal(response)
		if err != nil {
			return nil, fmt.Errorf("failed to marshal batch response: %w", err)
		}

		if err := json.Unmarshal(responseBytes, batchResponse); err != nil {
			return nil, fmt.Errorf("failed to unmarshal batch response: %w", err)
		}

		return batchResponse, nil
	case err := <-pending.Error:
		rm.removePendingRequest(batchID)
		return nil, err
	case <-pending.Context.Done():
		rm.removePendingRequest(batchID)
		return nil, fmt.Errorf("batch request timeout")
	}
}

// ProcessResponse processes incoming response messages
func (rm *RequestManager) ProcessResponse(message *WebSocketMessage) error {
	var response *RequestResponseData

	// Handle different response types
	switch message.Op {
	case OpCodeRequestResponse:
		response = &RequestResponseData{}
		dataBytes, err := json.Marshal(message.D)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(dataBytes, response); err != nil {
			return err
		}

	case OpCodeRequestBatchResponse:
		batchResponse := &BatchResponseData{}
		dataBytes, err := json.Marshal(message.D)
		if err != nil {
			return err
		}

		if err := json.Unmarshal(dataBytes, batchResponse); err != nil {
			return err
		}

		// Convert batch response to individual response for processing
		response = &RequestResponseData{
			RequestId: batchResponse.RequestId,
			RequestStatus: RequestStatus{
				Result: true,
				Code:   100,
			},
			ResponseData: batchResponse,
		}

	default:
		return fmt.Errorf("invalid response op code: %d", message.Op)
	}

	// Find and notify pending request
	rm.mutex.RLock()
	pending, exists := rm.pendingRequests[response.RequestId]
	rm.mutex.RUnlock()

	if !exists {
		return fmt.Errorf("no pending request found for ID: %s", response.RequestId)
	}

	// Check if request was successful
	if !response.RequestStatus.Result {
		pending.Error <- fmt.Errorf("request failed: %s (code: %d)",
			response.RequestStatus.Comment, response.RequestStatus.Code)
		return nil
	}

	// Send response
	select {
	case pending.Response <- response:
	default:
		// Channel might be closed or full
	}

	return nil
}

// removePendingRequest removes a pending request from the map
func (rm *RequestManager) removePendingRequest(requestID string) {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	if pending, exists := rm.pendingRequests[requestID]; exists {
		if pending.Cancel != nil {
			pending.Cancel()
		}
		delete(rm.pendingRequests, requestID)
	}
}

// CancelAllRequests cancels all pending requests
func (rm *RequestManager) CancelAllRequests() {
	rm.mutex.Lock()
	defer rm.mutex.Unlock()

	for requestID, pending := range rm.pendingRequests {
		if pending.Cancel != nil {
			pending.Cancel()
		}
		delete(rm.pendingRequests, requestID)
	}
}

// GetPendingRequestCount returns the number of pending requests
func (rm *RequestManager) GetPendingRequestCount() int {
	rm.mutex.RLock()
	defer rm.mutex.RUnlock()

	return len(rm.pendingRequests)
}

// Common request data structures

// SetCurrentSceneRequestData represents data for setting current scene
type SetCurrentSceneRequestData struct {
	SceneName string `json:"sceneName"`
}

// CreateSceneRequestData represents data for creating a scene
type CreateSceneRequestData struct {
	SceneName string `json:"sceneName"`
}

// RemoveSceneRequestData represents data for removing a scene
type RemoveSceneRequestData struct {
	SceneName string `json:"sceneName"`
}

// SetStreamSettingsRequestData represents data for setting stream settings
type SetStreamSettingsRequestData struct {
	StreamSettings map[string]interface{} `json:"streamSettings"`
}

// SetRecordDirectoryRequestData represents data for setting record directory
type SetRecordDirectoryRequestData struct {
	RecordDirectory string `json:"recordDirectory"`
}

// Helper methods for common requests

// CreateGetVersionRequest creates a GetVersion request
func (rm *RequestManager) CreateGetVersionRequest() *RequestData {
	return rm.CreateRequest(RequestTypeGetVersion, nil)
}

// CreateGetSceneListRequest creates a GetSceneList request
func (rm *RequestManager) CreateGetSceneListRequest() *RequestData {
	return rm.CreateRequest(RequestTypeGetSceneList, nil)
}

// CreateSetCurrentSceneRequest creates a SetCurrentProgramScene request
func (rm *RequestManager) CreateSetCurrentSceneRequest(sceneName string) *RequestData {
	return rm.CreateRequest(RequestTypeSetCurrentProgramScene, &SetCurrentSceneRequestData{
		SceneName: sceneName,
	})
}

// CreateStartStreamRequest creates a StartStream request
func (rm *RequestManager) CreateStartStreamRequest() *RequestData {
	return rm.CreateRequest(RequestTypeStartStream, nil)
}

// CreateStopStreamRequest creates a StopStream request
func (rm *RequestManager) CreateStopStreamRequest() *RequestData {
	return rm.CreateRequest(RequestTypeStopStream, nil)
}

// CreateStartRecordRequest creates a StartRecord request
func (rm *RequestManager) CreateStartRecordRequest() *RequestData {
	return rm.CreateRequest(RequestTypeStartRecord, nil)
}

// CreateStopRecordRequest creates a StopRecord request
func (rm *RequestManager) CreateStopRecordRequest() *RequestData {
	return rm.CreateRequest(RequestTypeStopRecord, nil)
}
