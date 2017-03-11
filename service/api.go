package service

// Request is the shape of an ASCII art conversion request.
// (Note that encoding/json will automatically base64-encode Image).
type Request struct {
	Width uint   `json:"width"`
	Image []byte `json:"image"`
}

// Success is the shape of a successful ASCII art conversion response
// (indicated with a 200 OK status).
type Success struct {
	Art string `json:"art"`
}

// Failure is the shape of a failed ASCII art conversion response
// (indicated with an HTTP status ≥ 400 and an application/json Content-Type).
// Note that internal server errors or gateway errors may result in clients
// seeing non-JSON responses or JSON responses not matching this shape.
type Failure struct {
	Error string `json:"error"`
}
