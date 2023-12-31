package node

import "fmt"

type RedirectToLeaderError struct {
	LeaderMember *Server
}

func (e RedirectToLeaderError) Error() string {
	return fmt.Sprintf("need to redirect to leader %q at %q", e.LeaderMember.ID, e.LeaderMember.Address)
}
