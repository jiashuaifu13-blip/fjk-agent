package model

type Request struct {
	Model       string    `json:"model,omitempty"`
	Messages    []Message `json:"messages"`
	Temperature *float64  `json:"temperature,omitempty"`
}

type Response struct {
	Content      string `json:"content"`
	FinishReason string `json:"finish_reason,omitempty"`
}
