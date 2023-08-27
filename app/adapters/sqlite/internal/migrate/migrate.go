package migrate

import (
	"database/sql"
	"embed"
	"io/fs"

	"github.com/pietjan/migrate"
	"github.com/pietjan/migrate/database/sqlite"
	"github.com/pietjan/migrate/source/file"
)

//go:embed files/*.sql
var assets embed.FS

func Run(db *sql.DB) error {
	fsys, err := fs.Sub(assets, `files`)
	if err != nil {
		return err
	}

	migrate := migrate.New(
		migrate.FromFile(file.FS(fsys)),
		migrate.ToSqlite(sqlite.DB(db)),
	)

	return migrate.Run()
}
