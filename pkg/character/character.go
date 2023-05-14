package character

import (
    "fmt"
    "github.com/korchasa/spilka/pkg/chat_gpt"
    "github.com/korchasa/spilka/pkg/types"
)

type Character struct {
    Name        string
    Role        string
    Description string
    Commands    []types.Command
}

func (c *Character) ProcessMessage(prompt string) (*types.Action, error) {
    act, err := chat_gpt.GPTClient.AskChatGPT(prompt)
    if err != nil {
        return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
    }
    act.ToTeam.From = c.Name
    act.ToUser.From = c.Name
    //js, _ := yaml.Marshal(act)
    //log.Debugf("===============")
    //log.Debugf("From: %s\n%s", c.Name, string(js))
    //log.Debugf("===============")
    return act, nil
}
