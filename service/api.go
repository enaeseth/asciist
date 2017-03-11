package service

type Request struct {
	Width uint   `json:"width"`
	Image []byte `json:"image"`
}

type Success struct {
	Art string `json:"art"`
}

type Failure struct {
	Error string `json:"error"`
}
