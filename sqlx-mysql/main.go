package main

import (
	"context"
	"database/sql"
	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

type Person struct {
	Id         int       `db:"Id"`
	Name       string    `db:"Name"`
	City       string    `db:"City"`
	AddTime    time.Time `db:"AddTime"`
	UpdateTime time.Time `db:"UpdateTime"`
}

func main() {
	db, err := sqlx.Connect("mysql", "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4&parseTime=true&loc=Local")
	if err != nil {
		log.Println(err)
		return
	}

	// Clear all data
	_, err = db.Exec("truncate table Person")
	if err != nil {
		log.Println(err)
		return
	}

	// Insert
	insertResult := db.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	lastInsertId, _ := insertResult.LastInsertId()
	log.Println("Insert Id is ", lastInsertId)

	insertPerson := &Person{
		Name:       "Li Si",
		City:       "Shanghai",
		AddTime:    time.Now(),
		UpdateTime: time.Now(),
	}
	insertPersonResult, err := db.NamedExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES(:Name, :City, :AddTime, :UpdateTime)", insertPerson)
	if err != nil {
		log.Println(err)
		return
	}
	lastInsertPersonId, _ := insertPersonResult.LastInsertId()
	log.Println("InsertPerson Id is ", lastInsertPersonId)

	insertMap := map[string]interface{}{
		"n": "Wang Wu",
		"c": "HongKong",
		"a": time.Now(),
		"u": time.Now(),
	}
	insertMapResult, err := db.NamedExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES(:n, :c, :a, :u)", insertMap)
	if err != nil {
		log.Println(err)
		return
	}
	lastInsertMapId, _ := insertMapResult.LastInsertId()
	log.Println("InsertMap Id is ", lastInsertMapId)

	insertPersonArray := []Person{
		{Name: "BOSIMA", City: "Wu Han", AddTime: time.Now(), UpdateTime: time.Now()},
		{Name: "BOSSMA", City: "Xi An", AddTime: time.Now(), UpdateTime: time.Now()},
		{Name: "BOMA", City: "Cheng Du", AddTime: time.Now(), UpdateTime: time.Now()},
	}
	insertPersonArrayResult, err := db.NamedExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES(:Name, :City, :AddTime, :UpdateTime)", insertPersonArray)
	if err != nil {
		log.Println(err)
		return
	}
	insertPersonArrayId, _ := insertPersonArrayResult.LastInsertId()
	log.Println("InsertPersonArray Id is ", insertPersonArrayId)

	// Query
	row := db.QueryRowx("select * from Person where Name=?", "Zhang San")
	if row.Err() == sql.ErrNoRows {
		log.Println("Not found Zhang San")
	} else {
		rowMap := make(map[string]interface{})
		err = row.MapScan(rowMap)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("QueryRowx-MapScan:", string(rowMap["City"].([]byte)))
	}

	row = db.QueryRowx("select * from Person where Name=?", "Zhang San")
	if row.Err() == sql.ErrNoRows {
		log.Println("Not found Zhang San")
	} else {
		queryPerson := &Person{}
		err = row.StructScan(queryPerson)
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("QueryRowx-StructScan:", queryPerson.City)
	}

	row = db.QueryRowx("select * from Person where Name=?", "Zhang San")
	if row.Err() == sql.ErrNoRows {
		log.Println("Not found Zhang San")
	} else {
		rowSlice, err := row.SliceScan()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("QueryRowx-SliceScan:", string(rowSlice[2].([]byte)))
	}

	rows, err := db.Queryx("select * from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		rowSlice, err := rows.SliceScan()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("Queryx-SliceScan:", string(rowSlice[2].([]byte)))
	}

	rows, err = db.NamedQuery("select * from Person where Name=:n", map[string]interface{}{"n": "Zhang San"})
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		rowSlice, err := rows.SliceScan()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("NamedQuery-SliceScan:", string(rowSlice[2].([]byte)))
	}

	getPerson := &Person{}
	err = db.Get(getPerson, "select * from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Get:", getPerson.City)

	getId := new(int64)
	err = db.Get(getId, "select Id from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Get-Id:", *getId)

	selectPersons := []Person{}
	err = db.Select(&selectPersons, "select * from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	selectPerson := selectPersons[0]
	log.Println("Select:", selectPerson.City)

	selectTowFieldSlice := []Person{}
	err = db.Select(&selectTowFieldSlice, "select Id,Name from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	selectTwoField := selectTowFieldSlice[0]
	log.Println("Select-TowField:", selectTwoField.Id, selectTwoField.Name)

	selectNameSlice := []string{}
	err = db.Select(&selectNameSlice, "select Name from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	selectName := selectNameSlice[0]
	log.Println("Select-Name:", selectName)

	rows, err = db.NamedQuery("select * from Person where Name=:n", map[string]interface{}{"n": "Zhang San"})
	if err != nil {
		log.Println(err)
		return
	}
	for rows.Next() {
		rowSlice, err := rows.SliceScan()
		if err != nil {
			log.Println(err)
			return
		}
		log.Println("NamedQuery-SliceScan:", string(rowSlice[2].([]byte)))
	}

	// Update
	updateResult := db.MustExec("Update Person set City=?, UpdateTime=? where Id=?", "Shanghai", time.Now(), 1)
	log.Print("Update-MustExec:")
	log.Println(updateResult.RowsAffected())

	updateMapResult, err := db.NamedExec("Update Person set City=:City, UpdateTime=:UpdateTime where Id=:Id",
		map[string]interface{}{"City": "Chong Qing", "UpdateTime": time.Now(), "Id": 1})
	if err != nil {
		log.Println(err)
	}
	log.Print("Update-NamedExec:")
	log.Println(updateMapResult.RowsAffected())

	// Delete
	deleteResult := db.MustExec("Delete from Person where Id=?", 1)
	log.Print("Delete-MustExec:")
	log.Println(deleteResult.RowsAffected())

	deleteMapResult, err := db.NamedExec("Delete from Person where Id=:Id",
		map[string]interface{}{"Id": 1})
	if err != nil {
		log.Println(err)
		return
	}
	log.Print("Delete-NamedExec:")
	log.Println(deleteMapResult.RowsAffected())

	// Stmt
	bosima := Person{}
	bossma := Person{}

	nstmt, err := db.PrepareNamed("SELECT * FROM Person WHERE Name = :n")
	if err != nil {
		log.Println(err)
		return
	}
	err = nstmt.Get(&bossma, map[string]interface{}{"n": "BOSSMA"})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("NamedStmt-Get1:", bossma.City)
	err = nstmt.Get(&bosima, map[string]interface{}{"n": "BOSIMA"})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("NamedStmt-Get2:", bosima.City)

	stmt, err := db.Preparex("SELECT * FROM Person WHERE Name=?")
	if err != nil {
		log.Println(err)
		return
	}
	err = stmt.Get(&bosima, "BOSIMA")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Stmt-Get1:", bosima.City)
	err = stmt.Get(&bossma, "BOSSMA")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("Stmt-Get2:", bossma.City)

	// Transaction
	tx := db.MustBegin()
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Li Si Hai", "Dong Bei", time.Now(), time.Now())
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tx-MustBegin is successful")

	tx, err = db.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Li Si Hai", "Dong Bei", time.Now(), time.Now())
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tx-Beginx is successful")

	tx, err = db.Beginx()
	if err != nil {
		log.Println(err)
		return
	}
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Li Si Hai", "Dong Bei", time.Now(), time.Now())
	err = tx.Rollback()
	if err != nil {
		log.Println(err)
		return
	}

	rowCountMap := 0
	err = db.Get(&rowCountMap, "select count(*) As RowCount from Person where Name=?", "Zhang San")
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tx-Rollback is successful?", rowCountMap == 2)

	tx, err = db.BeginTxx(context.Background(), &sql.TxOptions{Isolation: sql.LevelRepeatableRead})
	if err != nil {
		log.Println(err)
		return
	}
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Zhang San", "Beijing", time.Now(), time.Now())
	tx.MustExec("INSERT INTO Person (Name, City, AddTime, UpdateTime) VALUES (?, ?, ?, ?)", "Li Si Hai", "Dong Bei", time.Now(), time.Now())
	err = tx.Commit()
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("tx-BeginTxx is successful")
}
