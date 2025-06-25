package httpclient

import (
	"context"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io"
	"net/url"
	"reflect"
	"strings"
	"time"
)

// requestBuilder implements the RequestBuilder interface
type requestBuilder struct {
	client  *httpClient
	request *Request
}

// NewRequestBuilder tạo một RequestBuilder mới
func NewRequestBuilder(client *httpClient, method HTTPMethod, url string) RequestBuilder {
	rb := &requestBuilder{
		client: client,
		request: &Request{
			Method:      method,
			URL:         url,
			Headers:     make(map[string]string),
			QueryParams: make(map[string]string),
			Metadata:    make(map[string]interface{}),
			Context:     context.Background(),
		},
	}
	return rb
}

// URL and path methods
func (rb *requestBuilder) URL(url string) RequestBuilder {
	rb.request.URL = url
	return rb
}

func (rb *requestBuilder) Path(path string) RequestBuilder {
	if strings.HasPrefix(path, "/") {
		rb.request.URL = rb.request.URL + path
	} else {
		rb.request.URL = rb.request.URL + "/" + path
	}
	return rb
}

func (rb *requestBuilder) Pathf(format string, args ...interface{}) RequestBuilder {
	path := fmt.Sprintf(format, args...)
	return rb.Path(path)
}

// Headers methods
func (rb *requestBuilder) Header(key, value string) RequestBuilder {
	rb.request.Headers[key] = value
	return rb
}

func (rb *requestBuilder) Headers(headers map[string]string) RequestBuilder {
	for key, value := range headers {
		rb.request.Headers[key] = value
	}
	return rb
}

func (rb *requestBuilder) ContentType(contentType ContentType) RequestBuilder {
	rb.request.ContentType = contentType
	rb.request.Headers["Content-Type"] = string(contentType)
	return rb
}

func (rb *requestBuilder) Accept(accept string) RequestBuilder {
	rb.request.Headers["Accept"] = accept
	return rb
}

func (rb *requestBuilder) UserAgent(userAgent string) RequestBuilder {
	rb.request.Headers["User-Agent"] = userAgent
	return rb
}

// Query parameters methods
func (rb *requestBuilder) Query(key, value string) RequestBuilder {
	rb.request.QueryParams[key] = value
	return rb
}

func (rb *requestBuilder) QueryParams(params map[string]string) RequestBuilder {
	for key, value := range params {
		rb.request.QueryParams[key] = value
	}
	return rb
}

func (rb *requestBuilder) QueryStruct(v interface{}) RequestBuilder {
	params := structToMap(v)
	return rb.QueryParams(params)
}

// Body methods
func (rb *requestBuilder) Body(body interface{}) RequestBuilder {
	rb.request.Body = body
	return rb
}

func (rb *requestBuilder) BodyReader(reader io.Reader) RequestBuilder {
	rb.request.BodyReader = reader
	return rb
}

func (rb *requestBuilder) JSON(v interface{}) RequestBuilder {
	rb.request.Body = v
	rb.request.ContentType = ContentTypeJSON
	rb.request.Headers["Content-Type"] = string(ContentTypeJSON)
	return rb
}

func (rb *requestBuilder) XML(v interface{}) RequestBuilder {
	rb.request.Body = v
	rb.request.ContentType = ContentTypeXML
	rb.request.Headers["Content-Type"] = string(ContentTypeXML)
	return rb
}

func (rb *requestBuilder) Form(data map[string]string) RequestBuilder {
	rb.request.Body = data
	rb.request.ContentType = ContentTypeForm
	rb.request.Headers["Content-Type"] = string(ContentTypeForm)
	return rb
}

func (rb *requestBuilder) FormData(data map[string][]string) RequestBuilder {
	rb.request.Body = data
	rb.request.ContentType = ContentTypeMultipart
	rb.request.Headers["Content-Type"] = string(ContentTypeMultipart)
	return rb
}

func (rb *requestBuilder) File(fieldName, fileName string, reader io.Reader) RequestBuilder {
	// TODO: Implement file upload
	rb.request.BodyReader = reader
	rb.request.ContentType = ContentTypeMultipart
	return rb
}

// Authentication methods
func (rb *requestBuilder) Auth(auth *AuthConfig) RequestBuilder {
	rb.request.Auth = auth
	return rb
}

func (rb *requestBuilder) BasicAuth(username, password string) RequestBuilder {
	rb.request.Auth = &AuthConfig{
		Type:     AuthTypeBasic,
		Username: username,
		Password: password,
	}
	return rb
}

func (rb *requestBuilder) BearerToken(token string) RequestBuilder {
	rb.request.Auth = &AuthConfig{
		Type:  AuthTypeBearer,
		Token: token,
	}
	return rb
}

func (rb *requestBuilder) APIKey(key, value string) RequestBuilder {
	rb.request.Auth = &AuthConfig{
		Type:   AuthTypeAPIKey,
		Header: key,
		APIKey: value,
	}
	return rb
}

// Request options methods
func (rb *requestBuilder) Timeout(timeout time.Duration) RequestBuilder {
	rb.request.Timeout = timeout
	return rb
}

func (rb *requestBuilder) Context(ctx context.Context) RequestBuilder {
	rb.request.Context = ctx
	return rb
}

func (rb *requestBuilder) FollowRedirects(follow bool) RequestBuilder {
	rb.request.FollowRedirects = follow
	return rb
}

func (rb *requestBuilder) MaxRedirects(max int) RequestBuilder {
	rb.request.MaxRedirects = max
	return rb
}

