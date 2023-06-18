package log

type Entry struct {
	Index int
	// TODO: Command is temporary.
	Command []byte
	Term    int
}

// TODO: Checksum for command.
// TODO: Command.
// TODO: Validate command with checksum (check entry has been stored correctly).
