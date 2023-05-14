package character

type Team interface {
    Characters() []*Character
    Start(q string) error
    TextHistory() []string
}
