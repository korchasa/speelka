package prompt

import (
    "fmt"
    "github.com/Masterminds/sprig/v3"
    "github.com/korchasa/spilka/pkg/actions"
    "os"
    "strings"
    "text/template"
)

const templatePath = "./pkg/character/prompt/character.toml.gotpl"

type Generator struct {
    tpl *template.Template
}

func NewGenerator() (*Generator, error) {
    cnt, err := os.ReadFile(templatePath)
    if err != nil {
        return nil, fmt.Errorf("failed to read template: %v", err)
    }
    tpl, err := template.New("prompt").Funcs(sprig.FuncMap()).Parse(string(cnt))
    if err != nil {
        return nil, fmt.Errorf("failed to parse template: %v", err)
    }
    return &Generator{
        tpl: tpl,
    }, nil
}

func (g *Generator) GeneratePrompt(problem string, char *CharacterSpec, team []*CharacterSpec, history []actions.Action) (string, error) {
    view := View{
        Problem:        problem,
        Character:      char,
        TeamCharacters: team,
        History:        history,
    }
    stringWriter := &strings.Builder{}
    err := g.tpl.Execute(stringWriter, view)
    if err != nil {
        return "", fmt.Errorf("failed to execute template: %v", err)
    }
    return stringWriter.String(), nil
}

type View struct {
    Problem        string
    Character      *CharacterSpec
    TeamCharacters []*CharacterSpec
    History        []actions.Action
}

type CharacterSpec struct {
    Name        string
    Role        string
    Description string
    Commands    []actions.Command
}
