package answer_parser

import (
    "github.com/korchasa/spilka/pkg/actions"
    "regexp"
    "strings"
)

type TextParser struct {
}

func NewTextParser() *TextParser {
    return &TextParser{}
}

func (t *TextParser) Parse(from string, modelAnswer string) []actions.Action {
    var acts []actions.Action
    acts = append(acts, &actions.TeamMessage{
        From: from,
        Text: removeCharacterPrefix(from, modelAnswer),
    })

    for _, call := range extractCommandCalls(modelAnswer) {
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
