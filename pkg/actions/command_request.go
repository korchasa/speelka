package actions

import (
    "fmt"
    "strings"
)

const ActionTypeCommandRequest = "command_request"

type CommandRequest struct {
    From        string
    CommandName string
    Arguments   map[string]string
}

func (c *CommandRequest) Type() string {
    return ActionTypeCommandRequest
}

func (c *CommandRequest) Log() string {
    args := make([]string, 0, len(c.Arguments))
    for k, v := range c.Arguments {
        args = append(args, fmt.Sprintf("%s=%s", k, v))
    }
    return fmt.Sprintf("@call %s %s", c.CommandName, strings.Join(args, " "))
}
