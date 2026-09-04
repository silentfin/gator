package main

import (
	"fmt"
	"os"

	"github.com/silentfin/gator/internal/config"
)

func main() {
	var s state
	userConfig := config.Read()
	s.conf = &userConfig
	availableCommands := commands{cmds: map[string]func(*state, command) error{}}
	availableCommands.register("login", handlerLogin)
	userArgs := os.Args
	if len(userArgs) < 2 {
		fmt.Println("insufficient args provided")
		os.Exit(1)
	} else {
		cmdName := os.Args[1]
		cmdArgs := os.Args[2:]
		err := availableCommands.run(&s, command{name: cmdName, args: cmdArgs})
		if err != nil {
			fmt.Printf("error occured: %s\n", err)
			os.Exit(1)
		} else {
			os.Exit(0)
		}
	}
}
