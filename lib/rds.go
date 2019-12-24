package lib

import (
	"../common"
	"database/sql"
	"fmt"
	_ "github.com/denisenkom/go-mssqldb"
	_ "github.com/go-sql-driver/mysql"
	"strings"
)

func MysqlQuery(sqlcommand string) error {
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config["User"], Config["Passwd"], Config["Server"], Config["Port"], Config["Database"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}

	rows, err := db.Query(sqlcommand)
	if err != nil {
		return fmt.Errorf("Query failed:", err.Error())
	}

	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
		fmt.Print(cols[i])
		fmt.Print("\t")
	}

	for rows.Next() {
		err = rows.Scan(colsdata...)
		common.PrintRow(colsdata)
	}
	return nil
}

func MysqlInsert(sqlcommand string) error {
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config["User"], Config["Passwd"], Config["Server"], Config["Port"], Config["Database"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}

	for _, command := range strings.Split(sqlcommand, ";;;") {

		if strings.TrimSpace(command) == "" {
			continue
		}
		_, err := db.Exec(command)
		if err != nil {
			return fmt.Errorf("Exec failed:", err.Error())
		}
	}
	fmt.Println("Exec success")
	return nil
}

func MysqlGetTableName() error {
	connstring := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", Config["User"], Config["Passwd"], Config["Server"], Config["Port"], Config["Database"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	for {
		tablename := common.GetRandomString(10)
		_, err = db.Query("select 1 from " + tablename + " limit 1")
		if err != nil {
			Config["TableName"] = tablename
			Config["Temp"] = Config["Temp"].(string) + "RdsName:" + Config["TableName"].(string) + "\n"
			return nil
		}
	}

}

func MssqlQuery(sqlcommand string) error {
	connstring := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", Config["Server"], Config["Port"], Config["Database"], Config["User"], Config["Passwd"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		fmt.Println("Open Connection failed:", err.Error())
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	rows, err := db.Query(sqlcommand)
	if err != nil {
		return fmt.Errorf("Exec failed:", err.Error())
	}
	cols, err := rows.Columns()
	var colsdata = make([]interface{}, len(cols))
	for i := 0; i < len(cols); i++ {
		colsdata[i] = new(interface{})
		fmt.Print(cols[i])
		fmt.Print("\t")
	}
	fmt.Println()

	for rows.Next() {
		err = rows.Scan(colsdata...)
		common.PrintRow(colsdata)
	}
	return nil

}

func MssqlInsert(sqlcommand string) error {
	connstring := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", Config["Server"], Config["Port"], Config["Database"], Config["User"], Config["Passwd"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	for _, command := range strings.Split(sqlcommand, ";;;") {
		if strings.TrimSpace(command) == "" {
			continue
		}
		_, err := db.Exec(command)
		if err != nil {
			return fmt.Errorf("Exec failed:", err.Error())
		}
	}
	fmt.Println("Exec success")
	return nil
}

func MssqlGetTableName() error {
	connstring := fmt.Sprintf("server=%s;port%d;database=%s;user id=%s;password=%s", Config["Server"], Config["Port"], Config["Database"], Config["User"], Config["Passwd"])
	db, err := sql.Open(Config["DatabaseType"].(string), connstring)
	if err != nil {
		return fmt.Errorf("Open Connection failed:", err.Error())
	}
	for {
		tablename := common.GetRandomString(10)
		_, err = db.Query("select 1 from " + tablename + " limit 1")
		if err != nil {
			Config["TableName"] = tablename
			Config["Temp"] = Config["Temp"].(string) + "RdsName:" + Config["TableName"].(string) + "\n"
			return nil
		}
	}

}

func RdsInsert(sqlcommand string) error {
	switch Config["DatabaseType"] {
	case "mysql":
		{
			err := MysqlInsert(sqlcommand)
			return err
		}
	case "mssql":
		{
			err := MssqlInsert(sqlcommand)
			return err
		}
	default:
		{
			return fmt.Errorf("rdstype only enter mysql or mssql")
		}
	}

}

func RdsQuery(sqlcommand string) error {
	switch Config["DatabaseType"] {
	case "mysql":
		{
			err := MysqlQuery(sqlcommand)
			return err
		}
	case "mssql":
		{
			err := MssqlQuery(sqlcommand)
			return err
		}
	default:
		{
			return fmt.Errorf("rdstype only enter mysql or mssql")
		}
	}
}

func RdsGetTableName() error {
	switch Config["DatabaseType"] {
	case "mysql":
		{
			err := MysqlGetTableName()
			return err
		}
	case "mssql":
		{
			err := MssqlGetTableName()
			return err
		}
	default:
		{
			return fmt.Errorf("rdstype only enter mysql or mssql")
		}
	}
}
