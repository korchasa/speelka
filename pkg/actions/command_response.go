package actions

import (
    "fmt"
)

const ActionTypeCommandResponse = "command_response"

type CommandResponse struct {
    Request *CommandRequest
    Success bool   `json:"success"`
    Output  string `json:"output"`
    Errors  string `json:"errors"`
}

func (r *CommandResponse) Type() string {
    return ActionTypeCommandResponse
}

func (r *CommandResponse) TextDescription() string {
    return fmt.Sprintf(`Command: %s
Success: %v
Stdout: %s
Stderr: %s`, r.Request.Log(), r.Success, r.Output, r.Errors)
}

func (r *CommandResponse) Log() string {
    return fmt.Sprintf("Command: %s\nSuccess: %v\nStdout: %s\nStderr: %s", r.Request.Log(), r.Success, r.Output, r.Errors)
}
