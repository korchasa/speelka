package actions

import (
    "fmt"
)

const ActionTypeUserAnswer = "user_answer"

type UserAnswer struct {
    Question *UserQuestion
    Answer   string
}

func (u *UserAnswer) Type() string {
    return ActionTypeUserAnswer
}

func (u *UserAnswer) Log() string {
    return fmt.Sprintf("operator answer: %s", u.Answer)
}
