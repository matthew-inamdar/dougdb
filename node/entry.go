package node

type operation int

const (
	operationPut operation = iota
	operationDelete
)

type Entry struct {
	term      uint64
	operation operation
	key       []byte
	value     []byte
}
