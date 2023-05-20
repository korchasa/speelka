package actions

import (
    "bytes"
    "context"
    "fmt"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "strings"
    "time"
)

const ActionTypeCommandCall = "command_call"

type CommandCall struct {
    From        string
    CommandName string
    Arguments   map[string]string
}

func (c *CommandCall) Type() string {
    return ActionTypeCommandCall
}

func (c *CommandCall) Log() string {
    args := make([]string, 0, len(c.Arguments))
    for k, v := range c.Arguments {
        args = append(args, fmt.Sprintf("%s=%s", k, v))
    }
    return fmt.Sprintf("@call %s %s", c.CommandName, strings.Join(args, " "))
}

func (c *CommandCall) Call() *CommandResult {
    resp := &CommandResult{
        Call: c,
    }

    switch c.CommandName {
    case "console":
        callConsole(c, resp)
        break
    case "save_file":
        saveFile(c, c.Arguments["filename"], c.Arguments["content"], resp)
        break
    }

    return resp
}

func callConsole(c *CommandCall, resp *CommandResult) {
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, "bash", `-c`, c.Arguments["query"])
    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb
    err := cmd.Run()
    resp.Output = limitStringLength(outb.String(), 500)
    resp.Errors += limitStringLength(errb.String(), 500)
    if exitErr, ok := err.(*exec.ExitError); ok {
        resp.Errors += fmt.Sprintf("command failed: exit code `%d`:%s", exitErr.ExitCode())
    } else if err != nil {
        log.Warnf("Error running command: %s", err)
    } else {
        resp.Success = true
    }

}
