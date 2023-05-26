package main

import (
	"github.com/fatih/color"
	"github.com/korchasa/speelka/pkg/character"
	"github.com/korchasa/speelka/pkg/character/text_answer_char"
	"github.com/korchasa/speelka/pkg/team/round_robin_team"
	"github.com/korchasa/speelka/pkg/tool"
	"github.com/korchasa/speelka/pkg/ui"
	log "github.com/sirupsen/logrus"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetReportCaller(false)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(
		&log.TextFormatter{
			ForceColors: true,
		},
	)
}

func main() {
	uin := ui.NewConsole()
	tm := round_robin_team.New([]character.Character{
		text_answer_char.MustNew(
			character.Spec{
				Name:        "Leady",
				Description: "specializes in problem-solving and team leadership",
				Role:        "I want you to act as a team leader. Break the task down into subtasks and supervise their execution. Repeat the task and subtasks in each turn. Monitor the progress of the task.",
				Color:       color.FgHiRed,
			},
		),
		text_answer_char.MustNew(
			character.Spec{
				Name:        "Consolleri",
				Description: "skilled in working with the operation system utilities",
				Role:        "I want you to act as an experienced macos user who knows how to work the console.",
				Color:       color.FgHiGreen,
				Tools: []tool.Tool{
					tool.NewConsole(),
				},
			},
		),
		text_answer_char.MustNew(
			character.Spec{
				Name:        "Charty",
				Description: "senior frontend developer",
				Role:        "I want you to act as a senior frontend developer.",
				Color:       color.FgHiBlue,
			},
		),
		text_answer_char.MustNew(
			character.Spec{
				Name:        "Critic",
				Description: "able to identify errors in the team's discussions",
				Role:        "Now as a proofreader, your task is to read through the team discussion and identify any errors they made. Monitor the progress of the task.",
				Color:       color.FgHiWhite,
			},
		),
	}, uin)

	if err := tm.Start("Which 3 processes are taking up the most memory on this device?"); err != nil {
		panic(err)
	}
}
