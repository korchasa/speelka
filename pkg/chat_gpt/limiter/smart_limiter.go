package limiter

import (
    "fmt"
    "github.com/korchasa/speelka/pkg/chat_gpt/tokens_counter"
    "github.com/sashabaranov/go-openai"
    log "github.com/sirupsen/logrus"
)

type SmartLimiter struct {
    maxSymbols int
    maxTokens  int
    counter    *tokens_counter.Counter
}

func NewSmartLimiter(maxSymbols, maxTokens int) (*SmartLimiter, error) {
    counter, err := tokens_counter.NewCounter()
    if err != nil {
        return nil, fmt.Errorf("failed to create counter: %w", err)
    }
    return &SmartLimiter{
        maxSymbols: maxSymbols,
        maxTokens:  maxTokens,
        counter:    counter,
    }, nil
}

// Limit ограничивает длину сообщений в запросе. Для чатов он удаляет сообщения(начиная с самых старых, но не system), чтобы сумма не превышала лимит токенов. Если не может посчитать по токенам, то ограничивает по символам.
func (l *SmartLimiter) Limit(req *openai.ChatCompletionRequest) error {
    infos, err := l.counter.CalcTokensPerMessage(req)
    if err != nil {
        log.WithError(err).Warn("failed to count tokens")
        infos = l.counter.CalcSymbolsPerMessage(req)
        req.Messages, err = removeOverweightMessages(req.Messages, infos, l.maxSymbols)
        return err
    }
    req.Messages, err = removeOverweightMessages(req.Messages, infos, l.maxTokens)
    return err
}
