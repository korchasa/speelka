package ui

import (
    "fmt"
    "github.com/AlecAivazis/survey/v2"
)

type Console struct {
}

func NewConsole() *Console {
    return &Console{}
}

func (u *Console) Confirmation(text string) (answer bool, err error) {
    prm := &survey.Confirm{
        Message: text,
    }
    err = survey.AskOne(prm, &answer, survey.WithValidator(survey.Required))
    if err != nil {
        return false, fmt.Errorf("failed to ask confirmation: %v", err)
    }
    return answer, nil
}

func (u *Console) Question(text string) (answer string, err error) {
    prm := &survey.Input{
        Message: text,
    }
    err = survey.AskOne(prm, &answer, survey.WithValidator(survey.Required))
    if err != nil {
        return "", fmt.Errorf("failed to ask question: %v", err)
    }
    return
}
