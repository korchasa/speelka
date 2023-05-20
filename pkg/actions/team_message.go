package actions

import (
    "fmt"
)

const ActionTypeTeamMessage = "team_message"

type TeamMessage struct {
    From string
    Text string
}

func (t *TeamMessage) Type() string {
    return ActionTypeTeamMessage
}

func (t *TeamMessage) Log() string {
    return fmt.Sprintf("%s said: %s", t.From, t.Text)
}
