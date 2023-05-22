package actions

type Action interface {
    Type() string
    Log() string
}
