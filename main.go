package main

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq"
	"github.com/silentfin/gator/internal/config"
	"github.com/silentfin/gator/internal/database"
)

func main() {
	var s state
	userConfig := config.Read()
	s.conf = &userConfig
	db, err := sql.Open("postgres", s.conf.DbUrl)
	if err != nil {
		fmt.Printf("error occured: %v\n", &err)
	}
	dbQueries := database.New(db)
	s.db = dbQueries

	availableCommands := commands{cmds: map[string]func(*state, command) error{}}
	availableCommands.register("login", handlerLogin)
	availableCommands.register("register", handlerRegister)
	availableCommands.register("reset", handleReset)
	availableCommands.register("users", handleUsers)
	availableCommands.register("agg", handleAgg)
	availableCommands.register("addfeed", middlewareLoggedIn(handleAddFeed))
	availableCommands.register("feeds", handleFeeds)
	availableCommands.register("follow", middlewareLoggedIn(handleFollow))
	availableCommands.register("following", middlewareLoggedIn(handleFollowing))
	availableCommands.register("unfollow", middlewareLoggedIn(handleUnfollow))
	availableCommands.register("browse", middlewareLoggedIn(handleBrowse))

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
