package httpclient

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// responseProcessor xử lý response
type responseProcessor struct {
	processors map[string]func(*Response) error
}

// NewResponseProcessor tạo response processor mới
func NewResponseProcessor() *responseProcessor {
	rp := &responseProcessor{
		processors: make(map[string]func(*Response) error),
	}

	// Register default processors
	rp.RegisterProcessor("application/json", rp.processJSON)
	rp.RegisterProcessor("application/xml", rp.processXML)
	rp.RegisterProcessor("text/xml", rp.processXML)
	rp.RegisterProcessor("text/plain", rp.processText)
	rp.RegisterProcessor("text/html", rp.processHTML)

	return rp
}

// RegisterProcessor đăng ký processor cho content type
func (rp *responseProcessor) RegisterProcessor(contentType string, processor func(*Response) error) {
	rp.processors[contentType] = processor
}

// Process xử lý response
func (rp *responseProcessor) Process(resp *Response) error {
	contentType := resp.ContentType
	if contentType == "" {
		return nil
	}

	// Extract main content type (remove charset, etc.)
	mainType := strings.Split(contentType, ";")[0]
	mainType = strings.TrimSpace(mainType)

	if processor, exists := rp.processors[mainType]; exists {
		return processor(resp)
	}

	// Default processing - just read body
	return rp.processDefault(resp)
}

// Default processors

func (rp *responseProcessor) processJSON(resp *Response) error {
	// Body is already read in buildResponse
	return nil
}

func (rp *responseProcessor) processXML(resp *Response) error {
	// Body is already read in buildResponse
	return nil
}

func (rp *responseProcessor) processText(resp *Response) error {
	// Body is already read in buildResponse
	return nil
}

func (rp *responseProcessor) processHTML(resp *Response) error {
	// Body is already read in buildResponse
	return nil
}

func (rp *responseProcessor) processDefault(_ *Response) error {
	// Body is already read in buildResponse
	return nil
}

// Response helper methods

// JSON unmarshals response body to JSON
func (r *Response) JSON(v interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}

	return json.Unmarshal(r.Body, v)
}

// XML unmarshals response body to XML
func (r *Response) XML(v interface{}) error {
	if len(r.Body) == 0 {
		return fmt.Errorf("empty response body")
	}

	return xml.Unmarshal(r.Body, v)
}

// String returns response body as string
func (r *Response) String() string {
	return string(r.Body)
}

// Bytes returns response body as bytes
func (r *Response) Bytes() []byte {
	return r.Body
}

// IsSuccess checks if response is successful (2xx)
func (r *Response) IsSuccess() bool {
	return r.StatusCode >= 200 && r.StatusCode < 300
}

// IsRedirect checks if response is redirect (3xx)
func (r *Response) IsRedirect() bool {
	return r.StatusCode >= 300 && r.StatusCode < 400
}

// IsClientError checks if response is client error (4xx)
func (r *Response) IsClientError() bool {
	return r.StatusCode >= 400 && r.StatusCode < 500
}

// IsServerError checks if response is server error (5xx)
func (r *Response) IsServerError() bool {
	return r.StatusCode >= 500 && r.StatusCode < 600
}

// IsError checks if response is error (4xx or 5xx)
func (r *Response) IsError() bool {
	return r.IsClientError() || r.IsServerError()
}

// Header gets header value
func (r *Response) Header(key string) string {
	values := r.Headers[key]
	if len(values) > 0 {
		return values[0]
	}
	return ""
}

// HeaderValues gets all header values
func (r *Response) HeaderValues(key string) []string {
	return r.Headers[key]
}

// HasHeader checks if header exists
func (r *Response) HasHeader(key string) bool {
	_, exists := r.Headers[key]
	return exists
}

// GetContentLength returns content length
func (r *Response) GetContentLength() int64 {
	if r.ContentLength > 0 {
		return r.ContentLength
	}

	// Try to parse from header
	if lengthStr := r.Header("Content-Length"); lengthStr != "" {
		if length, err := strconv.ParseInt(lengthStr, 10, 64); err == nil {
			return length
		}
	}

	// Return body length as fallback
	return int64(len(r.Body))
}

// GetContentType returns content type
func (r *Response) GetContentType() string {
	return r.ContentType
}

// GetCharset returns charset from content type
func (r *Response) GetCharset() string {
	contentType := r.ContentType
	if contentType == "" {
		return ""
	}

	parts := strings.Split(contentType, ";")
	for _, part := range parts {
		part = strings.TrimSpace(part)
		if strings.HasPrefix(part, "charset=") {
			return strings.TrimPrefix(part, "charset=")
		}
	}

	return ""
}

// GetLocation returns Location header for redirects
func (r *Response) GetLocation() string {
	return r.Header("Location")
}

