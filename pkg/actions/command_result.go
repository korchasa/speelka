package actions

import (
    "fmt"
)

const ActionTypeCommandResponse = "command_result"

type CommandResult struct {
    Call    *CommandCall
    Success bool   `json:"success"`
    Output  string `json:"output"`
    Errors  string `json:"errors"`
}

func (r *CommandResult) Type() string {
    return ActionTypeCommandResponse
}

func (r *CommandResult) TextDescription() string {
    return fmt.Sprintf(`Command: %s
Success: %v
Stdout: %s
Stderr: %s`, r.Call.Log(), r.Success, r.Output, r.Errors)
}

func (r *CommandResult) Log() string {
    return fmt.Sprintf("Command: %s\nSuccess: %v\nStdout: %s\nStderr: %s", r.Call.Log(), r.Success, r.Output, r.Errors)
}
