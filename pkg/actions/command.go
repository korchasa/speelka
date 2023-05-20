package actions

import (
    "fmt"
    "strings"
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
