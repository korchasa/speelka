package team

import (
	"fmt"
	"github.com/AlecAivazis/survey/v2"
	"github.com/Masterminds/sprig/v3"
	"github.com/fatih/color"
	"github.com/korchasa/spilka/pkg/actions"
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
	history []actions.Action
	problem string
}

func NewTeam(chars []*character.Character) *Team {
	return &Team{
		chars:   chars,
		history: make([]actions.Action, 0),
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
			characterTurn(t, c)
		}
	}
}

func characterTurn(team *Team, character *character.Character) {
	p := team.generatePromptFromTemplate(tpl, character)
	//log.Debugf("Prompt for `%s`:\n%s", c.Name, p)
	acts, err := character.ProcessMessage(team.problem, p, team.history)
	if err != nil {
		log.Fatalf("failed to process message: %v", err)
	}
	for _, act := range acts {
		team.AddToHistory(act)
		switch x := act.(type) {
		case *actions.TeamMessage:
			_, _ = color.New(character.Color).Println(act.ToLog())
		case *actions.CommandCall:
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
			team.AddToHistory(resp)
			log.Warn(resp.ToLog())
			characterTurn(team, character)
		case *actions.UserQuestion:
			at := ""
			prompt := &survey.Input{
				Message: fmt.Sprintf("%s?", x.Question),
			}
			err := survey.AskOne(prompt, &at, survey.WithValidator(survey.Required))
			if err != nil {
				log.Fatalf("failed to ask: %v", err)
				continue
			}
			answer := &actions.UserAnswer{
				Question: x,
				Answer:   at,
			}
			team.AddToHistory(answer)
			log.Warn(answer.ToLog())
			characterTurn(team, character)
		case nil:
			log.Warn("nil action")
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

func (t *Team) AddToHistory(act actions.Action) {
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
	}
	err := tpl.Execute(stringWriter, view)
	if err != nil {
		panic(err)
	}
	return stringWriter.String()
}
