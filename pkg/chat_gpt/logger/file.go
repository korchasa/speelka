package logger

import (
    "fmt"
    "github.com/sashabaranov/go-openai"
    "gopkg.in/yaml.v3"
    "os"
)

type FileLogger struct {
    logFile *os.File
}

func (l *FileLogger) Log(req *openai.ChatCompletionRequest, resp *openai.ChatCompletionResponse) error {
    rec := struct {
        Request  *openai.ChatCompletionRequest
        Response *openai.ChatCompletionResponse
    }{
        Request:  req,
        Response: resp,
    }
    yr, err := yaml.Marshal(rec)
    if err != nil {
        return fmt.Errorf("failed to marshal yaml: %w", err)
    }
    if _, err := l.logFile.WriteString(string(yr) + "\n===================="); err != nil {
        return fmt.Errorf("failed to write to log file: %w", err)
    }
    return nil
}

func NewFileLogger(logPath string) (*FileLogger, error) {
    f, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY, 0644)
    if err != nil {
        return nil, fmt.Errorf("failed to open log file `%s`: %w", logPath, err)
    }
    if err = f.Truncate(0); err != nil {
        return nil, fmt.Errorf("failed to truncate log file `%s`: %w", logPath, err)
    }
    if _, err = f.Seek(0, 0); err != nil {
        return nil, fmt.Errorf("failed to seek log file `%s`: %w", logPath, err)
    }
    return &FileLogger{
        logFile: f,
    }, nil
}
