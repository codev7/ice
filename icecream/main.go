package main

import (
	"fmt"
	"github.com/nirandas/ice"
	"github.com/nirandas/ice/icecream/migration"
	"github.com/nirandas/ice/icecream/routes"
	"os"
)

func main() {
	fmt.Println("ICE version 0.01")
	ice.LoadConfig()
	fmt.Println(ice.Config.Name, " version ", ice.Config.Version)

	if len(os.Args) == 1 {
		fmt.Println("supported commands: 'migration'")
		return
	}

	if os.Args[1] == "migration" {
		migration.Process(os.Args[2:])
		return
	}

	if os.Args[1] == "routes" {
		routes.Process(os.Args[2:])
		return
	}

	fmt.Println("Cannot understand the command")
}
