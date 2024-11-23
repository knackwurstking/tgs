package tgs

type API interface {
	Token() string
	SetToken(token string)
	URL(command Command) string
	Send(request Request) ([]byte, error)
}

type Request interface {
	Command() Command
}
