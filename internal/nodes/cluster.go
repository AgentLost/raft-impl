package nodes

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"sync"
)

func InitCluster() {
	NodeCluster.AddNode("6000")
	NodeCluster.AddNode("7000")
	NodeCluster.AddNode("8000")
}

var NodeCluster = &Cluster{
	nodes: make(map[string]*Node),
	mu:    &sync.RWMutex{},
}

type Cluster struct {
	nodes map[string]*Node
	mu    *sync.RWMutex
}

func (c *Cluster) AddNode(port string) error {
	node := NewNode(port)

	c.mu.Lock()
	defer c.mu.Unlock()

	_, ok := c.nodes[port]
	if ok {
		return errors.New(fmt.Sprintf("node this port %s already exists", port))
	}

	c.nodes[port] = node

	log.Println("add node", port)
	return nil
}

func (c *Cluster) DeleteNode(port string) bool {
	close(c.nodes[port].sd)
	delete(c.nodes, port)

	log.Println("delete node ", port)
	return true
}

func (c *Cluster) GetInfo() string {
	m := map[string]string{}

	for _, node := range c.nodes {
		m[node.Port] = node.String()
	}

	data, err := json.Marshal(m)
	if err != nil {
		log.Fatal(err)
	}

	return string(data)
}
