package commands

type Command interface {
	Handle(data []byte)
	GetName() string
}
