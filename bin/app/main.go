package main

import (
	"github.com/fatih/color"
	"github.com/korchasa/spilka/pkg/actions"
	"github.com/korchasa/spilka/pkg/character"
	"github.com/korchasa/spilka/pkg/team"
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
	team := team.NewTeam([]*character.Character{
		{
			Name:        "Leady",
			Role:        "I want you to act as a team leader. Break the task down into subtasks and supervise their execution. Repeat the task and subtasks in each turn.",
			Description: "specializes in problem-solving and team leadership",
			Color:       color.FgCyan,
		},
		{
			Name:        "Consolleri",
			Role:        "I want you to act as an experienced macos user who knows how to work the console.",
			Description: "skilled in working with the operation system utilities",
			Color:       color.FgHiBlue,
			Commands: []actions.Command{
				{
					Name:        "console",
					Description: "execute bash expressions in macos terminal",
					Arguments: []actions.CommandArgument{
						{
							Name:        "query",
							Description: "console_command_to_execute",
						},
					},
				},
			},
		},
		{
			Name:        "Charty",
			Role:        "I want you to act as a senior frontend developer.",
			Description: "senior frontend developer",
			Color:       color.FgHiRed,
		},
		//{
		//    Name:        "Failly",
		//    Role:        "I want you to act as a file system commander.",
		//    Description: "knows how to work with files",
		//    Color:       color.FgYellow,
		//    Commands: []actions.Command{
		//        {
		//            Name:        "save_file",
		//            Description: "save file",
		//            Arguments: []actions.CommandArgument{
		//                {
		//                    Name:        "filename",
		//                    Description: "file_name",
		//                },
		//                {
		//                    Name:        "content",
		//                    Description: "file_content",
		//                },
		//            },
		//        },
		//    },
		//},
		{
			Name:        "Critic",
			Role:        "Now as a proofreader, your task is to read through the team discussion and identify any errors they made.",
			Color:       color.FgHiWhite,
			Description: "able to identify errors in the team's discussions",
		},
	})

	err := team.Start("I need a HTML page with ring chart of the memory occupied by the 10 largest processes of the operating system")
	if err != nil {
		panic(err)
	}
}
