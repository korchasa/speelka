package memory

import (
    "github.com/korchasa/spilka/pkg/types"
)

type Memory struct {
    actions []types.Action
}

func NewMemory() *Memory {
    return &Memory{}
}

func (m *Memory) AddAction(act types.Action) {
    m.actions = append(m.actions, act)
}
