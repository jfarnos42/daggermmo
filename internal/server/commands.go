package server

import (
    "errors"
    "strings"
)

type CommandFunc func(args []string) (string, error)

type CommandManager struct {
    commands map[string]CommandFunc
}

func NewCommandManager() *CommandManager {
    return &CommandManager{
        commands: make(map[string]CommandFunc),
    }
}

func (cm *CommandManager) Register(name string, fn CommandFunc) {
    cm.commands[strings.ToLower(name)] = fn
}

func (cm *CommandManager) Execute(input string) (string, error) {
    parts := strings.Fields(input)
    if len(parts) == 0 {
        return "", errors.New("no command provided")
    }

    cmdName := strings.ToLower(parts[0])
    args := parts[1:]

    cmd, ok := cm.commands[cmdName]
    if !ok {
        return "", errors.New("unknown command: " + cmdName)
    }

    return cmd(args)
}
