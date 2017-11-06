package dbLogger

import (
  _ "github.com/mattn/go-sqlite3"
  "lib/heartbeat"
  "database/sql"
  "log"
  "time"
  //_ "github.com/go-sql-driver/mysql"
)

func checkError(err error) {
  if err != nil {
    log.Panicf("[Error] %v\n", err)
  }
}

func InitializeDB(fn string) (*sql.DB) {
  db, err := sql.Open("sqlite3", fn)
  checkError(err)
  CreateTables(db)
  return db
}

func CreateTables(db *sql.DB) {
	usersTable := `
	CREATE TABLE IF NOT EXISTS users(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
		email VARCHAR(255) NOT NULL,
		passwordHash CHAR(60) NOT NULL,
    CONSTRAINT UNIQ_email UNIQUE (email)
	);
	`
  clientsTable := `
	CREATE TABLE IF NOT EXISTS clients(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
		user_id INTEGER NOT NULL,
		client_name VARCHAR(100) NOT NULL,
    CONSTRAINT UNIQ_ts_client_id UNIQUE (user_id, client_name),
    CONSTRAINT FK_UserClient FOREIGN KEY (user_id) REFERENCES users (id)
	);
	`
  heartbeatsTable := `
	CREATE TABLE IF NOT EXISTS heartbeats(
    id INTEGER PRIMARY KEY AUTOINCREMENT,
		ts timestamp NOT NULL,
		client_id INTEGER NOT NULL,
    CONSTRAINT UNIQ_ts_client_id UNIQUE (ts, client_id)
    CONSTRAINT FK_ClientHeartbeat FOREIGN KEY (client_id) REFERENCES clients (id)
	);
	`

	_, err := db.Exec(usersTable)
	checkError(err)
	_, err = db.Exec(clientsTable)
	checkError(err)
	_, err = db.Exec(heartbeatsTable)
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
  db := InitializeDB("server/db.sqlite3")
	rows, err := db.Query("SELECT client_id, ts FROM heartbeats")
  checkError(err)

	return fetchHeartbeats(rows)
}

func GetHeartbeats(clientID string) []*heartbeat.Heartbeat {
  db := InitializeDB("server/db.sqlite3")
	rows, err := db.Query("SELECT client_id, ts FROM heartbeats WHERE client_id=?", clientID)
  checkError(err)

  return fetchHeartbeats(rows)
}
