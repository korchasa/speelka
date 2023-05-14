package roundrobin_team

import (
    "fmt"
    "github.com/Masterminds/sprig/v3"
    "github.com/korchasa/spilka/pkg/character"
    "github.com/korchasa/spilka/pkg/types"
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
    "text/template"
)

var (
    tpl *template.Template
)

func init() {
    cnt, err := os.ReadFile("./pkg/roundrobin_team/char.tmpl.md")
    if err != nil {
        log.Panicf("failed to read template: %v", err)
    }
    tpl, err = template.New("prompt").Funcs(sprig.FuncMap()).Parse(string(cnt))
    if err != nil {
        log.Panicf("failed to parse template: %v", err)
    }
}

type RoundRobinTeam struct {
    chars   []*character.Character
    history []*types.Action
    problem string
}

func NewTeam(chars []*character.Character) character.Team {
    return &RoundRobinTeam{
        chars:   chars,
        history: make([]*types.Action, 0),
    }
}

func (t *RoundRobinTeam) Start(query string) error {
    t.problem = query
    t.loop()
    return nil
}

func (t *RoundRobinTeam) loop() {
    for {
        for _, c := range t.chars {
            log.Debugf("Character: %s", c.Name)
            p := t.generatePromptFromTemplate(tpl, c)
            //log.Debugf("Prompt for `%s`:\n%s", c.Name, p)
            resp, err := c.ProcessMessage(p)
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

func (t *RoundRobinTeam) generatePromptFromTemplate(tpl *template.Template, char *character.Character) string {
    stringWriter := &strings.Builder{}
    var view = struct {
        Problem        string
        Self           *character.Character
        TeamCharacters []*character.Character
        TeamHistory    []string
    }{
        Problem:        t.problem,
        Self:           char,
        TeamCharacters: t.Characters(),
        TeamHistory:    t.TextHistory(),
    }
    err := tpl.Execute(stringWriter, view)
    if err != nil {
        panic(err)
    }
    return stringWriter.String()
}
