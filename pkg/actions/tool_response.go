package actions

import (
    "fmt"
)

const ActionTypeToolResponse = "tool_response"

type ToolResponse struct {
    Request *ToolRequest
    Success bool   `json:"success"`
    Output  string `json:"output"`
    Errors  string `json:"errors"`
}

func (r *ToolResponse) Type() string {
    return ActionTypeToolResponse
}

func (r *ToolResponse) TextDescription() string {
    return fmt.Sprintf(`Tool: %s
Success: %v
Stdout: %s
Stderr: %s`, r.Request.Log(), r.Success, r.Output, r.Errors)
}

func (r *ToolResponse) Log() string {
    return fmt.Sprintf("Tool: %s\nSuccess: %v\nStdout: %s\nStderr: %s", r.Request.Log(), r.Success, r.Output, r.Errors)
}
