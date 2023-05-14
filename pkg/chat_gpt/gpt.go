package chat_gpt

import (
    "context"
    "encoding/json"
    "fmt"
    "github.com/korchasa/spilka/pkg/types"
    "github.com/korchasa/spilka/pkg/utils"
    "github.com/sashabaranov/go-openai"
    log "github.com/sirupsen/logrus"
    "strings"
)

var (
    GPTClient *GPT
)

func init() {
    GPTClient = NewGPT(openai.NewClient(utils.EnsureEnv("OPENAI_API_KEY")))
}

type GPT struct {
    client *openai.Client
}

func (g *GPT) AskChatGPT(p string) (*types.Action, error) {
    req := openai.ChatCompletionRequest{
        Model: openai.GPT3Dot5Turbo,
        Messages: []openai.ChatCompletionMessage{
            {
                Role:    openai.ChatMessageRoleUser,
                Content: p,
            },
        },
    }
    resp, err := g.client.CreateChatCompletion(context.TODO(), req)
    if err != nil {
        panic(err)
    }
    content := resp.Choices[0].Message.Content
    content = strings.Trim(content, "`")

    log.Debugf("Raw model response: %s", content)

    var msg types.Action
    err = json.Unmarshal([]byte(content), &msg)
    if err != nil {
        log.Warnf("Raw model response: %s", content)
        return nil, fmt.Errorf("failed to unmarshal model response: %w", err)
    }
    return &msg, nil
}

func NewGPT(client *openai.Client) *GPT {
    return &GPT{client: client}
}
