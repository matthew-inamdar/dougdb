package node

type Member struct {
	ID      string
	Address string
}

type Config struct {
	Node    *Member
	Members []*Member
}
