package response

// HTTP Status Code Constants
// This package provides standardized HTTP status code constants for consistent use across the application.

// 1xx Informational responses
const (
	// StatusContinue indicates that the initial part of a request has been received
	// and has not yet been rejected by the server.
	StatusContinue = 100

	// StatusSwitchingProtocols indicates that the server understands and is willing
	// to comply with the client's request for a change in the application protocol.
	StatusSwitchingProtocols = 101

	// StatusProcessing indicates that the server has received and is processing the request,
	// but no response is available yet.
	StatusProcessing = 102

	// StatusEarlyHints indicates that the server is likely to send a final response
	// with the header fields included in the informational response.
	StatusEarlyHints = 103
)

// 2xx Success responses
const (
	// StatusOK indicates that the request has succeeded.
	StatusOK = 200

	// StatusCreated indicates that the request has been fulfilled and has resulted
	// in one or more new resources being created.
	StatusCreated = 201

	// StatusAccepted indicates that the request has been accepted for processing,
	// but the processing has not been completed.
	StatusAccepted = 202

	// StatusNonAuthoritativeInfo indicates that the request was successful but the
	// enclosed payload has been modified from that of the origin server's 200 OK response.
	StatusNonAuthoritativeInfo = 203

	// StatusNoContent indicates that the server has successfully fulfilled the request
	// and that there is no additional content to send in the response payload body.
	StatusNoContent = 204

	// StatusResetContent indicates that the server has fulfilled the request and
	// desires that the user agent reset the "document view".
	StatusResetContent = 205

	// StatusPartialContent indicates that the server is successfully fulfilling a
	// range request for the target resource.
	StatusPartialContent = 206

	// StatusMultiStatus provides status for multiple independent operations.
	StatusMultiStatus = 207

	// StatusAlreadyReported is used inside a DAV: propstat response element to avoid
	// enumerating the internal members of multiple bindings to the same collection repeatedly.
	StatusAlreadyReported = 208

	// StatusIMUsed indicates that the server has fulfilled a GET request for the resource,
	// and the response is a representation of the result of one or more instance-manipulations.
	StatusIMUsed = 226
)

// 3xx Redirection responses
const (
	// StatusMultipleChoices indicates that the target resource has more than one representation.
	StatusMultipleChoices = 300

	// StatusMovedPermanently indicates that the target resource has been assigned a new
	// permanent URI and any future references to this resource ought to use one of the returned URIs.
	StatusMovedPermanently = 301

	// StatusFound indicates that the target resource resides temporarily under a different URI.
	StatusFound = 302

	// StatusSeeOther indicates that the server is redirecting the user agent to a different
	// resource that is intended to provide an indirect response to the original request.
	StatusSeeOther = 303

	// StatusNotModified indicates that a conditional GET or HEAD request has been received
	// and would have resulted in a 200 OK response if it were not for the fact that the condition evaluated to false.
	StatusNotModified = 304

	// StatusUseProxy indicates that the requested resource must be accessed through the proxy given by the Location field.
	StatusUseProxy = 305

	// StatusTemporaryRedirect indicates that the target resource resides temporarily under a different URI
	// and the user agent MUST NOT change the request method if it performs an automatic redirection to that URI.
	StatusTemporaryRedirect = 307

	// StatusPermanentRedirect indicates that the target resource has been assigned a new permanent URI
	// and any future references to this resource ought to use one of the enclosed URIs.
	StatusPermanentRedirect = 308
)

