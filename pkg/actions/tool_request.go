package actions

import (
    "fmt"
    "strings"
)

const ActionTypeToolRequest = "tool_request"

type ToolRequest struct {
    From      string
    ToolName  string
    Arguments map[string]string
}

func (c *ToolRequest) Type() string {
    return ActionTypeToolRequest
}

func (c *ToolRequest) Log() string {
    args := make([]string, 0, len(c.Arguments))
    for k, v := range c.Arguments {
        args = append(args, fmt.Sprintf("%s=%s", k, v))
    }
    return fmt.Sprintf("@call %s %s", c.ToolName, strings.Join(args, " "))
}
