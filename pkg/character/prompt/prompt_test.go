package prompt

import (
    "github.com/Masterminds/sprig/v3"
    "github.com/stretchr/testify/assert"
    "testing"
    "text/template"
)

const testTemplate = `{{.Problem}}`

func TestGenerator_GeneratePrompt(t *testing.T) {
    tpl, _ := template.New("prompt").Funcs(sprig.FuncMap()).Parse(testTemplate)
    g := &Generator{
        tpl: tpl,
    }

    view := struct {
        Problem string
    }{
        Problem: "Solve the problem",
    }

    prompt, err := g.Prompt(view)
    assert.NoError(t, err)
    assert.Equal(t, "Solve the problem", prompt)
}
