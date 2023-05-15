package main

import (
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
            Name:        "Leader",
            Role:        "Leader, i want you to act as a team lead. Use the strengths of team members.",
            Description: "knows how to solve project problems and how to lead team",
        },
        {
            Name:        "Consolleri",
            Role:        "Consolleri, i want you to act as a macos power user and senior admin.",
            Description: "can work with macos console",
            Commands: []character.Command{
                {
                    Name:        "console",
                    Description: "knows how to work with the macos console",
                    Arguments: []character.CommandArgument{
                        {
                            Name:        "query",
                            Description: "console_command_to_execute",
                        },
                    },
                },
            },
        },
        {
            Name:        "Charter",
            Role:        "Charter, i want you to act as a senior frontend developer.",
            Description: "can build charts in html from given data",
        },
        {
            Name:        "Failler",
            Role:        "Failler, i want you to act as a file system commander.",
            Description: "knows how to work with files",
            Commands: []character.Command{
                {
                    Name:        "save_file",
                    Description: "save file",
                    Arguments: []character.CommandArgument{
                        {
                            Name:        "path",
                            Description: "file_path",
                        },
                        {
                            Name:        "content",
                            Description: "file_content",
                        },
                    },
                },
            },
        },
        {
            Name:        "Critic",
            Role:        "Now as a proofreader, your task is to read through the team discussion and identify any errors they made.",
            Description: "can find errors in the discussion",
        },
    })

    err := team.Start("Construct a ring chart of the memory occupied by the 10 largest processes of the operating system")
    if err != nil {
        panic(err)
    }
}
