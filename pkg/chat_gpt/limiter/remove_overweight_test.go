package limiter

import (
    "fmt"
    "github.com/korchasa/spilka/pkg/chat_gpt/tokens_counter"
    "github.com/sashabaranov/go-openai"
    "github.com/stretchr/testify/assert"
    "testing"
)

//func TestRemoveOverweightMessages(t *testing.T) {
//    req := []openai.ChatCompletionMessage{
//        {Content: "1"},
//        {Content: "2"},
//        {Content: "3"},
//        {Content: "4"},
//        {Content: "5"},
//    }
//    infos := []MessageInfo{
//        {Length: 1, IsSystem: true},
//        {Length: 1, IsSystem: false},
//        {Length: 1, IsSystem: true},
//        {Length: 1, IsSystem: false},
//        {Length: 1, IsSystem: false},
//    }
//    res, err := removeOverweightMessages(req, infos, 3)
//    assert.NoError(t, err)
//    assert.Equal(t, []openai.ChatCompletionMessage{
//        {Content: "1"},
//        {Content: "3"},
//        {Content: "5"},
//    }, res)
//}

func Test_removeOverweightMessages(t *testing.T) {
    type args struct {
        msgs  []openai.ChatCompletionMessage
        infos []tokens_counter.MessageInfo
        limit int
    }
    tests := []struct {
        name    string
        args    args
        want    []openai.ChatCompletionMessage
        wantErr error
    }{
        {
            name: "should return all messages",
            args: args{
                msgs: []openai.ChatCompletionMessage{
                    {Content: "1"},
                    {Content: "2"},
                    {Content: "3"},
                },
                infos: []tokens_counter.MessageInfo{
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: false},
                    {Length: 1, IsSystem: false},
                },
                limit: 10,
            },
            want: []openai.ChatCompletionMessage{
                {Content: "1"},
                {Content: "2"},
                {Content: "3"},
            },
            wantErr: nil,
        },
        {
            name: "should return part of messages",
            args: args{
                msgs: []openai.ChatCompletionMessage{
                    {Content: "1"},
                    {Content: "2"},
                    {Content: "3"},
                    {Content: "4"},
                    {Content: "5"},
                },
                infos: []tokens_counter.MessageInfo{
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: false},
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: false},
                    {Length: 1, IsSystem: false},
                },
                limit: 3,
            },
            want: []openai.ChatCompletionMessage{
                {Content: "1"},
                {Content: "3"},
                {Content: "5"},
            },
            wantErr: nil,
        },
        {
            name: "should return part of messages, bu all system messages",
            args: args{
                msgs: []openai.ChatCompletionMessage{
                    {Content: "1"},
                    {Content: "2"},
                    {Content: "3"},
                    {Content: "4"},
                },
                infos: []tokens_counter.MessageInfo{
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: false},
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: false},
                },
                limit: 2,
            },
            want: []openai.ChatCompletionMessage{
                {Content: "1"},
                {Content: "3"},
            },
            wantErr: nil,
        },
        {
            name: "error because all messages are system",
            args: args{
                msgs: []openai.ChatCompletionMessage{
                    {Content: "1"},
                    {Content: "2"},
                },
                infos: []tokens_counter.MessageInfo{
                    {Length: 1, IsSystem: true},
                    {Length: 1, IsSystem: true},
                },
                limit: 1,
            },
            want:    nil,
            wantErr: fmt.Errorf("failed to remove overweight messages: all messages are system"),
        },
    }
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := removeOverweightMessages(tt.args.msgs, tt.args.infos, tt.args.limit)
            assert.Equal(t, tt.wantErr, err, "wrong error")
            assert.Equal(t, tt.want, got, "wrong messages")
        })
    }
}
