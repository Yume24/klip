package cli

type helpError struct {
	usageMessage string
}

func (e *helpError) Error() string {
	return e.usageMessage
}

func (e *helpError) ExitCode() int {
	return 0
}