// 4xx Client error responses
const (
	// StatusBadRequest indicates that the server cannot or will not process the request
	// due to something that is perceived to be a client error.
	StatusBadRequest = 400

	// StatusUnauthorized indicates that the request has not been applied because
	// it lacks valid authentication credentials for the target resource.
	StatusUnauthorized = 401

	// StatusPaymentRequired is reserved for future use.
	StatusPaymentRequired = 402

	// StatusForbidden indicates that the server understood the request but refuses to authorize it.
	StatusForbidden = 403

	// StatusNotFound indicates that the origin server did not find a current representation
	// for the target resource or is not willing to disclose that one exists.
	StatusNotFound = 404

	// StatusMethodNotAllowed indicates that the method received in the request-line
	// is known by the origin server but not supported by the target resource.
	StatusMethodNotAllowed = 405

	// StatusNotAcceptable indicates that the target resource does not have a current
	// representation that would be acceptable to the user agent.
	StatusNotAcceptable = 406

	// StatusProxyAuthRequired indicates that the client must first authenticate itself with the proxy.
	StatusProxyAuthRequired = 407

	// StatusRequestTimeout indicates that the server did not receive a complete request
	// message within the time that it was prepared to wait.
	StatusRequestTimeout = 408

	// StatusConflict indicates that the request could not be completed due to a conflict
	// with the current state of the target resource.
	StatusConflict = 409

	// StatusGone indicates that access to the target resource is no longer available
	// at the origin server and that this condition is likely to be permanent.
	StatusGone = 410

	// StatusLengthRequired indicates that the server refuses to accept the request
	// without a defined Content-Length.
	StatusLengthRequired = 411

	// StatusPreconditionFailed indicates that one or more conditions given in the
	// request header fields evaluated to false when tested on the server.
	StatusPreconditionFailed = 412

	// StatusRequestEntityTooLarge indicates that the server is refusing to process a request
	// because the request payload is larger than the server is willing or able to process.
	StatusRequestEntityTooLarge = 413

	// StatusRequestURITooLong indicates that the server is refusing to service the request
	// because the request-target is longer than the server is willing to interpret.
	StatusRequestURITooLong = 414

	// StatusUnsupportedMediaType indicates that the origin server is refusing to service
	// the request because the payload is in a format not supported by this method on the target resource.
	StatusUnsupportedMediaType = 415

	// StatusRequestedRangeNotSatisfiable indicates that none of the ranges in the request's
	// Range header field overlap the current extent of the selected resource.
	StatusRequestedRangeNotSatisfiable = 416

	// StatusExpectationFailed indicates that the expectation given in the request's
	// Expect header field could not be met by at least one of the inbound servers.
	StatusExpectationFailed = 417

	// StatusTeapot indicates that the server refuses the attempt to brew coffee with a teapot.
	// This is a reference to Hyper Text Coffee Pot Control Protocol defined in RFC 2324.
	StatusTeapot = 418

	// StatusMisdirectedRequest indicates that the request was directed at a server that
	// is not able to produce a response.
	StatusMisdirectedRequest = 421

	// StatusUnprocessableEntity indicates that the server understands the content type
	// of the request entity, but was unable to process the contained instructions.
	StatusUnprocessableEntity = 422

	// StatusLocked indicates that the source or destination resource of a method is locked.
	StatusLocked = 423

	// StatusFailedDependency indicates that the method could not be performed on the resource
	// because the requested action depended on another action and that action failed.
	StatusFailedDependency = 424

	// StatusTooEarly indicates that the server is unwilling to risk processing a request
	// that might be replayed.
	StatusTooEarly = 425

	// StatusUpgradeRequired indicates that the server refuses to perform the request
	// using the current protocol but might be willing to do so after the client upgrades to a different protocol.
	StatusUpgradeRequired = 426

	// StatusPreconditionRequired indicates that the origin server requires the request to be conditional.
	StatusPreconditionRequired = 428

	// StatusTooManyRequests indicates that the user has sent too many requests in a given amount of time.
	StatusTooManyRequests = 429

	// StatusRequestHeaderFieldsTooLarge indicates that the server is unwilling to process the request
	// because its header fields are too large.
	StatusRequestHeaderFieldsTooLarge = 431

	// StatusUnavailableForLegalReasons indicates that the server is denying access to the resource
	// as a consequence of a legal demand.
	StatusUnavailableForLegalReasons = 451
)

// 5xx Server error responses
const (
	// StatusInternalServerError indicates that the server encountered an unexpected condition
	// that prevented it from fulfilling the request.
	StatusInternalServerError = 500

	// StatusNotImplemented indicates that the server does not support the functionality
	// required to fulfill the request.
	StatusNotImplemented = 501

	// StatusBadGateway indicates that the server, while acting as a gateway or proxy,
	// received an invalid response from an inbound server it accessed while attempting to fulfill the request.
	StatusBadGateway = 502

	// StatusServiceUnavailable indicates that the server is currently unable to handle the request
	// due to a temporary overload or scheduled maintenance.
	StatusServiceUnavailable = 503

	// StatusGatewayTimeout indicates that the server, while acting as a gateway or proxy,
	// did not receive a timely response from an upstream server it needed to access in order to complete the request.
	StatusGatewayTimeout = 504

	// StatusHTTPVersionNotSupported indicates that the server does not support,
	// or refuses to support, the major version of HTTP that was used in the request message.
	StatusHTTPVersionNotSupported = 505

	// StatusVariantAlsoNegotiates indicates that the server has an internal configuration error.
	StatusVariantAlsoNegotiates = 506

	// StatusInsufficientStorage indicates that the method could not be performed on the resource
	// because the server is unable to store the representation needed to successfully complete the request.
	StatusInsufficientStorage = 507

	// StatusLoopDetected indicates that the server terminated an operation because
	// it encountered an infinite loop while processing a request.
	StatusLoopDetected = 508

	// StatusNotExtended indicates that further extensions to the request are required
	// for the server to fulfill it.
	StatusNotExtended = 510

	// StatusNetworkAuthenticationRequired indicates that the client needs to authenticate
	// to gain network access.
	StatusNetworkAuthenticationRequired = 511
)

// Helper functions for status code categories

// IsInformational returns true if the status code is in the 1xx range (informational responses).
func IsInformational(code int) bool {
	return code >= 100 && code < 200
}

// IsSuccess returns true if the status code is in the 2xx range (successful responses).
func IsSuccess(code int) bool {
	return code >= 200 && code < 300
}

// IsRedirection returns true if the status code is in the 3xx range (redirection responses).
func IsRedirection(code int) bool {
	return code >= 300 && code < 400
}

// IsClientError returns true if the status code is in the 4xx range (client error responses).
func IsClientError(code int) bool {
	return code >= 400 && code < 500
}

// IsServerError returns true if the status code is in the 5xx range (server error responses).
func IsServerError(code int) bool {
	return code >= 500 && code < 600
}

// IsError returns true if the status code indicates an error (4xx or 5xx).
func IsError(code int) bool {
	return IsClientError(code) || IsServerError(code)
}
