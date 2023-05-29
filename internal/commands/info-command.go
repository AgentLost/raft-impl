package commands

import (
	"log"
	"raft-impl/internal/nodes"
)

type GetInfoCommand struct {
	name string
}

func (g *GetInfoCommand) Handle(_ []byte) {
	log.Println(nodes.NodeCluster.GetInfo())
}

func (g *GetInfoCommand) GetName() string {
	return g.name
}

func NewGetInfoCommand() *GetInfoCommand {
	return &GetInfoCommand{
		name: "info",
	}
}
