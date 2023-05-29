package nodes

import (
	"io"
	"log"
	"net/http"
)

type Handler struct {
	node *Node
}

func (h *Handler) ServeHTTP(wr http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/heart-beat" && r.Method == http.MethodGet {
		h.node.heartBeat <- struct{}{}

		wr.WriteHeader(http.StatusOK)
		wr.Write([]byte("ok"))

		return
	} else if r.URL.Path == "/vote" && r.Method == http.MethodPost {
		h.node.vote <- struct{}{}

		if h.node.voted == false {
			h.node.voted = true

			data, err := io.ReadAll(r.Body)
			if err != nil {
				log.Println(err)
			}

			h.node.Votes = append(h.node.Votes, string(data))
		}

		wr.Write([]byte(h.node.Votes[len(h.node.Votes)-1]))
		return
	}

	wr.WriteHeader(http.StatusNotFound)
}

func NewDefaultHandler(node *Node) *Handler {
	return &Handler{
		node: node,
	}
}
