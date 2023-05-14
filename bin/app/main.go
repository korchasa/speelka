package main

import (
    "github.com/korchasa/spilka/pkg/character"
    "github.com/korchasa/spilka/pkg/roundrobin_team"
    "github.com/korchasa/spilka/pkg/types"
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
    team := roundrobin_team.NewTeam([]*character.Character{
        {
            Name:        "Speaky",
            Role:        "I want you to act as a project manager and team lead. Please discuss the problem with the team.",
            Description: "knows how to solve project problems and how to lead team",
        },
        {
            Name:        "Consi",
            Role:        "I want you to act as a macos power user and senior admin. Please discuss the problem with team.",
            Description: "can work with macos console",
            Commands: []types.Command{
                {
                    Name:        "console",
                    Description: "knows how to work with the macos console",
                    Arguments: []types.CommandArgument{
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
            Role:        "I want you to act as a senior frontend developer. Please discuss the problem with team.",
            Description: "can build charts in html from given data",
        },
        {
            Name:        "Filly",
            Role:        "I want you to act as a file system commander. Please discuss the problem with team.",
            Description: "knows how to work with files",
            Commands: []types.Command{
                {
                    Name:        "save_file",
                    Description: "save file",
                    Arguments: []types.CommandArgument{
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
            Name:        "Crit",
            Role:        "Now as a proofreader, your task is to read through the team discussion and identify any errors they made.",
            Description: "can find errors in the discussion",
        },
    })

    err := team.Start("Construct a ring diagram of the memory occupied by the operating system processes")
    if err != nil {
        panic(err)
    }
}
