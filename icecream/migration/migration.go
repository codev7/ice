package migration

import (
	"fmt"
	"io/ioutil"
	"log"
	"path"
	"time"

	"github.com/nirandas/ice"
)

var all bool

func Process(args []string) {
	if ice.Config.MigrationPath == "" {
		fmt.Println("please set the MigrationPath configuration in the ice config file")
		return
	}
	if len(args) == 0 {
		fmt.Println("Specify command arguments. Available commands are, make, up, down")
		return
	}

	for _, a := range args {
		if a == "-all" {
			all = true
		}
	}

	//cmd:make
	if args[0] == "make" {
		handleMake(args)
		return
	}

	//cmd:up
	if args[0] == "up" {
		handleUp(args)
		return
	}

	//cmd:down
	if args[0] == "down" {
		handleDown(args)
		return
	}

}

func handleMake(args []string) {
	newPath := path.Join(ice.Config.MigrationPath, fmt.Sprintf("%s-%s.sql", time.Now().Format("20060102-150405"), args[1]))
	err := ioutil.WriteFile(newPath, []byte(`--icecream:migration:up

--icecream:migration:down

`), 0)
	if err != nil {
		fmt.Println("Failed to write the migration ", err)
		return
	}
	fmt.Println("Migration file created at ", newPath)
}

func handleUp(args []string) {
	db := OpenDb()
	rec := ParseMigrationRecords(db)
	lm := LastMigration(rec)
	_, n := ListMigrations(lm.Migration)
	if len(n) > 0 {

		for len(n) > 0 {
			mc := ParseMigration(n[0])
			fmt.Println("Migrating ", mc.Name)
			if err := mc.RunUp(db, lm.Batch+1); err != nil {
				log.Fatalf("Failed %s", err.Error())
			}
			fmt.Println("Completed ", mc.Name)
			n = n[1:]
			if !all {
				break
			}
		}
		fmt.Println(len(n), " migrations remaining")
	} else {
		fmt.Println("No more migrations to run")
	}
}

func handleDown(args []string) {
	db := OpenDb()
	rec := ParseMigrationRecords(db)
	lm := LastMigration(rec)
	p, _ := ListMigrations(lm.Migration)
	if len(p) > 0 {
		for len(p) > 0 {
			mc := ParseMigration(p[len(p)-1])
			fmt.Println("Downgrading ", mc.Name)
			if err := mc.RunDown(db); err != nil {
				log.Fatalf("Failed %s", err.Error())
			}
			fmt.Println("Rolled back", mc.Name)
			p = p[:len(p)-1]
			if !all {
				break
			}
		}
		fmt.Println(len(p), " migrations remaining to be rolled back")
	} else {
		fmt.Println("No migrations to rollback")
	}
}