// GetETag returns ETag header
func (r *Response) GetETag() string {
	return r.Header("ETag")
}

// GetLastModified returns Last-Modified header
func (r *Response) GetLastModified() string {
	return r.Header("Last-Modified")
}

// GetCacheControl returns Cache-Control header
func (r *Response) GetCacheControl() string {
	return r.Header("Cache-Control")
}

// GetExpires returns Expires header
func (r *Response) GetExpires() string {
	return r.Header("Expires")
}

// GetServer returns Server header
func (r *Response) GetServer() string {
	return r.Header("Server")
}

// GetDate returns Date header
func (r *Response) GetDate() string {
	return r.Header("Date")
}

// Enhanced buildResponse implementation
func (c *httpClient) buildResponse(req *Request, httpResp *http.Response, duration time.Duration) (*Response, error) {
	resp := &Response{
		StatusCode:    httpResp.StatusCode,
		Status:        httpResp.Status,
		Headers:       make(map[string][]string),
		ContentType:   httpResp.Header.Get("Content-Type"),
		ContentLength: httpResp.ContentLength,
		Request:       req,
		Duration:      duration,
		Attempts:      req.attempt,
		TotalDuration: time.Since(req.startTime),
		Raw:           httpResp,
		Metadata:      make(map[string]interface{}),
	}

	// Copy headers
	for key, values := range httpResp.Header {
		resp.Headers[key] = make([]string, len(values))
		copy(resp.Headers[key], values)
	}

	// Read body
	if httpResp.Body != nil {
		bodyBytes, err := io.ReadAll(httpResp.Body)
		if err != nil {
			return resp, &HTTPError{
				Code:     1200,
				Message:  fmt.Sprintf("failed to read response body: %v", err),
				Type:     "body",
				Response: resp,
			}
		}
		resp.Body = bodyBytes

		// Create new reader for BodyReader
		resp.BodyReader = io.NopCloser(bytes.NewReader(bodyBytes))
	}

	// Process response based on content type
	processor := NewResponseProcessor()
	if err := processor.Process(resp); err != nil {
		return resp, err
	}

	// Check for HTTP errors
	if resp.IsError() {
		httpErr := &HTTPError{
			Code:       1201,
			Message:    fmt.Sprintf("HTTP %d: %s", resp.StatusCode, resp.Status),
			Type:       "http",
			StatusCode: resp.StatusCode,
			Response:   resp,
		}

		// Try to extract error message from response body
		if len(resp.Body) > 0 {
			var errorData map[string]interface{}
			if err := json.Unmarshal(resp.Body, &errorData); err == nil {
				if msg, ok := errorData["message"].(string); ok {
					httpErr.Message = msg
				} else if msg, ok := errorData["error"].(string); ok {
					httpErr.Message = msg
				}
			}
		}

		return resp, httpErr
	}

	return resp, nil
}

// Response validation helpers

// ValidateStatus validates response status code
func (r *Response) ValidateStatus(expectedCodes ...int) error {
	for _, code := range expectedCodes {
		if r.StatusCode == code {
			return nil
		}
	}

	return &HTTPError{
		Code:       1202,
		Message:    fmt.Sprintf("unexpected status code %d", r.StatusCode),
		Type:       "status",
		StatusCode: r.StatusCode,
		Response:   r,
	}
}

// ValidateContentType validates response content type
func (r *Response) ValidateContentType(expectedTypes ...string) error {
	contentType := strings.Split(r.ContentType, ";")[0]
	contentType = strings.TrimSpace(contentType)

	for _, expectedType := range expectedTypes {
		if contentType == expectedType {
			return nil
		}
	}

	return &HTTPError{
		Code:     1203,
		Message:  fmt.Sprintf("unexpected content type %s", r.ContentType),
		Type:     "content_type",
		Response: r,
	}
}

// ValidateHeader validates response header
func (r *Response) ValidateHeader(key, expectedValue string) error {
	actualValue := r.Header(key)
	if actualValue != expectedValue {
		return &HTTPError{
			Code:     1204,
			Message:  fmt.Sprintf("unexpected header %s: expected %s, got %s", key, expectedValue, actualValue),
			Type:     "header",
			Response: r,
		}
	}

	return nil
}

// ValidateJSON validates response as JSON and unmarshals
func (r *Response) ValidateJSON(v interface{}) error {
	if err := r.ValidateContentType("application/json"); err != nil {
		return err
	}

	return r.JSON(v)
}

// ValidateXML validates response as XML and unmarshals
func (r *Response) ValidateXML(v interface{}) error {
	if err := r.ValidateContentType("application/xml", "text/xml"); err != nil {
		return err
	}

	return r.XML(v)
}
