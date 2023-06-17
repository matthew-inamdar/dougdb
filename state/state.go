package state

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

type State interface {
	HandleAppendEntriesRPC(rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error)
	HandleRequestVoteRPC(rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error)
}
