package dbLogger

import (
  "../udpUtils"
  "database/sql"
  _ "github.com/mattn/go-sqlite3"
  "log"
  //_ "github.com/go-sql-driver/mysql"
)

func InitializeDB(fn string) (*sql.DB) {
  db, err := sql.Open("sqlite3", fn)
  if err != nil {
    panic(err)
  } else if db == nil {
    panic("DB could not be initialized")
  }
  CreateTables(db)
  return db
}

func CreateTables(db *sql.DB) {
  //TODO: User management, client_id shouldn't be arbitrary, but a foreign
  //key to a table of clients. Heartbeats that are not authorized should be
  //rejected.

	heartbeatsTable := `
	CREATE TABLE IF NOT EXISTS heartbeats(
		ts timestamp NOT NULL PRIMARY KEY,
		client_id VARCHAR(100) NOT NULL
	);
	`
	_, err := db.Exec(heartbeatsTable)
	if err != nil { panic(err) }
}

func LogHeartbeat(db *sql.DB, heartbeat *udpUtils.Heartbeat) {
  tx, err := db.Begin()
  if err != nil {
    log.Fatal(err)
  }
  stmt, err := tx.Prepare("INSERT INTO heartbeats (client_id, ts) VALUES (?, ?)")
  if err != nil {
    log.Fatal(err)
  }
  defer stmt.Close()
  _, err = stmt.Exec(heartbeat.ClientID, heartbeat.Timestamp)
  if err != nil {
    log.Fatal(err)
  }
  tx.Commit()
  stmt.Close()
}
