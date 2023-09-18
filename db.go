package main

import (
	"database/sql"
	"fmt"
	"reflect"
	"runtime"
	"time"
)

func createUsersTable(sqlCon *sql.DB) {
	_, err := sqlCon.Exec("CREATE TABLE IF NOT EXISTS users (" +
		"id int PRIMARY KEY AUTO_INCREMENT, " +
		"name VARCHAR(100) COLLATE utf8_bin UNIQUE, " +
		"password CHAR(64), " +
		"userPower INT, " +
		"createdAt DATETIME(3)" +
		")")
	if err != nil {
		HandleUnexpectedError(err, createUsersTable, "error in create query")
	}
}

func createDisconnectsTable(sqlCon *sql.DB) {
	_, err := sqlCon.Exec("CREATE TABLE IF NOT EXISTS disconnects (" +
		"connectsId int AUTO_INCREMENT PRIMARY KEY," +
		"createdAt DATETIME(3)," +
		"FOREIGN KEY (connectsId) REFERENCES connects(id)" +
		")")
	if err != nil {
		HandleUnexpectedError(err, createDisconnectsTable, "error in create query")
	}
}
func createConnectsTable(sqlCon *sql.DB) {
	_, err := sqlCon.Exec("CREATE TABLE IF NOT EXISTS connects (" +
		"id int AUTO_INCREMENT PRIMARY KEY," +
		"userId int," +
		"ip VARCHAR(45)," +
		"createdAt DATETIME(3)," +
		"FOREIGN KEY (userId) REFERENCES users(id)" +
		")")
	if err != nil {
		HandleUnexpectedError(err, createConnectsTable, "error in create query")
	}
}

func createActivitiesTable(sqlCon *sql.DB) {
	_, err := sqlCon.Exec("CREATE TABLE IF NOT EXISTS activities (" +
		"id INT AUTO_INCREMENT PRIMARY KEY," +
		"connectsId INT," +
		"messageType varchar(100)," +
		"messageArgs mediumtext," +
		"receivedAt DATETIME(3)," +
		"FOREIGN KEY (connectsId) REFERENCES connects(id)" +
		")")
	if err != nil {
		HandleUnexpectedError(err, createActivitiesTable, "error in create query")
	}
}

func HandleUnexpectedError(err error, i interface{}, customMsg string) {
	if err != nil {
		panic("\"" + runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name() + "\" failed at " + time.Now().Format(time.RFC3339) + "\n" + err.Error() + "\n" + customMsg)
	}
}

func ConnectToLocalDb(username string, password string, dbName string) *sql.DB {
	sqlCon, err := sql.Open("mysql", username+":"+password+"@/"+dbName+"?parseTime=true")
	if err != nil {
		HandleUnexpectedError(err, ConnectToLocalDb, "error in sql.Open")
	}
	if !pingMysqlConnection(sqlCon) {
		HandleUnexpectedError(err, ConnectToLocalDb, "error in PingMysqlConnection")
	}
	return sqlCon
}

func pingMysqlConnection(sqlCon *sql.DB) bool {
	if err := sqlCon.Ping(); err != nil {
		fmt.Println(err.Error())
		return false
	}
	return true
}

func SelectUserIdByUsernameAndPassword(sqlConn *sql.DB, username string, password string) int {
	rows, err := sqlConn.Query("SELECT id "+
		"FROM users "+
		"WHERE name = ? AND password = ?",
		username, password)
	if err != nil {
		HandleUnexpectedError(err, SelectUserIdByUsernameAndPassword, "error in select query")
	}
	var userId int = 0
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&userId)
		if err != nil {
			HandleUnexpectedError(err, SelectUserIdByUsernameAndPassword, "error in scan")
		}
	}
	return userId
}

func InsertUser(sqlConn *sql.DB, name string, password string, userPower int, createdAt time.Time) int {

	res, err := sqlConn.Exec("INSERT INTO users ("+
		"name, password, userPower, createdAt) VALUES (?, ?, ?, ?)",
		name, password, userPower, TimeToSqlString(createdAt))
	if err != nil {
		HandleUnexpectedError(err, InsertUser, "error in insert query")
	}
	userId, err := res.LastInsertId()
	if err != nil {
		HandleUnexpectedError(err, InsertUser, "error in last insert id")
	}
	return int(userId)
}
func TimeToSqlString(t time.Time) string {
	return t.Format("2006-01-02 15:04:05.000")
}

func InsertConnect(sqlConn *sql.DB, userId int, ip string, createdAt time.Time) int {
	res, err := sqlConn.Exec("INSERT INTO connects ("+
		"userId, ip, createdAt"+
		") VALUES (?, ?, ?)",
		userId, ip, TimeToSqlString(createdAt))
	if err != nil {
		HandleUnexpectedError(err, InsertConnect, "error in insert query")
	}
	id, err := res.LastInsertId()
	if err != nil {
		HandleUnexpectedError(err, InsertConnect, "error in last insert id")
	}
	return int(id)
}

func InsertDisconnect(sqlConn *sql.DB, connnectsId int, createdAt time.Time) int {
	res, err := sqlConn.Exec("INSERT INTO disconnects ("+
		"connectsId, createdAt"+
		") VALUES (?, ?)",
		connnectsId, TimeToSqlString(createdAt))
	if err != nil {
		HandleUnexpectedError(err, InsertDisconnect, "error in insert query")
	}
	id, err := res.LastInsertId()
	if err != nil {
		HandleUnexpectedError(err, InsertDisconnect, "error in last insert id")
	}
	return int(id)
}

func InsertActivity(sqlConn *sql.DB, connnectsId int, messageType string, messageArgs string, receivedAt time.Time) {
	_, err := sqlConn.Exec("INSERT INTO activities ("+
		"connectsId, messageType, messageArgs, receivedAt"+
		") VALUES (?, ?, ?, ?)",
		connnectsId, messageType, messageArgs, TimeToSqlString(receivedAt))
	if err != nil {
		HandleUnexpectedError(err, InsertActivity, "error in insert query")
	}
}

func SelectUsernameByUserId(sqlConn *sql.DB, userId int) string {
	rows, err := sqlConn.Query("SELECT name "+
		"FROM users "+
		"WHERE id = ?",
		userId)
	if err != nil {
		HandleUnexpectedError(err, SelectUsernameByUserId, "error in select query")
	}
	var username string = ""
	defer rows.Close()
	if rows.Next() {
		err = rows.Scan(&username)
		if err != nil {
			HandleUnexpectedError(err, SelectUsernameByUserId, "error in scan")
		}
	}
	return username
}
