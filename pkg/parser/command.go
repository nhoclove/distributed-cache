package parser

// Command represents instructions issue from client
type (
	Command struct {
		Op   Op
		Args []string
	}
)
