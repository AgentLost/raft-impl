package nodes

import (
	"bytes"
	"fmt"
	"io"
	"log"
	"math/rand"
	"net/http"
	"time"
)

type NodeStatus int

func (ns NodeStatus) String() string {
	switch ns {
	case Follower:
		return "FOLLOWER"
	case Candidate:
		return "CANDIDATE"
	case Leader:
		return "LEADER"
	default:
		return "Bad status"
	}
}

const (
	Follower NodeStatus = iota
	Candidate
	Leader
)

type Node struct {
	Port   string
	Status NodeStatus
	State  int
	Votes  []string

	voted     bool
	vote      chan struct{}
	heartBeat chan struct{}
	sd        chan struct{}
}

func (n *Node) run() {
	for {
		t := time.Duration(rand.Intn(10)+5) * time.Second
		select {
		case <-time.After(t):

			n.Status = n.voting()
			if n.Status == Leader {
				n.heart()
			}
		case <-n.vote:
		case <-n.heartBeat:
			n.Status = Follower
			n.voted = false
		case <-n.sd:
			return
		}
	}
}

func (n *Node) voting() NodeStatus {
	n.voted = true
	n.Votes = append(n.Votes, n.Port)
	client := http.Client{Timeout: 1 * time.Second}
	m := map[string]int{}
	m[n.Port]++

	for key, node := range NodeCluster.nodes {
		if key == n.Port {
			continue
		}

		buffer := bytes.NewBuffer([]byte(n.Port))

		r, err := client.Post("http://localhost:"+node.Port+"/vote", "text/plain", buffer)
		if err != nil {
			log.Println(err)
			continue
		}
		data, err := io.ReadAll(r.Body)

		if err != nil {
			log.Println(err)
			continue
		}

		m[string(data)]++
	}

	keyMax := n.Port
	var max int
	for key, val := range m {
		if val > max {
			keyMax = key
			max = val
		}
	}

	log.Println(m)

	if keyMax == n.Port {
		return Leader
	}

	return Follower
}

func (n *Node) heart() {
	for {
		t := 3 * time.Second
		select {
		case <-time.After(t):
			client := http.Client{}
			for key, node := range NodeCluster.nodes {
				if key == n.Port {
					continue
				}
				_, err := client.Get("http://localhost:" + node.Port + "/heart-beat")
				if err != nil {
					log.Println(err)
				}
			}
		case <-n.sd:
			log.Println("shutdown heart")
			return
		}
	}
}

func (n *Node) String() string {
	return fmt.Sprintf(
		"port : %s, status : %s, state : %d",
		n.Port,
		n.Status.String(),
		n.State)
}

func NewNode(port string) *Node {
	server := Server{}
	shutdown := make(chan struct{})
	node := &Node{
		Port:   port,
		Status: Candidate,
		State:  -1,
		Votes:  []string{},

		voted:     false,
		sd:        shutdown,
		heartBeat: make(chan struct{}),
		vote:      make(chan struct{}, 10),
	}

	handler := NewDefaultHandler(node)

	go server.Run(port, handler, shutdown)
	go node.run()

	return node
}
