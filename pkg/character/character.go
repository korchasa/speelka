package character

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/korchasa/spilka/pkg/actions"
    "github.com/korchasa/spilka/pkg/answer_parser"
    "github.com/korchasa/spilka/pkg/character/prompt"
    "github.com/korchasa/spilka/pkg/chat_gpt"
)

type Character struct {
    Name            string
    Role            string
    Description     string
    Commands        []actions.Command
    Color           color.Attribute
    promptGenerator *prompt.Generator
}

func (c *Character) Init() error {
    gen, err := prompt.NewGenerator()
    if err != nil {
        return fmt.Errorf("failed to create prompt generator: %v", err)
    }
    c.promptGenerator = gen
    return nil
}

func (c *Character) Respond(problem string, teamChars []*Character, history []actions.Action) (acts []actions.Action, err error) {
    req, err := c.prompt(problem, teamChars, history)
    if err != nil {
        return nil, err
    }
    resp, err := chat_gpt.GPTClient.AskChatGPT(req)
    if err != nil {
        return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
    }

    acts = answer_parser.NewTextParser().Parse(c.Name, resp)
    return acts, nil
}

func (c *Character) prompt(problem string, chars []*Character, history []actions.Action) (string, error) {
    var team []*prompt.CharacterSpec
    for _, ch := range chars {
        team = append(team, ch.viewSpec())
    }
    p, err := c.promptGenerator.GeneratePrompt(problem, c.viewSpec(), team, history)
    if err != nil {
        return "", fmt.Errorf("failed to generate prompt: %v", err)
    }
    return p, nil
}

func (c *Character) viewSpec() *prompt.CharacterSpec {
    return &prompt.CharacterSpec{
        Name:        c.Name,
        Role:        c.Role,
        Description: c.Description,
        Commands:    c.Commands,
    }
}
