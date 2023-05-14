package roundrobin_team

import (
    "fmt"
    "github.com/korchasa/spilka/pkg/character"
    "github.com/korchasa/spilka/pkg/types"
    log "github.com/sirupsen/logrus"
)

type RoundRobinTeam struct {
    chars   []*character.Character
    history []*types.Action
}

func NewTeam(chars []*character.Character) character.Team {
    return &RoundRobinTeam{
        chars:   chars,
        history: make([]*types.Action, 0),
    }
}

func (t *RoundRobinTeam) Start(query string) error {
    msg := &types.Action{
        ToTeam: types.ToTeam{
            From: "user",
            Text: query,
        },
    }
    t.AddToHistory(msg)
    t.loop()
    return nil
}

func (t *RoundRobinTeam) loop() {
    for {
        for _, c := range t.chars {
            resp, err := c.ProcessMessage(t)
            if err != nil {
                log.Fatalf("failed to process message: %v", err)
            }
            if resp == nil {
                continue
            }
            t.AddToHistory(resp)
            if len(resp.TextHistory()) > 0 {
                for _, text := range resp.TextHistory() {
                    log.Info(text)
                }
            }
        }
    }
}

func (t *RoundRobinTeam) findCharacter(to string) (*character.Character, error) {
    for _, c := range t.chars {
        if c.Name == to {
            return c, nil
        }
    }
    return nil, fmt.Errorf("failed to find character: %s", to)
}

func (t *RoundRobinTeam) Characters() []*character.Character {
    return t.chars
}

func (t *RoundRobinTeam) TextHistory() (history []string) {
    for _, action := range t.history {
        history = append(history, action.TextHistory()...)
    }
    return history
}

func (t *RoundRobinTeam) AddToHistory(msg *types.Action) {
    t.history = append(t.history, msg)
}
