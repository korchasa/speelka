package tokens_counter

import (
    "github.com/sashabaranov/go-openai"
    "github.com/stretchr/testify/assert"
    "testing"
)

func TestCountTokens(t *testing.T) {
    req := openai.ChatCompletionRequest{
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleAssistant,
                Content: "this is a long string in english",
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: "це довгий рядок українською",
            },
            {
                Role:    openai.ChatMessageRoleSystem,
                Content: "это длинная строка на русском",
            },
            {
                Role:    openai.ChatMessageRoleUser,
                Content: "ist eine lange Zeile im Deutschen",
            },
        },
    }
    counter, err := NewCounter()
    assert.NoError(t, err)
    info, err := counter.CalcTokensPerMessage(&req)
    assert.NoError(t, err)
    assert.Equal(t, []MessageInfo{
        {Length: 7, IsSystem: false},
        {Length: 32, IsSystem: false},
        {Length: 30, IsSystem: true},
        {Length: 12, IsSystem: false},
    }, info)
}
