package env

import (
    "fmt"
    "os"
    "strconv"
)

func EnsureEnv(name string) string {
    value := os.Getenv(name)
    if value == "" {
        panic(fmt.Errorf("%s env is empty", name))
    }
    return value
}

func EnsureEnvAsInt(name string) int {
    v, err := strconv.Atoi(EnsureEnv(name))
    if err != nil {
        panic(fmt.Errorf(
            "cannot parse %s environment: %v",
            name,
            err,
        ))
    }
    return v
}

func EnsureEnvOrDefault(name string, defaultValue string) string {
    value := os.Getenv(name)
    if value == "" {
        return defaultValue
    }
    return value
}
func EnsureEnvAsIntOrDefault(name string, defaultValue string) int {
    v, err := strconv.Atoi(EnsureEnvOrDefault(name, defaultValue))
    if err != nil {
        panic(fmt.Errorf(
            "cannot parse %s environment: %v",
            name,
            err,
        ))
    }
    return v
}
