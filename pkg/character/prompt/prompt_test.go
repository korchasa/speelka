package prompt

import (
    "github.com/Masterminds/sprig/v3"
    "github.com/korchasa/spilka/pkg/actions"
    "github.com/stretchr/testify/assert"
    "testing"
    "text/template"
)

const testTemplate = `{{.Problem}} {{.Character.Name}} {{range .TeamCharacters}}{{.Name}} {{end}}{{range .History}}{{.Answer}}{{end}}`

func TestGenerator_GeneratePrompt(t *testing.T) {
    tpl, _ := template.New("prompt").Funcs(sprig.FuncMap()).Parse(testTemplate)
    g := &Generator{
        tpl: tpl,
    }

    char := &CharacterSpec{
        Name:        "Hero",
        Role:        "Leader",
        Description: "The leader of the team",
        Commands:    nil,
    }

    team := []*CharacterSpec{
        {
            Name:        "Sidekick",
            Role:        "Support",
            Description: "The sidekick of the team",
            Commands:    nil,
        },
    }

    history := []actions.Action{
        &actions.UserAnswer{
            Question: nil,
            Answer:   "Answer",
        },
    }

    prompt, err := g.Prompt("Solve the problem", char, team, history)
    assert.NoError(t, err)
    assert.Equal(t, "Solve the problem Hero Sidekick Answer", prompt)
}
