package answer_parser

import (
	"github.com/fatih/color"
	"github.com/korchasa/speelka/pkg/actions"
	"regexp"
	"strings"
)

type TextParser struct {
}

func NewTextParser() *TextParser {
	return &TextParser{}
}

func (t *TextParser) Parse(from string, color color.Attribute, modelAnswer string) []actions.Action {
	var acts []actions.Action
	acts = append(acts, &actions.Message{
		From:  from,
		Text:  removeCharacterPrefix(from, modelAnswer),
		Color: color,
	})

	for _, call := range extractToolCalls(modelAnswer) {
		call.From = from
		acts = append(acts, &call)
	}

	for _, question := range extractUserQuestions(modelAnswer) {
		question.From = from
		acts = append(acts, &question)
	}
	return acts
}

func removeCharacterPrefix(name string, resp string) string {
	if strings.HasPrefix(resp, name) {
		return strings.TrimPrefix(resp, name+":")
	}
	return resp
}

func extractToolCalls(text string) []actions.ToolRequest {
	var toolCalls []actions.ToolRequest

	// Match @call tool key="value" pattern
	regex := regexp.MustCompile(`@call\s+(\w+)\s*((?:\w+=['"][^"^']+['"]\s*)*)\s*`)
	matches := regex.FindAllStringSubmatch(text, -1)

	for _, match := range matches {
		toolRequest := actions.ToolRequest{
			ToolName:  match[1],
			Arguments: make(map[string]string),
		}

		// Match key="value" pattern within the tool arguments
		argRegex := regexp.MustCompile(`(\w+)=['"]([^'^"]+)['"]`)
		argMatches := argRegex.FindAllStringSubmatch(match[2], -1)

		for _, argMatch := range argMatches {
			toolRequest.Arguments[argMatch[1]] = argMatch[2]
		}

		toolCalls = append(toolCalls, toolRequest)
	}

	return toolCalls
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
