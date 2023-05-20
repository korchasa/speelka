package actions

import (
    "fmt"
)

const ActionTypeUserQuestion = "user_question"

type UserQuestion struct {
    From     string
    Question string
}

func (u *UserQuestion) Type() string {
    return ActionTypeUserQuestion
}

func (u *UserQuestion) Log() string {
    return fmt.Sprintf("ask operator: %s", u.Question)
}
