package character

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/korchasa/spilka/pkg/actions"
	"github.com/korchasa/spilka/pkg/chat_gpt"
	"regexp"
	"strings"
)

type Character struct {
	Name        string
	Role        string
	Description string
	Commands    []actions.Command
	Color       color.Attribute
}

func (c *Character) ProcessMessage(problem string, prompt string, history []actions.Action) (acts []actions.Action, err error) {
	resp, err := chat_gpt.GPTClient.AskChatGPT(c.Name, c.Role, problem, prompt, history)
	if err != nil {
		return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
	}

	acts = append(acts, &actions.TeamMessage{
		From: c.Name,
		Text: removeCharacterPrefix(c.Name, resp),
	})

	for _, call := range extractCommandCalls(resp) {
		call.From = c.Name
		acts = append(acts, &call)
	}

	for _, question := range extractUserQuestions(resp) {
		question.From = c.Name
		acts = append(acts, &question)
	}

	return acts, nil
}

func removeCharacterPrefix(name string, resp string) string {
	if strings.HasPrefix(resp, name) {
		return strings.TrimPrefix(resp, name+":")
	}
	return resp
}

func extractCommandCalls(text string) []actions.CommandCall {
	var commandCalls []actions.CommandCall

	// Match @call command key="value" pattern
	regex := regexp.MustCompile(`@call\s+(\w+)\s*((?:\w+=['"][^"^']+['"]\s*)*)\s*`)
	matches := regex.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		commandCall := actions.CommandCall{
			CommandName: match[1],
			Arguments:   make(map[string]string),
		}

		// Match key="value" pattern within the command arguments
		argRegex := regexp.MustCompile(`(\w+)=['"]([^'^"]+)['"]`)
		argMatches := argRegex.FindAllStringSubmatch(match[2], -1)

		for _, argMatch := range argMatches {
			commandCall.Arguments[argMatch[1]] = argMatch[2]
		}

		commandCalls = append(commandCalls, commandCall)
	}

	return commandCalls
}

func extractUserQuestions(resp string) (calls []actions.UserQuestion) {
	re := regexp.MustCompile(`@ask[,\s](.+)$`)
	for _, sub := range strings.Split(resp, "\n") {
		matches := re.FindAllStringSubmatch(sub, -1)
		for _, match := range matches {
			calls = append(calls, actions.UserQuestion{
				Question: strings.Trim(match[1], ", "),
			})
		}
	}
	return calls
}
