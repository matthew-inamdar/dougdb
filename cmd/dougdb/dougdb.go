package main

import (
	"flag"
	"github.com/matthew-inamdar/dougdb/node"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
)

func main() {
	var (
		nodeConfig    = flag.String("node", "", `the node configuration (e.g. "1,127.0.0.1:8001")`)
		clusterConfig = flag.String("cluster", "", `the cluster configuration (e.g. "2,127.0.0.1:8002;3,127.0.0.1:8003")`)
	)
	flag.Parse()

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		os.Exit(0)
	}()

	c := &node.Config{
		Node:    parseNodeConfig(*nodeConfig),
		Members: parseNodesConfig(*clusterConfig),
	}
	n, err := node.NewNode(c)
	if err != nil {
		panic(err)
	}

	if err := n.Start(); err != nil {
		panic(err)
	}
}

func parseNodeConfig(config string) *node.Member {
	c := strings.Split(config, ",")
	if len(c) != 2 {
		log.Fatalf("invalid format for member %q", config)
	}
	return &node.Member{ID: c[0], Address: c[1]}
}

func parseNodesConfig(configs string) []*node.Member {
	c := strings.Split(configs, ";")
	m := make([]*node.Member, 0, len(c))

	for _, _c := range c {
		m = append(m, parseNodeConfig(_c))
	}

	return m
}
