package chat_gpt

import (
    "context"
    "github.com/korchasa/spilka/pkg/utils"
    "github.com/sashabaranov/go-openai"
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

func (g *GPT) AskChatGPT(p string) (string, error) {
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

    return content, nil
}

func NewGPT(client *openai.Client) *GPT {
    return &GPT{client: client}
}
