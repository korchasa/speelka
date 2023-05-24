package main

import (
    "github.com/fatih/color"
    "github.com/korchasa/speelka/pkg/character"
    "github.com/korchasa/speelka/pkg/character/simple_answer"
    "github.com/korchasa/speelka/pkg/command"
    "github.com/korchasa/speelka/pkg/team"
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
    tm := team.NewTeam([]character.Character{
        simple_answer.NewSimpleFormat(
            "Leady",
            "specializes in problem-solving and team leadership",
            "I want you to act as a team leader. Break the task down into subtasks and supervise their execution. Repeat the task and subtasks in each turn. Monitor the progress of the task.",
            color.FgCyan,
            nil,
        ),
        simple_answer.NewSimpleFormat(
            "Consolleri",
            "skilled in working with the operation system utilities",
            "I want you to act as an experienced macos user who knows how to work the console.",
            color.FgHiBlue,
            []command.Command{
                command.NewConsole(),
            },
        ),
        simple_answer.NewSimpleFormat(
            "Charty",
            "senior frontend developer",
            "I want you to act as a senior frontend developer.",
            color.FgHiRed,
            nil,
        ),
        //{
        //    Name:        "Failly",
        //    Role:        "I want you to act as a file system commander.",
        //    Description: "knows how to work with files",
        //    Color:       color.FgYellow,
        //    Commands: []actions.Command{
        //        {
        //            Name:        "save_file",
        //            Description: "save file",
        //            Arguments: []actions.Argument{
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
        simple_answer.NewSimpleFormat(
            "Critic",
            "able to identify errors in the team's discussions",
            "Now as a proofreader, your task is to read through the team discussion and identify any errors they made. Monitor the progress of the task.",
            color.FgHiWhite,
            nil,
        ),
    }, uin)

    if err := tm.Start("Get the memory occupied by the 10 largest processes of the operating system"); err != nil {
        panic(err)
    }
}
