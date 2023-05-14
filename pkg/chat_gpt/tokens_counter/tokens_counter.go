package tokens_counter

import (
    "fmt"
    tokenizer "github.com/samber/go-gpt-3-encoder"
    "github.com/sashabaranov/go-openai"
)

type MessageInfo struct {
    Length   int
    IsSystem bool
}

type Counter struct {
    tokensEncoder *tokenizer.Encoder
}

func NewCounter() (*Counter, error) {
    enc, err := tokenizer.NewEncoder()
    if err != nil {
        return nil, fmt.Errorf("failed to create encoder: %w", err)
    }
    return &Counter{
        tokensEncoder: enc,
    }, nil
}

func (c *Counter) CalcTokensSum(req *openai.ChatCompletionRequest) (sum int, err error) {
    for _, msg := range req.Messages {
        encoded, err := c.tokensEncoder.Encode(msg.Content)
        if err != nil {
            return sum, fmt.Errorf("failed to encode message: %w", err)
        }
        sum += len(encoded)
    }
    return
}

func (c *Counter) CalcTokensPerMessage(req *openai.ChatCompletionRequest) (info []MessageInfo, err error) {
    for _, msg := range req.Messages {
        encoded, err := c.tokensEncoder.Encode(msg.Content)
        if err != nil {
            return nil, fmt.Errorf("failed to encode message: %w", err)
        }
        info = append(info, MessageInfo{
            Length:   len(encoded),
            IsSystem: msg.Role == openai.ChatMessageRoleSystem,
        })
    }
    return
}

func (c *Counter) CalcSymbolsSum(req *openai.ChatCompletionRequest) (sum int) {
    for _, msg := range req.Messages {
        sum += len(msg.Content)
    }
    return
}

func (c *Counter) CalcSymbolsPerMessage(req *openai.ChatCompletionRequest) (info []MessageInfo) {
    for _, msg := range req.Messages {
        info = append(info, MessageInfo{
            Length:   len(msg.Content),
            IsSystem: msg.Role == openai.ChatMessageRoleSystem,
        })
    }
    return
}
