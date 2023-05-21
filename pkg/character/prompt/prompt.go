package prompt

import (
    "fmt"
    "github.com/Masterminds/sprig/v3"
    "os"
    "strings"
    "text/template"
)

type Generator struct {
    tpl *template.Template
}

func NewGenerator(templatePath string) (*Generator, error) {
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

func (g *Generator) Prompt(variables interface{}) (string, error) {
    stringWriter := &strings.Builder{}
    err := g.tpl.Execute(stringWriter, variables)
    if err != nil {
        return "", fmt.Errorf("failed to execute template: %v", err)
    }
    return stringWriter.String(), nil
}
