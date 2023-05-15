package character

import (
    "fmt"
    "github.com/korchasa/spilka/pkg/chat_gpt"
    "regexp"
    "strings"
)

type Character struct {
    Name        string
    Role        string
    Description string
    Commands    []Command
}

func (c *Character) ProcessMessage(prompt string) (acts []Action, err error) {
    resp, err := chat_gpt.GPTClient.AskChatGPT(prompt)
    if err != nil {
        return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
    }

    acts = append(acts, &TeamText{
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

func extractCommandCalls(text string) []CommandCall {
    var commandCalls []CommandCall

    // Match @call command key="value" pattern
    regex := regexp.MustCompile(`@call\s+(\w+)\s*((?:\w+=['"][^"^']+['"]\s*)*)\s*`)
    matches := regex.FindAllStringSubmatch(text, -1)

    for _, match := range matches {
        commandCall := CommandCall{
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

func extractUserQuestions(resp string) (calls []UserQuestion) {
    re := regexp.MustCompile(`@user[,\s](.+)$`)
    for _, sub := range strings.Split(resp, "\n") {
        matches := re.FindAllStringSubmatch(sub, -1)
        for _, match := range matches {
            calls = append(calls, UserQuestion{
                Question: strings.Trim(match[1], ", "),
            })
        }
    }
    return calls
}
