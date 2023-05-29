package app

import (
	"bufio"
	"log"
	"os"
	"raft-impl/internal/commands"
	"raft-impl/internal/nodes"
	"strings"
)

func Run() {
	manager := commands.NewManager()

	manager.AddCommand(commands.NewAddNodeCommand())
	manager.AddCommand(commands.NewDeleteNodeCommand())
	manager.AddCommand(commands.NewGetInfoCommand())

	nodes.InitCluster()

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		data := strings.Split(scanner.Text(), " ")

		command, ok := manager.GetCommand(data[0])
		if !ok {
			log.Println("bad command")
			continue
		}

		command.Handle([]byte(data[1]))
	}
}
