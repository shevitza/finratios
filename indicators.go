package main

import (
	"database/sql"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

func readCalcIndicators(note string) map[string]string {
	//Brings out all calculated indicators
	db, err := sql.Open("mysql", CONN)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	sqlselect := "SELECT indicatorID, formula FROM indicator where notes='" + note + "' and type='calculation';"

	res, err := db.Query(sqlselect)
	defer res.Close()

	var key, val string
	dbvalues := make(map[string]string)
	for res.Next() {
		err := res.Scan(&key, &val)
		if err != nil {
			log.Fatal(err)
		}
		dbvalues[key] = val
	}

	return dbvalues

}

func readDateCompany(numrows int) []string {
	//Count  (companyID, reportDate) in finval after AFTERYEAR
	db, err := sql.Open("mysql", CONN)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	sqlselect := "select distinct CONCAT(companyID,',', reportDate) as c from finval where reportDate='" + forDate + "' order by reportDate, companyID;"
	res, err := db.Query(sqlselect)
	defer res.Close()
	dbvalues := make([]string, numrows)
	var c string
	i := 0

	for res.Next() {
		err := res.Scan(&c)
		if err != nil {
			log.Fatal(err)
		}
		dbvalues[i] = c

		i++
	}
	return dbvalues
}

func countDateCompany() int {
	//Calculate number of records for  (companyID, reportDate) in finval after AFTERYEAR
	db, err := sql.Open("mysql", CONN)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}

	sqlselect := "SELECT  COUNT(DISTINCT reportDate)*Count(DISTINCT companyID) as num FROM finval where reportDate='" +forDate + "';"
	res, err := db.Query(sqlselect)
	defer res.Close()
	var n, count int

	for res.Next() {
		err := res.Scan(&n)
		if err != nil {
			log.Fatal(err)
		}
		count = n
	}
	return count

}