// Retry methods
func (rb *requestBuilder) Retry(policy *RetryPolicy) RequestBuilder {
	rb.request.RetryPolicy = policy
	return rb
}

func (rb *requestBuilder) RetryAttempts(attempts int) RequestBuilder {
	if rb.request.RetryPolicy == nil {
		rb.request.RetryPolicy = &RetryPolicy{}
	}
	rb.request.RetryPolicy.MaxAttempts = attempts
	return rb
}

func (rb *requestBuilder) RetryDelay(delay time.Duration) RequestBuilder {
	if rb.request.RetryPolicy == nil {
		rb.request.RetryPolicy = &RetryPolicy{}
	}
	rb.request.RetryPolicy.InitialDelay = delay
	return rb
}

// Cache methods
func (rb *requestBuilder) Cache(ttl time.Duration) RequestBuilder {
	rb.request.CacheTTL = ttl
	return rb
}

func (rb *requestBuilder) CacheKey(key string) RequestBuilder {
	rb.request.CacheKey = key
	return rb
}

func (rb *requestBuilder) NoCache() RequestBuilder {
	rb.request.NoCache = true
	return rb
}

// Metadata methods
func (rb *requestBuilder) Metadata(key string, value interface{}) RequestBuilder {
	rb.request.Metadata[key] = value
	return rb
}

// Execute methods
func (rb *requestBuilder) Send() (*Response, error) {
	return rb.SendWithContext(context.Background())
}

func (rb *requestBuilder) SendWithContext(ctx context.Context) (*Response, error) {
	// Build final request
	req, err := rb.Build()
	if err != nil {
		return nil, err
	}

	// Set context
	if req.Context == nil {
		req.Context = ctx
	}

	// Execute request
	return rb.client.DoWithContext(ctx, req)
}

// Response helpers
func (rb *requestBuilder) Expect(statusCode int) (*Response, error) {
	resp, err := rb.Send()
	if err != nil {
		return resp, err
	}

	if resp.StatusCode != statusCode {
		return resp, &HTTPError{
			Code:       1100,
			Message:    fmt.Sprintf("expected status %d, got %d", statusCode, resp.StatusCode),
			Type:       "status",
			StatusCode: resp.StatusCode,
			Response:   resp,
		}
	}

	return resp, nil
}

func (rb *requestBuilder) ExpectJSON(v interface{}) (*Response, error) {
	resp, err := rb.Send()
	if err != nil {
		return resp, err
	}

	if len(resp.Body) > 0 {
		if err := json.Unmarshal(resp.Body, v); err != nil {
			return resp, &HTTPError{
				Code:     1101,
				Message:  fmt.Sprintf("failed to unmarshal JSON: %v", err),
				Type:     "json",
				Response: resp,
			}
		}
	}

	return resp, nil
}

func (rb *requestBuilder) ExpectXML(v interface{}) (*Response, error) {
	resp, err := rb.Send()
	if err != nil {
		return resp, err
	}

	if len(resp.Body) > 0 {
		if err := xml.Unmarshal(resp.Body, v); err != nil {
			return resp, &HTTPError{
				Code:     1102,
				Message:  fmt.Sprintf("failed to unmarshal XML: %v", err),
				Type:     "xml",
				Response: resp,
			}
		}
	}

	return resp, nil
}

func (rb *requestBuilder) ExpectBytes() ([]byte, error) {
	resp, err := rb.Send()
	if err != nil {
		return nil, err
	}

	return resp.Body, nil
}

func (rb *requestBuilder) ExpectString() (string, error) {
	bytes, err := rb.ExpectBytes()
	if err != nil {
		return "", err
	}

	return string(bytes), nil
}

// Build request without sending
func (rb *requestBuilder) Build() (*Request, error) {
	// Apply query parameters to URL
	if len(rb.request.QueryParams) > 0 {
		u, err := url.Parse(rb.request.URL)
		if err != nil {
			return nil, &HTTPError{
				Code:    1103,
				Message: fmt.Sprintf("invalid URL: %v", err),
				Type:    "url",
			}
		}

		q := u.Query()
		for key, value := range rb.request.QueryParams {
			q.Set(key, value)
		}
		u.RawQuery = q.Encode()
		rb.request.URL = u.String()
	}

	// Clone request to avoid mutations
	req := *rb.request

	// Copy maps to avoid shared references
	req.Headers = make(map[string]string)
	for k, v := range rb.request.Headers {
		req.Headers[k] = v
	}

	req.QueryParams = make(map[string]string)
	for k, v := range rb.request.QueryParams {
		req.QueryParams[k] = v
	}

	req.Metadata = make(map[string]interface{})
	for k, v := range rb.request.Metadata {
		req.Metadata[k] = v
	}

	return &req, nil
}

// Helper functions

// structToMap converts struct to map[string]string for query parameters
func structToMap(v interface{}) map[string]string {
	result := make(map[string]string)

	val := reflect.ValueOf(v)
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	if val.Kind() != reflect.Struct {
		return result
	}

	typ := val.Type()
	for i := 0; i < val.NumField(); i++ {
		field := val.Field(i)
		fieldType := typ.Field(i)

		// Skip unexported fields
		if !field.CanInterface() {
			continue
		}

		// Get field name from tag or field name
		name := fieldType.Tag.Get("query")
		if name == "" {
			name = fieldType.Tag.Get("json")
		}
		if name == "" {
			name = strings.ToLower(fieldType.Name)
		}

		// Convert field value to string
		if field.IsValid() && !field.IsZero() {
			result[name] = fmt.Sprintf("%v", field.Interface())
		}
	}

	return result
}
