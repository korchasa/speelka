package round_robin_team

import (
	"github.com/korchasa/speelka/pkg/actions"
	"github.com/korchasa/speelka/pkg/character"
	"github.com/korchasa/speelka/pkg/team"
	"github.com/korchasa/speelka/pkg/ui"
	log "github.com/sirupsen/logrus"
)

type RoundRobinTeam struct {
	chars   []character.Character
	history []actions.Action
	ui      *ui.Console
}

func New(chars []character.Character, uin *ui.Console) team.Team {
	return &RoundRobinTeam{
		history: make([]actions.Action, 0),
		ui:      uin,
		chars:   chars,
	}
}

func (t *RoundRobinTeam) Start(problem string) error {
	for {
		for _, c := range t.chars {
			t.characterTurn(problem, c)
		}
	}
}

func (t *RoundRobinTeam) characterTurn(problem string, character character.Character) {
	acts, err := character.Respond(problem, t.Characters(), t.history, t.ui)
	if err != nil {
		log.Fatalf("failed to process message: %v", err)
	}
	t.history = append(t.history, acts...)
}

func (t *RoundRobinTeam) Characters() []character.Character {
	return t.chars
}
