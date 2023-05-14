package limiter

import (
    "fmt"
    "github.com/korchasa/spilka/pkg/chat_gpt/tokens_counter"
    "github.com/sashabaranov/go-openai"
)

func removeOverweightMessages(msgs []openai.ChatCompletionMessage, infos []tokens_counter.MessageInfo, limit int) ([]openai.ChatCompletionMessage, error) {
    if calcSum(infos) <= limit {
        return msgs, nil
    }
    for i := 0; i < len(msgs); i++ {
        if infos[i].IsSystem {
            continue
        }
        msgs = removeElementFromSlice(msgs, i)
        infos = removeElementFromSlice(infos, i)
        return removeOverweightMessages(msgs, infos, limit)
    }
    return nil, fmt.Errorf("failed to remove overweight messages: all messages are system")
}

func removeElementFromSlice[T any](msgs []T, i int) []T {
    return append(msgs[:i], msgs[i+1:]...)
}

func calcSum(infos []tokens_counter.MessageInfo) (sum int) {
    for _, info := range infos {
        sum += info.Length
    }
    return
}
