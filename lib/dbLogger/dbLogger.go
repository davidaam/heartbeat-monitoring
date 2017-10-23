package dbLogger

import (
  "../heartbeat"
  "database/sql"
  "log"
  "time"
  _ "github.com/mattn/go-sqlite3"
  //_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
  if err != nil {
    log.Fatalln("[Error] %v", err)
  }
}

func InitializeDB(fn string) (*sql.DB) {
  db, err := sql.Open("sqlite3", fn)

  checkError(err)
  if db == nil {
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
    id INTEGER PRIMARY KEY AUTOINCREMENT,
		ts timestamp NOT NULL,
		client_id VARCHAR(100) NOT NULL,
    CONSTRAINT UNIQ_ts_client_id UNIQUE (ts, client_id)
	);
	`
	_, err := db.Exec(heartbeatsTable)
	checkError(err)
}

func LogHeartbeats(db *sql.DB, heartbeats chan *heartbeat.Heartbeat) {
  for {
    heartbeat := <-heartbeats

    tx, err := db.Begin()
    checkError(err)

    stmt, err := tx.Prepare("INSERT INTO heartbeats (client_id, ts) VALUES (?, ?)")
    checkError(err)

    _, err = stmt.Exec(heartbeat.ClientID, heartbeat.Timestamp)
    checkError(err)

    tx.Commit()
  }
}

func fetchHeartbeats(rows *sql.Rows) []*heartbeat.Heartbeat {
  defer rows.Close()

  // TODO: Check if row count can be retrieved, to allocate exact memory beforehand
  var heartbeats []*heartbeat.Heartbeat

	for rows.Next() {
		var clientID string
		var timestamp time.Time
		checkError(rows.Scan(&clientID, &timestamp))
    hb := &heartbeat.Heartbeat { ClientID: clientID, Timestamp: timestamp.Unix() }
    heartbeats = append(heartbeats, hb)
	}
  checkError(rows.Err())

  return heartbeats
}

func ListHeartbeats() []*heartbeat.Heartbeat {
  db := InitializeDB("../../server/heartbeats.sqlite3")
	rows, err := db.Query("SELECT client_id, ts FROM heartbeats")
  checkError(err)

	return fetchHeartbeats(rows)
}

func GetHeartbeats(clientID string) []*heartbeat.Heartbeat {
  db := InitializeDB("../../server/heartbeats.sqlite3")
	rows, err := db.Query("SELECT client_id, ts FROM heartbeats WHERE client_id=?", clientID)
  checkError(err)

  return fetchHeartbeats(rows)
}
