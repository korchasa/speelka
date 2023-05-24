package answer_parser

import (
    "github.com/korchasa/speelka/pkg/actions"
    "reflect"
    "testing"
)

func TestExtractToolCalls(t *testing.T) {
    tests := []struct {
        name          string
        text          string
        expectedCalls []actions.ToolRequest
    }{
        {
            name: "SingleToolCall",
            text: `some text.
            @call console query="memory usage"
            some text`,
            expectedCalls: []actions.ToolRequest{
                {
                    ToolName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage",
                    },
                },
            },
        },
        {
            name: "MultipleToolCalls",
            text: `- user: hello
- @call console query="memory usage of operating system processes in macOS"
- @call save_file path="/path/to/memory.csv" content="memory usage data extracted from activity monitor"`,
            expectedCalls: []actions.ToolRequest{
                {
                    ToolName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage of operating system processes in macOS",
                    },
                },
                {
                    ToolName: "save_file",
                    Arguments: map[string]string{
                        "path":    "/path/to/memory.csv",
                        "content": "memory usage data extracted from activity monitor",
                    },
                },
            },
        },
        {
            name: "MultipleToolCallsSingleQuote",
            text: `some text
- @call console query='memory usage of operating system processes in macOS'
- @call save_file path='/path/to/memory.csv' content='memory usage data extracted from activity monitor'`,
            expectedCalls: []actions.ToolRequest{
                {
                    ToolName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage of operating system processes in macOS",
                    },
                },
                {
                    ToolName: "save_file",
                    Arguments: map[string]string{
                        "path":    "/path/to/memory.csv",
                        "content": "memory usage data extracted from activity monitor",
                    },
                },
            },
        },
        {
            name: "ComplexTool",
            text: "some text `@call console query='top -l 1 -o MEM | head -n 11 | tail -n 10 | awk '{print \"<tr><td>\" $1 \"</td><td>\" $11 \"</td></tr>\"}'`",
            expectedCalls: []actions.ToolRequest{
                {
                    ToolName: "console",
                    Arguments: map[string]string{
                        "query": "top -l 1 -o MEM | head -n 11 | tail -n 10 | awk '{print \"<tr><td>\" $1 \"</td><td>\" $11 \"</td></tr>\"}'",
                    },
                },
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := extractToolCalls(test.text)

            if !reflect.DeepEqual(result, test.expectedCalls) {
                t.Errorf("extractToolCalls() failed: expected %v, got %v", test.expectedCalls, result)
            }
        })
    }
}

func TestExtractUserQuestions(t *testing.T) {
    resp := `@ask, What is your name?
        @ask How old are you?
        @ask Where are you from?`
    expected := []actions.UserQuestion{
        {Question: "What is your name?"},
        {Question: "How old are you?"},
        {Question: "Where are you from?"},
    }

    result := extractUserQuestions(resp)

    if len(result) != len(expected) {
        t.Errorf("Expected %d user questions, but got %d", len(expected), len(result))
    }

    for i, question := range result {
        if question.Question != expected[i].Question {
            t.Errorf("Expected question '%s', but got '%s'", expected[i].Question, question)
        }
    }
}
