package chat_gpt

import (
    "context"
    "fmt"
    "github.com/korchasa/speelka/pkg/chat_gpt/limiter"
    logger "github.com/korchasa/speelka/pkg/chat_gpt/logger"
    "github.com/korchasa/speelka/pkg/utils"
    "github.com/pelletier/go-toml/v2"
    "github.com/sashabaranov/go-openai"
    log "github.com/sirupsen/logrus"
    "strings"
)

var (
    GPTClient *GPT
)

func init() {
    var err error
    GPTClient, err = NewGPT(
        openai.NewClient(utils.EnsureEnv("OPENAI_API_KEY")),
        "./var/work.log",
    )
    if err != nil {
        log.Error(err)
    }
}

type GPT struct {
    client  *openai.Client
    limiter *limiter.SmartLimiter
    logger  *logger.FileLogger
}

func (g *GPT) AskChatGPT(prompt string) (string, error) {
    req := &openai.ChatCompletionRequest{}
    err := toml.Unmarshal([]byte(prompt), req)
    if err != nil {
        return "", fmt.Errorf("failed to unmarshal toml: %v", err)
    }
    err = g.limiter.Limit(req)
    if err != nil {
        return "", fmt.Errorf("failed to limit: %w", err)
    }
    resp, err := g.client.CreateChatCompletion(context.TODO(), *req)
    if err != nil {
        log.Panic(err)
    }
    if err = g.logger.Log(req, &resp); err != nil {
        log.Warnf("failed to log: %v", err)
    }
    content := resp.Choices[0].Message.Content
    content = strings.Trim(content, "`")

    return content, nil
}

func NewGPT(client *openai.Client, logPath string) (*GPT, error) {
    logr, err := logger.NewFileLogger(logPath)
    if err != nil {
        return nil, fmt.Errorf("failed to create logger: %w", err)
    }
    lim, err := limiter.NewSmartLimiter(3000, 12000)
    if err != nil {
        return nil, fmt.Errorf("failed to create limiter: %w", err)
    }
    return &GPT{
        client:  client,
        logger:  logr,
        limiter: lim,
    }, nil
}
