package node

type LeaderRole struct{}

func (r LeaderRole) HandleAppendEntriesRPC(n *Node, rpc AppendEntriesRPCRequest) (AppendEntriesRPCResponse, error) {
	// Handle new leader.
	if rpc.Term > n.currentTerm.Get() {
		if err := n.currentTerm.Set(rpc.Term); err != nil {
			return AppendEntriesRPCResponse{}, err
		}
		n.currentLeader = n.config.FindServer(rpc.LeaderID)
		n.role = &FollowerRole{}
		return n.role.HandleAppendEntriesRPC(n, rpc)
	}

	// Handle previous term leader.
	if rpc.Term < n.currentTerm.Get() {
		return AppendEntriesRPCResponse{Term: n.currentTerm.Get(), Success: false}, nil
	}

	// Byzantine fault.
	panic("two leaders present")
}

func (r LeaderRole) HandleRequestVoteRPC(n *Node, rpc RequestVoteRPCRequest) (RequestVoteRPCResponse, error) {
	// Handle previous leader.
	if rpc.Term < n.currentTerm.Get() {
		return RequestVoteRPCResponse{Term: n.currentTerm.Get(), VoteGranted: false}, nil
	}

	// Handle losing vote candidate.
	thisNodeLogAtLeastUpToDate := false
	if n.votedFor != nil && thisNodeLogAtLeastUpToDate {
		return RequestVoteRPCResponse{Term: n.currentTerm.Get(), VoteGranted: false}, nil
	}

	// Handle winning vote candidate.
	return RequestVoteRPCResponse{Term: n.currentTerm.Get(), VoteGranted: true}, nil
}

func (r LeaderRole) HandleExists(n *Node, key string) (bool, error) {
	return n.db.Has(key), nil
}

func (r LeaderRole) HandleGet(n *Node, key string) ([]byte, error) {
	v, ok := n.db.Get(key)
	if !ok {
		return nil, nil
	}
	return v, nil
}

func (r LeaderRole) HandleSet(n *Node, key string, value []byte) error {
	n.db.Set(key, value)
	return nil
}
