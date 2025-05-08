package database

import (
	"database/sql"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strings"

	_ "github.com/lib/pq"
)

func (d *Database) Deploy(listFile string) (err error) {

	// let's get absolute path for the list file
	if listFile, err = filepath.Abs(listFile); err != nil {
		return err
	}

	// let's get the directory fo rthe database, so we could use it with relative paths in the list file
	databaseDirectory := filepath.Dir(filepath.Dir(listFile))
	databaseDirectoryRoot := filepath.Dir(databaseDirectory)
	databaseDirectoryRootLength := len(databaseDirectoryRoot)

	// let's deploy the list file
	data, err := os.ReadFile(listFile)
	if err != nil {
		return err
	}

	// let's iterate each line in the file
	for _, item := range strings.Split(string(data), "\n") {
		// drop the comment
		if p := strings.Index(item, `#`); p >= 0 {
			item = item[0:p]
		}
		// trim the line
		item = strings.TrimSpace(item)
		// do we still have a line?
		if len(item) > 0 {
			// deploy file
			filename := path.Join(databaseDirectory, item)
			fmt.Printf("...%s\n", filename[databaseDirectoryRootLength:])
			if err := d.deployFile(filename); err != nil {
				return err
			}
		}
	}

	return nil

}

func (d *Database) deployFile(filename string) error {

	query, err := os.ReadFile(filename)
	if err != nil {
		return err
	}

	// opern db connection...
	db, err := sql.Open(dbDriverName, d.ConnectionString)
	if err != nil {
		return err
	}
	defer db.Close()

	// ...run query from file
	if _, err := db.Exec(string(query)); err != nil {
		return err
	}

	return nil

}
