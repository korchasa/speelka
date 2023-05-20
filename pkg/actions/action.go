package actions

import (
    "fmt"
    "os"
)

type Action interface {
    Type() string
    Log() string
}

func saveFile(c *CommandCall, path string, content string, resp *CommandResult) {
    err := os.WriteFile(path, []byte(content), 0644)
    if err != nil {
        resp.Success = false
        resp.Errors = fmt.Sprintf("error saving file: %s", err)
    } else {
        resp.Success = true
    }
}

func limitStringLength(text string, length int) string {
    if len(text) > length {
        return text[:length] + "..."
    }
    return text
}
