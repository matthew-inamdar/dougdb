package node

type Entry = string

type AppendEntriesRPCRequest struct {
	Term         int
	LeaderID     string
	PrevLogIndex int
	PrevLogTerm  int
	Entries      []Entry
	LeaderCommit int
}

type AppendEntriesRPCResponse struct {
	Term    int
	Success bool
}

type RequestVoteRPCRequest struct {
	Term         int
	CandidateID  string
	LastLogIndex int
	LastLogTerm  int
}

type RequestVoteRPCResponse struct {
	Term        int
	VoteGranted bool
}

type Role interface {
	HandleAppendEntriesRPC(n *Node, rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error)
	HandleRequestVoteRPC(n *Node, rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error)

	HandleExists(n *Node, key string) (bool, error)    // TODO: Error to contain leader config for redirect.
	HandleGet(n *Node, key string) ([]byte, error)     // TODO: Error to contain leader config for redirect.
	HandleSet(n *Node, key string, value []byte) error // TODO: Error to contain leader config for redirect.
}
