package simple_answer

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/korchasa/speelka/pkg/actions"
    "github.com/korchasa/speelka/pkg/character"
    "github.com/korchasa/speelka/pkg/character/prompt"
    "github.com/korchasa/speelka/pkg/character/simple_answer/answer_parser"
    "github.com/korchasa/speelka/pkg/chat_gpt"
    "github.com/korchasa/speelka/pkg/tool"
)

const templatePath = "./pkg/character/simple_answer/character.toml.gotpl"

type Character struct {
    name            string
    role            string
    description     string
    tools           []tool.Tool
    color           color.Attribute
    promptGenerator *prompt.Generator
}

func NewSimpleFormat(name string, description string, role string, textColor color.Attribute, tools []tool.Tool) character.Character {
    if tools == nil {
        tools = make([]tool.Tool, 0)
    }
    return &Character{
        name:        name,
        role:        role,
        description: description,
        tools:       tools,
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

func (c *Character) Tools() []tool.Tool {
    return c.tools
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

func (c *Character) RunTool(req *actions.ToolRequest) *actions.ToolResponse {
    for _, cmd := range c.tools {
        if cmd.Name() == req.ToolName {
            return cmd.Call(req)
        }
    }
    return &actions.ToolResponse{
        Request: req,
        Success: false,
        Errors:  fmt.Sprintf("tool `%s` not found", req.ToolName),
    }
}
