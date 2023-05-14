package types

import (
    "fmt"
    "strings"
)

type Action struct {
    Thoughts Thoughts `json:"thoughts,omitempty"`
    ToUser   ToUser   `json:"toUser" yaml:"toUser,omitempty""`
    ToTeam   ToTeam   `json:"toTeam" yaml:"toTeam,omitempty"`
    Command  Command  `json:"command,omitempty" yaml:"command,omitempty"`
}

func (a *Action) TextHistory() (resp []string) {
    if a.ToTeam.Text != "" {
        resp = append(resp, fmt.Sprintf("%s: %s", a.ToTeam.From, a.ToTeam.Text))
    }
    if a.ToUser.Text != "" {
        resp = append(resp, fmt.Sprintf("%s(to User): %s", a.ToUser.From, a.ToUser.Text))
    }
    if a.Command.Name != "" {
        var args []string
        for _, arg := range a.Command.Arguments {
            args = append(args, fmt.Sprintf("\"%s\": \"<%s>\"", arg.Name, arg.Description))
        }
        resp = append(resp, fmt.Sprintf("%s: call `%s`, arguments: %s", a.Command.Name, a.Command.Description, strings.Join(args, ", ")))
    }
    return resp
}

type Thoughts struct {
    Text      string `json:"text" yaml:"text,omitempty"`
    Reasoning string `json:"reasoning" yaml:"reasoning,omitempty"`
    Plan      string `json:"plan" yaml:"plan,omitempty"`
    Criticism string `json:"criticism" yaml:"criticism,omitempty"`
}

type ToUser struct {
    From   string `json:"from" yaml:"from,omitempty"`
    Text   string `json:"text" yaml:"text,omitempty"`
    Reason string `json:"reason" yaml:"reason,omitempty"`
}

type ToTeam struct {
    From   string `json:"from" yaml:"from,omitempty"`
    Text   string `json:"text" yaml:"text,omitempty"`
    Reason string `json:"reason" yaml:"reason,omitempty"`
}

type CommandArgument struct {
    Name        string `json:"name" yaml:"name,omitempty"`
    Description string `json:"description" yaml:"description,omitempty"`
}

type Command struct {
    Name        string                 `json:"name" yaml:"name,omitempty"`
    Description string                 `json:"description" yaml:"description,omitempty"`
    Config      map[string]interface{} `json:"config,omitempty" yaml:"config,omitempty"`
    Arguments   []CommandArgument      `json:"arguments,omitempty" yaml:"arguments,omitempty"`
}

func (t *Command) String() string {
    var args []string
    for _, arg := range t.Arguments {
        args = append(args, fmt.Sprintf("\"%s\": \"<%s>\"", arg.Name, arg.Description))
    }
    return fmt.Sprintf("%s: %s, arguments: %s", t.Name, t.Description, strings.Join(args, ", "))
}

type Message struct {
    From string `json:"from"`
    To   string `json:"to"`
    Text string `json:"text"`
}

func (m Message) CommandRequest() *CommandRequest {
    return nil
}

func (m Message) String() string {
    return fmt.Sprintf("message from `%s` to `%s`: %s", m.From, m.To, m.Text)
}

func (m Message) Type() string {
    return "message"
}

func (m Message) Message() *Message {
    return &m
}

func (m Message) ToolRequest() *CommandRequest {
    return nil
}

func (m Message) CommandResponse() *CommandResponse {
    return nil
}

type CommandRequest struct {
    Command   string            `json:"name"`
    Arguments map[string]string `json:"arguments"`
}

func (r CommandRequest) Type() string {
    return "tool_request"
}

func (r CommandRequest) Message() *Message {
    return nil
}

func (r CommandRequest) CommandRequest() *CommandRequest {
    return &r
}

func (r CommandRequest) CommandResponse() *CommandResponse {
    return nil
}

func (r CommandRequest) String() string {
    return fmt.Sprintf("command `%s` request: query=%+v", r.Command, r.Arguments)
}

type CommandResponse struct {
    Command   string            `json:"command"`
    Arguments map[string]string `json:"arguments"`
    Success   bool              `json:"success"`
    Output    string            `json:"output"`
    Errors    string            `json:"errors"`
}

func (r CommandResponse) String() string {
    return fmt.Sprintf("command `%s` response: success=%t, output=%s, errors=%s", r.Command, r.Success, r.Output, r.Errors)
}
