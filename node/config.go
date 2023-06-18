package node

type Server struct {
	ID      string
	Address string
}

type Config struct {
	ThisServer *Server
	Servers    []*Server
}

func (c *Config) FindServer(id string) *Server {
	for _, m := range c.Servers {
		if m.ID == id {
			return m
		}
	}

	return nil
}
