package obs_studio

import (
	"encoding/json"
	"reflect"
	"sync"
)

// EventHandler is a function type for handling events
type EventHandler func(eventData interface{})

// EventEmitter manages event listeners and emission
type EventEmitter struct {
	listeners map[string][]EventHandler
	mutex     sync.RWMutex
}

// NewEventEmitter creates a new event emitter
func NewEventEmitter() *EventEmitter {
	return &EventEmitter{
		listeners: make(map[string][]EventHandler),
	}
}

// On registers an event handler for a specific event type
func (e *EventEmitter) On(eventType string, handler EventHandler) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	e.listeners[eventType] = append(e.listeners[eventType], handler)
}

// Off removes an event handler for a specific event type
func (e *EventEmitter) Off(eventType string, handler EventHandler) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	handlers := e.listeners[eventType]
	for i, h := range handlers {
		if reflect.ValueOf(h).Pointer() == reflect.ValueOf(handler).Pointer() {
			e.listeners[eventType] = append(handlers[:i], handlers[i+1:]...)
			break
		}
	}
}

// Emit emits an event to all registered handlers
func (e *EventEmitter) Emit(eventType string, eventData interface{}) {
	e.mutex.RLock()
	handlers := make([]EventHandler, len(e.listeners[eventType]))
	copy(handlers, e.listeners[eventType])
	e.mutex.RUnlock()

	for _, handler := range handlers {
		go handler(eventData)
	}
}

// RemoveAllListeners removes all listeners for a specific event type
func (e *EventEmitter) RemoveAllListeners(eventType string) {
	e.mutex.Lock()
	defer e.mutex.Unlock()

	delete(e.listeners, eventType)
}

// ListenerCount returns the number of listeners for a specific event type
func (e *EventEmitter) ListenerCount(eventType string) int {
	e.mutex.RLock()
	defer e.mutex.RUnlock()

	return len(e.listeners[eventType])
}

// EventProcessor handles processing of incoming events
type EventProcessor struct {
	emitter *EventEmitter
}

// NewEventProcessor creates a new event processor
func NewEventProcessor() *EventProcessor {
	return &EventProcessor{
		emitter: NewEventEmitter(),
	}
}

// ProcessEvent processes an incoming event message
func (ep *EventProcessor) ProcessEvent(message *WebSocketMessage) error {
	eventData := &EventData{}

	// Convert message data to EventData
	dataBytes, err := json.Marshal(message.D)
	if err != nil {
		return err
	}

	if err := json.Unmarshal(dataBytes, eventData); err != nil {
		return err
	}

	// Emit the event
	ep.emitter.Emit(eventData.EventType, eventData.EventData)

	// Also emit generic "event" for all events
	ep.emitter.Emit("event", eventData)

	return nil
}

// On registers an event handler
func (ep *EventProcessor) On(eventType string, handler EventHandler) {
	ep.emitter.On(eventType, handler)
}

// Off removes an event handler
func (ep *EventProcessor) Off(eventType string, handler EventHandler) {
	ep.emitter.Off(eventType, handler)
}

// Emit emits an event
func (ep *EventProcessor) Emit(eventType string, eventData interface{}) {
	ep.emitter.Emit(eventType, eventData)
}

// Common event data structures

// SceneChangedEventData represents scene change event data
type SceneChangedEventData struct {
	SceneName string `json:"sceneName"`
}

// StreamStateChangedEventData represents stream state change event data
type StreamStateChangedEventData struct {
	OutputActive bool   `json:"outputActive"`
	OutputState  string `json:"outputState"`
}

// RecordStateChangedEventData represents record state change event data
type RecordStateChangedEventData struct {
	OutputActive bool   `json:"outputActive"`
	OutputState  string `json:"outputState"`
	OutputPath   string `json:"outputPath,omitempty"`
}

// SceneCreatedEventData represents scene created event data
type SceneCreatedEventData struct {
	SceneName  string `json:"sceneName"`
	IsGroup    bool   `json:"isGroup"`
	SceneIndex int    `json:"sceneIndex"`
}

// SceneRemovedEventData represents scene removed event data
type SceneRemovedEventData struct {
	SceneName string `json:"sceneName"`
	IsGroup   bool   `json:"isGroup"`
}

// ConnectionEventData represents connection-related events
type ConnectionEventData struct {
	ConnectionStatus string `json:"connectionStatus"`
	Message          string `json:"message,omitempty"`
}

// Built-in event types for common OBS events
const (
	// Connection events
	EventConnectionOpened = "ConnectionOpened"
	EventConnectionClosed = "ConnectionClosed"
	EventConnectionError  = "ConnectionError"
	EventIdentified       = "Identified"
	EventReidentified     = "Reidentified"

	// Scene events
	EventCurrentProgramSceneChanged = "CurrentProgramSceneChanged"
	EventCurrentPreviewSceneChanged = "CurrentPreviewSceneChanged"
	EventSceneCreated               = "SceneCreated"
	EventSceneRemoved               = "SceneRemoved"
	EventSceneNameChanged           = "SceneNameChanged"

	// Stream/Record events
	EventStreamStateChanged = "StreamStateChanged"
	EventRecordStateChanged = "RecordStateChanged"

	// General events
	EventExitStarted            = "ExitStarted"
	EventStudioModeStateChanged = "StudioModeStateChanged"

	// Generic event for all events
	EventGeneric = "event"
)
