package commands

import "sync"

type CommandManager struct {
	commands map[string]Command
	mu       *sync.RWMutex
}

func (cm *CommandManager) AddCommand(command Command) {
	cm.mu.Lock()
	defer cm.mu.Unlock()
	cm.commands[command.GetName()] = command
}

func (cm *CommandManager) GetCommand(name string) (Command, bool) {
	cm.mu.RLock()
	c, ok := cm.commands[name]
	cm.mu.RUnlock()
	return c, ok
}

func NewManager() *CommandManager {
	return &CommandManager{
		commands: make(map[string]Command),
		mu:       &sync.RWMutex{},
	}
}
