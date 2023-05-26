package character

import (
	"github.com/fatih/color"
	"github.com/korchasa/speelka/pkg/actions"
	"github.com/korchasa/speelka/pkg/tool"
	"github.com/korchasa/speelka/pkg/ui"
)

type Character interface {
	Respond(problem string, teamChars []Character, history []actions.Action, u *ui.Console) ([]actions.Action, error)
}

type Spec struct {
	Name        string
	Description string
	Role        string
	Color       color.Attribute
	Tools       []tool.Tool
}
