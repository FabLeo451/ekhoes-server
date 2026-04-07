package herenow

import (
	"ekhoes-server/db"
	"ekhoes-server/utils"
)

func Install() error {
	utils.Log(thisModule, "Opening database...")

	if err := db.OpenDatabase(); err != nil {
		return err
	}

	utils.Log(thisModule, "Creating schema...")

	err := db.ExecuteSQL(SqlFS, "install.sql")

	if err != nil {
		return err
	}

	db.CloseDatabase()

	return nil
}
