package chat_gpt

import (
	"context"
	"fmt"
	"github.com/korchasa/spilka/pkg/actions"
	"github.com/korchasa/spilka/pkg/chat_gpt/limiter"
	"github.com/korchasa/spilka/pkg/utils"
	"github.com/sashabaranov/go-openai"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v2"
	"os"
	"strings"
)

var (
	GPTClient *GPT
)

func init() {
	GPTClient = NewGPT(
		openai.NewClient(utils.EnsureEnv("OPENAI_API_KEY")),
		"./var/work.log",
	)
}

type GPT struct {
	client  *openai.Client
	logFile *os.File
	limiter *limiter.SmartLimiter
}

func (g *GPT) AskChatGPT(receiverName string, receiverRole string, problem string, system string, history []actions.Action) (string, error) {
	req := openai.ChatCompletionRequest{
		Model: openai.GPT3Dot5Turbo,
		Messages: []openai.ChatCompletionMessage{
			{
				Role:    openai.ChatMessageRoleSystem,
				Content: system,
			},
			{
				Role:    openai.ChatMessageRoleUser,
				Content: fmt.Sprintf("%s, what is you first turn?", receiverName),
			},
		},
	}
	for _, act := range history {
		m := act.ToChatCompletionMessage()
		if m != nil {
			req.Messages = append(req.Messages, *m)
		}
	}
	err := g.limiter.Limit(&req)
	if err != nil {
		return "", fmt.Errorf("failed to limit: %w", err)
	}
	resp, err := g.client.CreateChatCompletion(context.TODO(), req)
	if err != nil {
		log.Panic(err)
	}
	rec := struct {
		Receiver string
		Request  openai.ChatCompletionRequest
		Response openai.ChatCompletionResponse
	}{
		Receiver: receiverName,
		Request:  req,
		Response: resp,
	}
	yr, _ := yaml.Marshal(rec)
	if _, err := g.logFile.WriteString(string(yr) + "\n===================="); err != nil {
		log.Error(err)
	}
	content := resp.Choices[0].Message.Content
	content = strings.Trim(content, "`")

	return content, nil
}

func NewGPT(client *openai.Client, logPath string) *GPT {
	f, err := os.OpenFile(logPath,
		os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Panic(err)
	}
	f.Truncate(0)
	f.Seek(0, 0)
	lim, err := limiter.NewSmartLimiter(3000, 12000)
	if err != nil {
		log.Panic(err)
	}
	return &GPT{
		client:  client,
		logFile: f,
		limiter: lim,
	}
}
