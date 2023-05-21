package team

import (
    "fmt"
    "github.com/fatih/color"
    "github.com/korchasa/spilka/pkg/actions"
    "github.com/korchasa/spilka/pkg/character"
    "github.com/korchasa/spilka/pkg/ui"
    log "github.com/sirupsen/logrus"
)

type Team struct {
    chars   []*character.Character
    history []actions.Action
    ui      *ui.Console
}

func NewTeam(chars []*character.Character, uin *ui.Console) *Team {
    return &Team{
        history: make([]actions.Action, 0),
        ui:      uin,
        chars:   chars,
    }
}

func (t *Team) Start(problem string) error {
    for _, c := range t.chars {
        if err := c.Init(); err != nil {
            return fmt.Errorf("failed to init character %s: %v", c.Name, err)
        }
    }
    t.loop(problem)
    return nil
}

func (t *Team) loop(problem string) {
    for {
        for _, c := range t.chars {
            t.characterTurn(problem, c)
        }
    }
}

func (t *Team) characterTurn(problem string, character *character.Character) {
    acts, err := character.Respond(problem, t.Characters(), t.history)
    if err != nil {
        log.Fatalf("failed to process message: %v", err)
    }
    for _, act := range acts {
        t.AddToHistory(act)
        switch x := act.(type) {
        case *actions.TeamMessage:
            _, _ = color.New(character.Color).Println(act.Log())
        case *actions.CommandCall:
            run, err := t.ui.Confirmation(fmt.Sprintf("Execute `%v`?", x.Log()))
            if err != nil {
                log.Error(err)
                continue
            }
            if !run {
                continue
            }
            resp := x.Call()
            t.AddToHistory(resp)
            log.Warn(resp.Log())
            t.characterTurn(problem, character)
        case *actions.UserQuestion:
            answer, err := t.ui.Question(x.Question)
            if err != nil {
                log.Error(err)
                continue
            }
            action := &actions.UserAnswer{
                Question: x,
                Answer:   answer,
            }
            t.AddToHistory(action)
            t.characterTurn(problem, character)
        case nil:
            log.Warn("nil action")
        }
    }
}

func (t *Team) Characters() []*character.Character {
    return t.chars
}

func (t *Team) AddToHistory(act actions.Action) {
    t.history = append(t.history, act)
}
