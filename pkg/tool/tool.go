package tool

import (
    "github.com/korchasa/speelka/pkg/actions"
)

type CommonTool struct {
    Name        string
    Description string
    Arguments   []Argument
}

type Tool interface {
    Name() string
    Description() string
    Arguments() []Argument
    Call(r *actions.ToolRequest) *actions.ToolResponse
    String() string
}

type Argument struct {
    Name        string `json:"name" yaml:"name,omitempty"`
    Description string `json:"description" yaml:"description,omitempty"`
}

//func (t *Tool) String() string {
//    var args []string
//    for _, arg := range t.Arguments {
//        args = append(args, fmt.Sprintf("%s='<%s>'", arg.Name, arg.Description))
//    }
//    return fmt.Sprintf("@call %s %s - %s", t.Name, strings.Join(args, " "), t.Description)
//}
