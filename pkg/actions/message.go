package actions

import (
    "github.com/fatih/color"
)

const ActionTypeMessage = "message"

type Message struct {
    From  string
    Text  string
    Color color.Attribute
}

func (m *Message) Type() string {
    return ActionTypeMessage
}

func (m *Message) Log() string {
    return color.New(m.Color).Sprintf("%s said: %s", m.From, m.Text)
}
