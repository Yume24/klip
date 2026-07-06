package cli

// Error returned when the -h or --help flag is passed
// Apart from the message it also defines a custom exit code
// so calls to -h/--help don't exit with 1
type helpError struct {
	usageMessage string
}

func (e *helpError) Error() string {
	return e.usageMessage
}

func (e *helpError) ExitCode() int {
	return 0
}
