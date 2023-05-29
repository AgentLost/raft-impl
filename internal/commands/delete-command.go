package commands

import (
	"raft-impl/internal/nodes"
)

type DeleteNodeCommand struct {
	name string
}

func (d *DeleteNodeCommand) Handle(data []byte) {
	nodes.NodeCluster.DeleteNode(string(data))
}

func (d *DeleteNodeCommand) GetName() string {
	return d.name
}

func NewDeleteNodeCommand() *DeleteNodeCommand {
	return &DeleteNodeCommand{
		name: "delete",
	}
}
