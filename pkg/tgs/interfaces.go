package tgs

type Request interface {
	Command() Command
}
