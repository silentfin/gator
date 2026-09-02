package main

import (
	"fmt"

	"github.com/silentfin/gator/internal/config"
)

func main() {
	userConfig := config.Read()
	fmt.Println(userConfig.CurrentUserName)
	fmt.Println(userConfig.DbUrl)
	if err := userConfig.SetUser("geralt"); err != nil {
		fmt.Printf("Something went wrong: %v", err)
	} else {
		updatedConfig := config.Read()
		fmt.Println(updatedConfig.CurrentUserName)
		fmt.Println(updatedConfig.DbUrl)
	}
}
