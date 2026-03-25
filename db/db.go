package db

import (
	"embed"
	"fmt"
	"io/fs"
	"log"

	"database/sql"

	_ "github.com/lib/pq"

	"websocket-server/config"
)

//go:embed sql
var all embed.FS

var SqlFS fs.FS

var _connection *sql.DB

// Builtin module init
func init() {
    sub, err := fs.Sub(all, "sql")
    if err != nil {
        panic(err)
    }
    SqlFS = sub
}

// Exported custom init
func Init(module string, flagCreateAdmin bool) error {
	log.Printf("Initializing database (%s)...", config.Runtime.Database)

    sql, err := LoadSQL("init.sql")

    if err != nil {
        log.Fatal(err)
    }

    //fmt.Println(sql)

    _, err = DB_GetConnection().Exec(sql)

    if err != nil {
        return err
    }

    if flagCreateAdmin {
		log.Print("Creating Administrator user...")

        if err := CreateAdmin(); err != nil {
            return err
        }
    }
	
	return nil
}

func LoadSQL(filename string) (string, error) {
    folder := "postgres"
	if config.Local() { folder = "local" }

	path := fmt.Sprintf("%s/%s", folder, filename)

    content, err := fs.ReadFile(SqlFS, path)

    if err != nil {
        return "", err
    }

    sql := string(content)

	return sql, nil
}

func Close(db *sql.DB) {
	if db != nil {
		db.Close()
	}
}

func DB_GetConnection() *sql.DB {
	return _connection
}

func DB_Ping() bool {
	conn := DB_GetConnection()

	err := conn.Ping()

	return err == nil
}

func Open(app string) error {
	if config.Local() {
		config.Runtime.Database = "Local"
		config.Runtime.Local = true

		if app == "" {
			app = "ekhoes"
		}

		log.Printf("Opening local database '%s'... ", app)

		conn, err := openLocal(app)
		_connection = conn

		return err
	} else {
		if config.Conf.DB.Enabled {
			config.Runtime.Database = "External"

			_, err := ConnectAndKeepAlive()

			return err
		}

		if config.Conf.Redis.Enabled {
			_, err := RedisConnect()
			
			return err
		}
	}
	
	return nil
}

func OpenAndInit(app string, flagCreateAdmin bool) error {
	err := Open(app)
	
	if err != nil {
		return err
	}
	
	err = Init(app, flagCreateAdmin)
	
	return err
}

func CreateAdmin() error {
	if config.Local() {
		return createAdminLocal()
	}

	return nil
}
