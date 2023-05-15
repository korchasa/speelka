package team

import (
    "fmt"
    "github.com/AlecAivazis/survey/v2"
    "github.com/Masterminds/sprig/v3"
    "github.com/korchasa/spilka/pkg/character"
    log "github.com/sirupsen/logrus"
    "os"
    "strings"
    "text/template"
)

var (
    tpl *template.Template
)

func init() {
    cnt, err := os.ReadFile("./pkg/team/char.tmpl.md")
    if err != nil {
        log.Panicf("failed to read template: %v", err)
    }
    tpl, err = template.New("prompt").Funcs(sprig.FuncMap()).Parse(string(cnt))
    if err != nil {
        log.Panicf("failed to parse template: %v", err)
    }
}

type Team struct {
    chars   []*character.Character
    history []character.Action
    problem string
}

func NewTeam(chars []*character.Character) *Team {
    return &Team{
        chars:   chars,
        history: make([]character.Action, 0),
    }
}

func (t *Team) Start(query string) error {
    t.problem = query
    t.loop()
    return nil
}

func (t *Team) loop() {
    for {
        for _, c := range t.chars {
            p := t.generatePromptFromTemplate(tpl, c)
            //log.Debugf("Prompt for `%s`:\n%s", c.Name, p)
            acts, err := c.ProcessMessage(p)
            if err != nil {
                log.Fatalf("failed to process message: %v", err)
            }
            if len(acts) == 0 {
                continue
            }
            for _, act := range acts {
                t.AddToHistory(act)
                switch x := act.(type) {
                case *character.TeamText:
                    log.Info(act.ToLog())
                case *character.CommandCall:
                    run := false
                    prompt := &survey.Confirm{
                        Message: fmt.Sprintf("Execute `%v`?", x.ToLog()),
                    }
                    err := survey.AskOne(prompt, &run)
                    if err != nil {
                        log.Fatalf("failed to ask: %v", err)
                        continue
                    }
                    if !run {
                        continue
                    }
                    log.Warn(act.ToLog())
                    resp := x.Call()
                    t.AddToHistory(resp)
                    log.Warn(resp.ToLog())
                case *character.UserQuestion:
                    at := ""
                    prompt := &survey.Input{
                        Message: fmt.Sprintf("%s?", x.Question),
                    }
                    err := survey.AskOne(prompt, &at, survey.WithValidator(survey.Required))
                    if err != nil {
                        log.Fatalf("failed to ask: %v", err)
                        continue
                    }
                    answer := &character.UserAnswer{
                        Question: x,
                        Answer:   at,
                    }
                    t.AddToHistory(answer)
                    log.Warn(answer.ToLog())
                case nil:
                    log.Warn("nil action")
                }
            }
            //spew.Dump(t.history)
        }
    }
}

func (t *Team) findCharacter(to string) (*character.Character, error) {
    for _, c := range t.chars {
        if c.Name == to {
            return c, nil
        }
    }
    return nil, fmt.Errorf("failed to find character: %s", to)
}

func (t *Team) Characters() []*character.Character {
    return t.chars
}

func (t *Team) TextHistory() (history []string) {
    for _, action := range t.history {
        history = append(history, action.ToHistory())
    }
    return history
}

func (t *Team) AddToHistory(act character.Action) {
    t.history = append(t.history, act)
}

func (t *Team) generatePromptFromTemplate(tpl *template.Template, char *character.Character) string {
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
