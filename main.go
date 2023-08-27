package main

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/pietjan/tm/app/adapters/sqlite"
	"github.com/pietjan/tm/app/ports/cli"

	"github.com/kirsle/configdir"
	_ "modernc.org/sqlite"
)

func main() {
	configPath := configdir.LocalConfig(`tm`)
	if err := configdir.MakePath(configPath); err != nil {
		panic(err)
	}

	repo := sqlite.New(fmt.Sprintf(`file:%s`, filepath.Join(configPath, `db.sqlite`)))
	app := cli.New(repo)

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
