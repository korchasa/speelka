package character

import (
    "fmt"
    "github.com/Masterminds/sprig/v3"
    "github.com/korchasa/spilka/pkg/chat_gpt"
    "github.com/korchasa/spilka/pkg/types"
    log "github.com/sirupsen/logrus"
    "gopkg.in/yaml.v3"
    "os"
    "strings"
    "text/template"
)

var (
    tpl *template.Template
)

func init() {
    cnt, err := os.ReadFile("./pkg/character/char.tmpl.md")
    if err != nil {
        log.Panicf("failed to read template: %v", err)
    }
    tpl, err = template.New("prompt").Funcs(sprig.FuncMap()).Parse(string(cnt))
    if err != nil {
        log.Panicf("failed to parse template: %v", err)
    }
}

type Character struct {
    Name        string
    Role        string
    Description string
    Commands    []types.Command
}

func (c *Character) ProcessMessage(team Team) (*types.Action, error) {
    p := c.generatePromptFromTemplate(tpl, team)
    //log.Debugf("===============")
    //log.Debugf(p)
    //log.Debugf("===============")
    act, err := chat_gpt.GPTClient.AskChatGPT(p)
    if err != nil {
        return nil, fmt.Errorf("failed to ask chat gpt: %v", err)
    }
    act.ToTeam.From = c.Name
    act.ToUser.From = c.Name
    js, _ := yaml.Marshal(act)
    log.Debugf("===============")
    log.Debugf("From: %s\n%s", c.Name, string(js))
    log.Debugf("===============")
    return act, nil
}

func (c *Character) generatePromptFromTemplate(tpl *template.Template, team Team) string {
    stringWriter := &strings.Builder{}
    var view = struct {
        Self           *Character
        TeamCharacters []*Character
        TeamHistory    []string
    }{
        Self:           c,
        TeamCharacters: team.Characters(),
        TeamHistory:    team.TextHistory(),
    }
    err := tpl.Execute(stringWriter, view)
    if err != nil {
        panic(err)
    }
    return stringWriter.String()
}
