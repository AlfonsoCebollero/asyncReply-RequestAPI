package entities

type SuccessfulResponse struct {
	WorkflowID     string `json:"workflow-id"`
	WorkflowRunID  string `json:"workflow-run-id"`
	WorkflowStatus string `json:"workflow-status"`
	Message        string `json:"message"`
	StatusCode     string `json:"status-code"`
}

type ErrorResponse struct {
	StatusCode int    `json:"status-code"`
	Message    string `json:"message"`
}

type APIResponse struct {
	Timestamp string             `json:"timestamp"`
	Response  SuccessfulResponse `json:"Response"`
}

type APIError struct {
	Timestamp string        `json:"timestamp"`
	Response  ErrorResponse `json:"response"`
}
