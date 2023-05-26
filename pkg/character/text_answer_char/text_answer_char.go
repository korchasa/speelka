package text_answer_char

import (
	"fmt"
	"github.com/korchasa/speelka/pkg/actions"
	"github.com/korchasa/speelka/pkg/character"
	"github.com/korchasa/speelka/pkg/character/prompt"
	"github.com/korchasa/speelka/pkg/character/text_answer_char/answer_parser"
	"github.com/korchasa/speelka/pkg/chat_gpt"
	"github.com/korchasa/speelka/pkg/tool"
	"github.com/korchasa/speelka/pkg/ui"
	log "github.com/sirupsen/logrus"
)

const templatePath = "./pkg/character/text_answer_char/character.toml.gotpl"

type TextAnswerChar struct {
	character.Spec
	promptGenerator *prompt.Generator
}

func MustNew(spec character.Spec) character.Character {
	if spec.Tools == nil {
		spec.Tools = make([]tool.Tool, 0)
	}
	gen, err := prompt.NewGenerator(templatePath)
	if err != nil {
		log.Fatalf("failed to create prompt generator: %v", err)
	}
	return &TextAnswerChar{
		Spec:            spec,
		promptGenerator: gen,
	}
}

func (c *TextAnswerChar) Respond(problem string, teamChars []character.Character, history []actions.Action, u *ui.Console) (acts []actions.Action, err error) {
	req, err := c.prompt(problem, teamChars, history)
	if err != nil {
		return nil, err
	}
	resp, err := chat_gpt.GPTClient.AskChatGPT(req)
	if err != nil {
		return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
	}
	answerActs := answer_parser.NewTextParser().Parse(c.Name, c.Color, resp)

	for _, act := range answerActs {
		fmt.Println(act.Log())
		switch x := act.(type) {
		case *actions.Message:
			history = append(history, x)
			acts = append(acts, x)
		case *actions.ToolRequest:
			run, err := u.Confirmation(fmt.Sprintf("Execute `%v`?", x.Log()))
			if err != nil {
				return nil, fmt.Errorf("failed to confirm tool execution: %v", err)
			}
			if run {
				resp := c.runTool(x)
				acts = append(acts, resp)
				fmt.Println(resp.Log())
			}
		case *actions.UserQuestion:
			answer, err := u.Question(x.Question)
			if err != nil {
				log.Error(err)
				continue
			}
			action := &actions.UserAnswer{
				Question: x,
				Answer:   answer,
			}
			acts = append(acts, action)
		case nil:
			log.Warn("nil action")
		}
	}

	return acts, nil
}

func (c *TextAnswerChar) prompt(problem string, chars []character.Character, history []actions.Action) (string, error) {
	view := struct {
		Problem        string
		Character      *TextAnswerChar
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

func (c *TextAnswerChar) runTool(req *actions.ToolRequest) *actions.ToolResponse {
	for _, cmd := range c.Tools {
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
