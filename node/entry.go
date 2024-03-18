package node

type Operation int

const (
	OperationPut Operation = iota
	OperationDelete
)

type Entry struct {
	Term      uint64
	Operation Operation
	Key       []byte
	Value     []byte
}
