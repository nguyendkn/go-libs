package obs_studio

import (
	"time"
)

// WebSocket Operation Codes
type OpCode int

const (
	OpCodeHello                OpCode = 0
	OpCodeIdentify             OpCode = 1
	OpCodeIdentified           OpCode = 2
	OpCodeReidentify           OpCode = 3
	OpCodeEvent                OpCode = 5
	OpCodeRequest              OpCode = 6
	OpCodeRequestResponse      OpCode = 7
	OpCodeRequestBatch         OpCode = 8
	OpCodeRequestBatchResponse OpCode = 9
)

// Close Codes
const (
	CloseCodeUnknownReason         = 4000
	CloseCodeMessageDecodeError    = 4002
	CloseCodeMissingDataField      = 4003
	CloseCodeInvalidDataFieldType  = 4004
	CloseCodeInvalidDataFieldValue = 4005
	CloseCodeUnknownOpCode         = 4006
	CloseCodeNotIdentified         = 4007
	CloseCodeAlreadyIdentified     = 4008
	CloseCodeAuthenticationFailed  = 4009
	CloseCodeUnsupportedRpcVersion = 4010
	CloseCodeSessionInvalidated    = 4011
	CloseCodeUnsupportedFeature    = 4012
)

// WebSocket Message represents the basic message structure
type WebSocketMessage struct {
	Op OpCode      `json:"op"`
	D  interface{} `json:"d"`
}

// Hello message data
type HelloData struct {
	ObsWebSocketVersion string         `json:"obsWebSocketVersion"`
	RpcVersion          int            `json:"rpcVersion"`
	Authentication      *AuthChallenge `json:"authentication,omitempty"`
}

// AuthChallenge contains authentication challenge data
type AuthChallenge struct {
	Challenge string `json:"challenge"`
	Salt      string `json:"salt"`
}

// Identify message data
type IdentifyData struct {
	RpcVersion         int    `json:"rpcVersion"`
	Authentication     string `json:"authentication,omitempty"`
	EventSubscriptions int    `json:"eventSubscriptions,omitempty"`
}

// Identified message data
type IdentifiedData struct {
	NegotiatedRpcVersion int `json:"negotiatedRpcVersion"`
}

// Request data structure
type RequestData struct {
	RequestType string      `json:"requestType"`
	RequestId   string      `json:"requestId"`
	RequestData interface{} `json:"requestData,omitempty"`
}

// Request response data
type RequestResponseData struct {
	RequestType   string        `json:"requestType"`
	RequestId     string        `json:"requestId"`
	RequestStatus RequestStatus `json:"requestStatus"`
	ResponseData  interface{}   `json:"responseData,omitempty"`
}

// Request status
type RequestStatus struct {
	Result  bool   `json:"result"`
	Code    int    `json:"code"`
	Comment string `json:"comment,omitempty"`
}

// Event data structure
type EventData struct {
	EventType   string      `json:"eventType"`
	EventIntent int         `json:"eventIntent"`
	EventData   interface{} `json:"eventData,omitempty"`
}

// Batch request data
type BatchRequestData struct {
	RequestId     string        `json:"requestId"`
	HaltOnFailure bool          `json:"haltOnFailure,omitempty"`
	ExecutionType int           `json:"executionType,omitempty"`
	Requests      []RequestData `json:"requests"`
}

// Batch response data
type BatchResponseData struct {
	RequestId string                `json:"requestId"`
	Results   []RequestResponseData `json:"results"`
}

// Connection configuration
type ConnectionConfig struct {
	Address            string
	Password           string
	RpcVersion         int
	EventSubscriptions int
	ConnectTimeout     time.Duration
	RequestTimeout     time.Duration
}

// Scene data structure
type Scene struct {
	SceneIndex int    `json:"sceneIndex"`
	SceneName  string `json:"sceneName"`
}

// Scene list response
type SceneListResponse struct {
	CurrentProgramSceneName string  `json:"currentProgramSceneName"`
	CurrentPreviewSceneName string  `json:"currentPreviewSceneName"`
	Scenes                  []Scene `json:"scenes"`
}

// Source data structure
type Source struct {
	SourceName string `json:"sourceName"`
	SourceType string `json:"sourceType"`
	SourceKind string `json:"sourceKind"`
}

// Event subscription flags
const (
	EventSubscriptionGeneral     = 1 << 0
	EventSubscriptionConfig      = 1 << 1
	EventSubscriptionScenes      = 1 << 2
	EventSubscriptionInputs      = 1 << 3
	EventSubscriptionTransitions = 1 << 4
	EventSubscriptionFilters     = 1 << 5
	EventSubscriptionOutputs     = 1 << 6
	EventSubscriptionSceneItems  = 1 << 7
	EventSubscriptionMediaInputs = 1 << 8
	EventSubscriptionVendors     = 1 << 9
	EventSubscriptionUi          = 1 << 10
	EventSubscriptionAll         = (1 << 11) - 1
)

// Common request types
const (
	RequestTypeGetVersion             = "GetVersion"
	RequestTypeGetStats               = "GetStats"
	RequestTypeGetSceneList           = "GetSceneList"
	RequestTypeSetCurrentProgramScene = "SetCurrentProgramScene"
	RequestTypeSetCurrentPreviewScene = "SetCurrentPreviewScene"
	RequestTypeCreateScene            = "CreateScene"
	RequestTypeRemoveScene            = "RemoveScene"
	RequestTypeGetSourcesList         = "GetSourcesList"
	RequestTypeGetInputList           = "GetInputList"
	RequestTypeStartStream            = "StartStream"
	RequestTypeStopStream             = "StopStream"
	RequestTypeStartRecord            = "StartRecord"
	RequestTypeStopRecord             = "StopRecord"
)

// Common event types
const (
	EventTypeExitStarted                = "ExitStarted"
	EventTypeCurrentProgramSceneChanged = "CurrentProgramSceneChanged"
	EventTypeCurrentPreviewSceneChanged = "CurrentPreviewSceneChanged"
	EventTypeSceneCreated               = "SceneCreated"
	EventTypeSceneRemoved               = "SceneRemoved"
	EventTypeStreamStateChanged         = "StreamStateChanged"
	EventTypeRecordStateChanged         = "RecordStateChanged"
)
