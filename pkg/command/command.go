package command

import (
    "github.com/korchasa/spilka/pkg/actions"
)

type CommonCommand struct {
    Name        string
    Description string
    Arguments   []Argument
}

type Command interface {
    Name() string
    Description() string
    Arguments() []Argument
    Call(r *actions.CommandRequest) *actions.CommandResponse
    String() string
}

type Argument struct {
    Name        string `json:"name" yaml:"name,omitempty"`
    Description string `json:"description" yaml:"description,omitempty"`
}

//func (t *Command) String() string {
//    var args []string
//    for _, arg := range t.Arguments {
//        args = append(args, fmt.Sprintf("%s='<%s>'", arg.Name, arg.Description))
//    }
//    return fmt.Sprintf("@call %s %s - %s", t.Name, strings.Join(args, " "), t.Description)
//}
