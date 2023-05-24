package character

import (
    "github.com/fatih/color"
    "github.com/korchasa/speelka/pkg/actions"
    "github.com/korchasa/speelka/pkg/tool"
)

type Character interface {
    Init() error
    Respond(problem string, teamChars []Character, history []actions.Action) ([]actions.Action, error)
    Name() string
    Description() string
    Role() string
    Color() color.Attribute
    Tools() []tool.Tool
    RunTool(req *actions.ToolRequest) *actions.ToolResponse
}
