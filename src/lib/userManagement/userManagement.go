package userManagement

import (
  "log"
  "golang.org/x/crypto/bcrypt"
  "lib/dbLogger"
)

func checkError(err error) {
  if err != nil {
    log.Panicf("[Error] %v\n", err)
  }
}

func GeneratePasswordHash(password string) string {
  hash, _ := bcrypt.GenerateFromPassword([]byte(password), 10)
  return string(hash)
}

func GetPasswordHash(email string) (passwordHash string) {
  db := dbLogger.InitializeDB("server/db.sqlite3")
	rows, err := db.Query("SELECT passwordHash FROM users WHERE email=?", email)
  checkError(err)

  defer rows.Close()
  rows.Next()
  rows.Scan(&passwordHash)
  return
}

func Verify(email string, password string) bool {
  if bcrypt.CompareHashAndPassword([]byte(GetPasswordHash(email)), []byte(password)) == nil {
    return true
  } else {
    return false
  }
}

func CreateUser(email string, password string) {
  db := dbLogger.InitializeDB("server/db.sqlite3")
  tx, err := db.Begin()
  checkError(err)

  stmt, err := tx.Prepare("INSERT INTO users (email, passwordHash) VALUES (?, ?)")
  checkError(err)

  _, err = stmt.Exec(email, GeneratePasswordHash(password))
  checkError(err)

  tx.Commit()
}
