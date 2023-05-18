package actions

import (
	"bytes"
	"context"
	"fmt"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"os"
	"os/exec"
	"strings"
	"time"
)

type Command struct {
	Name        string                 `json:"name" yaml:"name,omitempty"`
	Description string                 `json:"description" yaml:"description,omitempty"`
	Config      map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
	Arguments   []CommandArgument      `json:"arguments,omitempty" yaml:"arguments,omitempty"`
}

func (t *Command) String() string {
	var args []string
	for _, arg := range t.Arguments {
		args = append(args, fmt.Sprintf("%s='<%s>'", arg.Name, arg.Description))
	}
	return fmt.Sprintf("@call %s %s - %s", t.Name, strings.Join(args, " "), t.Description)
}

type CommandArgument struct {
	Name        string `json:"name" yaml:"name,omitempty"`
	Description string `json:"description" yaml:"description,omitempty"`
}

type Action interface {
	ToChatCompletionMessage() *openai.ChatCompletionMessage
	ToLog() string
}

type TeamMessage struct {
	From string
	Text string
}

func (t *TeamMessage) ToChatCompletionMessage() *openai.ChatCompletionMessage {
	return &openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Name:    t.From,
		Content: t.Text,
	}
}

func (t *TeamMessage) ToLog() string {
	return fmt.Sprintf("%s said: %s", t.From, t.Text)
}

type CommandCall struct {
	From        string
	CommandName string
	Arguments   map[string]string
}

func (c *CommandCall) ToChatCompletionMessage() *openai.ChatCompletionMessage {
	return nil
}

func (c *CommandCall) ToLog() string {
	args := make([]string, 0, len(c.Arguments))
	for k, v := range c.Arguments {
		args = append(args, fmt.Sprintf("%s=%s", k, v))
	}
	return fmt.Sprintf("@call %s %s", c.CommandName, strings.Join(args, " "))
}

func (c *CommandCall) Call() *CommandResponse {
	resp := &CommandResponse{
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

func callConsole(c *CommandCall, resp *CommandResponse) {
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

func saveFile(c *CommandCall, path string, content string, resp *CommandResponse) {
	err := os.WriteFile(path, []byte(content), 0644)
	if err != nil {
		resp.Success = false
		resp.Errors = fmt.Sprintf("error saving file: %s", err)
	} else {
		resp.Success = true
	}
}

type CommandResponse struct {
	Call    *CommandCall
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Errors  string `json:"errors"`
}

func (r *CommandResponse) ToChatCompletionMessage() *openai.ChatCompletionMessage {
	return &openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleSystem,
		Name:    "system",
		Content: r.ToLog(),
	}
}

func (r *CommandResponse) ToLog() string {
	return fmt.Sprintf("Command: %s\nSuccess: %v\nStdout: %s\nStderr: %s", r.Call.ToLog(), r.Success, r.Output, r.Errors)
}

type UserQuestion struct {
	From     string
	Question string
}

func (u *UserQuestion) ToChatCompletionMessage() *openai.ChatCompletionMessage {
	return nil
}

func (u *UserQuestion) ToLog() string {
	return fmt.Sprintf("ask operator: %s", u.Question)
}

type UserAnswer struct {
	Question *UserQuestion
	Answer   string
}

func (u *UserAnswer) ToChatCompletionMessage() *openai.ChatCompletionMessage {
	return &openai.ChatCompletionMessage{
		Role:    openai.ChatMessageRoleAssistant,
		Content: u.Answer,
	}
}

func (u *UserAnswer) ToLog() string {
	return fmt.Sprintf("operator answer: %s", u.Answer)
}

func limitStringLength(text string, length int) string {
	if len(text) > length {
		return text[:length] + "..."
	}
	return text
}
