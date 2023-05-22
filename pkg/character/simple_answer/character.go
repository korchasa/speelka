package simple_answer

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/korchasa/spilka/pkg/actions"
    "github.com/korchasa/spilka/pkg/character"
    "github.com/korchasa/spilka/pkg/character/prompt"
    "github.com/korchasa/spilka/pkg/character/simple_answer/answer_parser"
    "github.com/korchasa/spilka/pkg/chat_gpt"
    "github.com/korchasa/spilka/pkg/command"
)

const templatePath = "./pkg/character/simple_answer/character.toml.gotpl"

type Character struct {
    name            string
    role            string
    description     string
    commands        []command.Command
    color           color.Attribute
    promptGenerator *prompt.Generator
}

func NewSimpleFormat(name string, description string, role string, textColor color.Attribute, commands []command.Command) character.Character {
    if commands == nil {
        commands = make([]command.Command, 0)
    }
    return &Character{
        name:        name,
        role:        role,
        description: description,
        commands:    commands,
        color:       textColor,
    }
}

func (c *Character) Name() string {
    return c.name
}

func (c *Character) Role() string {
    return c.role
}

func (c *Character) Description() string {
    return c.description
}

func (c *Character) Color() color.Attribute {
    return c.color
}

func (c *Character) Commands() []command.Command {
    return c.commands
}

func (c *Character) Init() error {
    gen, err := prompt.NewGenerator(templatePath)
    if err != nil {
        return fmt.Errorf("failed to create prompt generator: %v", err)
    }
    c.promptGenerator = gen
    return nil
}

func (c *Character) Respond(problem string, teamChars []character.Character, history []actions.Action) ([]actions.Action, error) {
    req, err := c.prompt(problem, teamChars, history)
    if err != nil {
        return nil, err
    }
    resp, err := chat_gpt.GPTClient.AskChatGPT(req)
    if err != nil {
        return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
    }
    acts := answer_parser.NewTextParser().Parse(c.name, c.color, resp)
    return acts, nil
}

func (c *Character) prompt(problem string, chars []character.Character, history []actions.Action) (string, error) {
    view := struct {
        Problem        string
        Character      *Character
        TeamCharacters []character.Character
        History        []actions.Action
    }{
        Problem:        problem,
        Character:      c,
        TeamCharacters: chars,
        History:        history,
    }
    p, err := c.promptGenerator.Prompt(view)
    if err != nil {
        return "", fmt.Errorf("failed to generate prompt: %v", err)
    }
    return p, nil
}

func (c *Character) RunCommand(req *actions.CommandRequest) *actions.CommandResponse {
    for _, cmd := range c.commands {
        if cmd.Name() == req.CommandName {
            return cmd.Call(req)
        }
    }
    return &actions.CommandResponse{
        Request: req,
        Success: false,
        Errors:  fmt.Sprintf("command `%s` not found", req.CommandName),
    }
}