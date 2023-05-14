package chat_gpt

import (
    "context"
    "errors"
    "fmt"
    "github.com/korchasa/spilka/pkg/chat_gpt/limiter"
    "github.com/sashabaranov/go-openai"
    "io"
)

var (
    MaxChatMessagesLengthSumInTokens  = 3000
    MaxChatMessagesLengthSumInSymbols = 12000
    ModelID                           = openai.GPT3Dot5Turbo
)

type ChatGPT struct {
    client  *openai.Client
    limiter *limiter.SmartLimiter
}

func NewChatGPT(client *openai.Client) (*ChatGPT, error) {
    lim, err := limiter.NewSmartLimiter(MaxChatMessagesLengthSumInSymbols, MaxChatMessagesLengthSumInTokens)
    if err != nil {
        return nil, fmt.Errorf("failed to create smart limiter: %w", err)
    }
    return &ChatGPT{
        client:  client,
        limiter: lim,
    }, nil
}

func (c *ChatGPT) AskChatGPTAndStreamResponse(ctx context.Context, req openai.ChatCompletionRequest, resultsCh chan<- StreamMessage) {
    req.Stream = true
    if req.Model == "" {
        req.Model = ModelID
    }
    stream, err := c.client.CreateChatCompletionStream(ctx, req)
    if err != nil {
        resultsCh <- StreamMessage{Err: fmt.Errorf("failed to create chat completion stream for template: %w", err)}
        return
    }
    defer stream.Close()

    for {
        response, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            resultsCh <- StreamMessage{Err: err}
            break
        }
        if err != nil {
            resultsCh <- StreamMessage{Err: fmt.Errorf("failed to receive chat completion message: %w", err)}
            return
        }
        resultsCh <- StreamMessage{RespMsg: response}
    }
}

type StreamMessage struct {
    RespMsg openai.ChatCompletionStreamResponse
    Err     error
}

func (c *ChatGPT) AskChatGPTAndStreamToWriter(ctx context.Context, req openai.ChatCompletionRequest, wr io.Writer) (string, error) {
    err := c.limiter.Limit(&req)
    if err != nil {
        return "", fmt.Errorf("failed to limit chat completion request: %w", err)
    }
    req.Stream = true
    stream, err := c.client.CreateChatCompletionStream(ctx, req)
    if err != nil {
        return "", fmt.Errorf("failed to create chat completion stream: %w", err)
    }
    defer stream.Close()

    var fullResponse string
    for {
        response, err := stream.Recv()
        if errors.Is(err, io.EOF) {
            return fullResponse, nil
        }
        if err != nil {
            return "", fmt.Errorf("failed to receive chat completion message: %w", err)
        }
        fullResponse += response.Choices[0].Delta.Content
        _, err = wr.Write([]byte(response.Choices[0].Delta.Content))
        if err != nil {
            return "", fmt.Errorf("failed to write chat completion message: %w", err)
        }
    }
}

type ParamValue struct {
    ID    string `json:"id"`
    Value string `json:"value"`
}
