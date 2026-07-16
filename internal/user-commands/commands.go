package user_commands

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

//Commands
/*
	add task -> description
	update task -> id, newStatus
*/

type CommandType string

const stopCommand CommandType = "stop"

type CommandHandler func(arg []string) error
type Command struct {
	commandName CommandType
	execute     func()
}

type UserCommands struct {
	commands map[CommandType]CommandHandler
}

func (u *UserCommands) Init() {
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("Enter command: \n > ")
		userInput, _ := reader.ReadString('\n')
		splittedCommand := strings.Split(strings.TrimSpace(userInput), " ")

		if len(splittedCommand) != 0 {
			command := CommandType(splittedCommand[0])

			commandArgs := []string{strings.Join(splittedCommand[1:], " ")}

			err := u.Run(command, commandArgs)

			if err != nil {
				fmt.Println(err)
			}

			if command == stopCommand {
				return
			}
			continue
		}

		fmt.Println("Command cannot be empty. Try again:")
	}
}

func (u *UserCommands) Run(command CommandType, args []string) error {
	if _, ok := u.commands[command]; ok {
		err := u.commands[command](args)

		return err
	}

	return errors.New(fmt.Sprintf("Command %s not found", command))
}

func (u *UserCommands) Register(command CommandType, execute CommandHandler) {
	u.commands[command] = execute
}

func NewUserCommands() UserCommands {
	userCommands := UserCommands{
		commands: make(map[CommandType]CommandHandler),
	}

	return userCommands
}
