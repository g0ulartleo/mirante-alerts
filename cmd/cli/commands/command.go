package commands

type Command interface {
	Name() string
	Run(args []string) error
}
