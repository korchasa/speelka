package character

import (
    "reflect"
    "testing"
)

func TestExtractCommandCalls(t *testing.T) {
    tests := []struct {
        name          string
        text          string
        expectedCalls []CommandCall
    }{
        {
            name: "SingleCommandCall",
            text: `some text.
            @call console query="memory usage"
            some text`,
            expectedCalls: []CommandCall{
                {
                    CommandName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage",
                    },
                },
            },
        },
        {
            name: "MultipleCommandCalls",
            text: `- user: hello
- @call console query="memory usage of operating system processes in macOS"
- @call save_file path="/path/to/memory.csv" content="memory usage data extracted from activity monitor"`,
            expectedCalls: []CommandCall{
                {
                    CommandName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage of operating system processes in macOS",
                    },
                },
                {
                    CommandName: "save_file",
                    Arguments: map[string]string{
                        "path":    "/path/to/memory.csv",
                        "content": "memory usage data extracted from activity monitor",
                    },
                },
            },
        },
        {
            name: "MultipleCommandCallsSingleQuote",
            text: `some text
- @call console query='memory usage of operating system processes in macOS'
- @call save_file path='/path/to/memory.csv' content='memory usage data extracted from activity monitor'`,
            expectedCalls: []CommandCall{
                {
                    CommandName: "console",
                    Arguments: map[string]string{
                        "query": "memory usage of operating system processes in macOS",
                    },
                },
                {
                    CommandName: "save_file",
                    Arguments: map[string]string{
                        "path":    "/path/to/memory.csv",
                        "content": "memory usage data extracted from activity monitor",
                    },
                },
            },
        },
    }

    for _, test := range tests {
        t.Run(test.name, func(t *testing.T) {
            result := extractCommandCalls(test.text)

            if !reflect.DeepEqual(result, test.expectedCalls) {
                t.Errorf("extractCommandCalls() failed: expected %v, got %v", test.expectedCalls, result)
            }
        })
    }
}

func TestExtractUserQuestions(t *testing.T) {
    resp := `@user, What is your name?
        @user How old are you?
        @user Where are you from?`
    expected := []UserQuestion{
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
            t.Errorf("Expected question '%s', but got '%s'", expected[i].Question, question.Question)
        }
    }
}
