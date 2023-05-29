package commands

import (
	"log"
	"raft-impl/internal/nodes"
)

type AddNodeCommand struct {
	name string
}

func (a *AddNodeCommand) Handle(data []byte) {
	err := nodes.NodeCluster.AddNode(string(data))
	if err != nil {
		log.Println(err)
		return
	}
}

func (a *AddNodeCommand) GetName() string {
	return a.name
}

func NewAddNodeCommand() *AddNodeCommand {
	return &AddNodeCommand{
		name: "add",
	}
}
