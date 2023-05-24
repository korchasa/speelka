package command

import (
    "bytes"
    "context"
    "fmt"
    "github.com/korchasa/speelka/pkg/actions"
    log "github.com/sirupsen/logrus"
    "os/exec"
    "strings"
    "time"
)

const ToolConsole = "console"

type Console struct {
}

func NewConsole() *Console {
    return &Console{}
}

func (c *Console) Name() string {
    return ToolConsole
}

func (c *Console) Description() string {
    return "Run bash command"
}

func (c *Console) Arguments() []Argument {
    return []Argument{
        {
            Name:        "query",
            Description: "bash command",
        },
    }
}

func (c *Console) String() string {
    var parts []string
    for _, arg := range c.Arguments() {
        parts = append(parts, fmt.Sprintf("%s='<%s>'", arg.Name, arg.Description))
    }
    return fmt.Sprintf("@call %s %s - %s", c.Name(), strings.Join(parts, " "), c.Description())
}

func (c *Console) Call(r *actions.CommandRequest) *actions.CommandResponse {
    resp := &actions.CommandResponse{
        Request: r,
    }
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    cmd := exec.CommandContext(ctx, "bash", `-c`, r.Arguments["query"])
    var outb, errb bytes.Buffer
    cmd.Stdout = &outb
    cmd.Stderr = &errb
    err := cmd.Run()
    resp.Output = limitStringLength(outb.String(), 500)
    resp.Errors += limitStringLength(errb.String(), 500)
    if exitErr, ok := err.(*exec.ExitError); ok {
        resp.Errors += fmt.Sprintf("command failed: exit code `%d`", exitErr.ExitCode())
    } else if err != nil {
        log.Warnf("Error running command: %s", err)
    } else {
        resp.Success = true
    }

    return resp
}

func limitStringLength(text string, length int) string {
    if len(text) > length {
        return text[:length] + "..."
    }
    return text
}
