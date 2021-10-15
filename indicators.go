package main

import (
	"database/sql"

	"log"
	"time"

	_ "github.com/go-sql-driver/mysql"
)

func readCalcIndicators(note string) map[string]string {
	//Brings out all calculated indicators

	db, err := sql.Open("mysql", conn)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Connection Error, readCalcIndicators: "+sql.ErrNoRows.Error()+"\n")
		log.Fatal("Connection error.")
	}
	defer db.Close()
	sqlselect := "SELECT indicatorID, formula FROM indicator where notes='" + note + "' and type='calculation';"

	res, err := db.Query(sqlselect)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Select Error, readCalcIndicators: "+sql.ErrNoRows.Error()+"\n")
		log.Fatal(err)
	}
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
	db, err := sql.Open("mysql", conn)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+" Connection Error, readDateCompany: "+sql.ErrConnDone.Error()+"\n")
		log.Fatal(err)
	}

	defer db.Close()
	sqlselect := "select distinct CONCAT(companyID,',', reportDate) as c from finval where reportDate='" + forDate + "' order by reportDate, companyID;"
	res, err := db.Query(sqlselect)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+" Scan Error, readDateCompany: "+sql.ErrConnDone.Error()+"\n")
		log.Fatal(err)
	}
	defer res.Close()
	dbvalues := make([]string, numrows)
	var c string
	i := 0

	for res.Next() {
		err := res.Scan(&c)
		if err != nil {
			logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Scan Error, readDateCompany: "+"\n")
			log.Fatal(err)
		}
		dbvalues[i] = c

		i++
	}
	return dbvalues
}

func countDateCompany() int {
	//Calculate number of records for  (companyID, reportDate) in finval after AFTERYEAR
	db, err := sql.Open("mysql", conn)

	if err != nil {
		log.Fatal(err)
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Connection Error, countDateCompany: "+sql.ErrConnDone.Error()+"\n")
	}
	defer db.Close()
	sqlselect := "SELECT  COUNT(DISTINCT reportDate)*Count(DISTINCT companyID) as num FROM finval where reportDate='" + forDate + "';"
	res, err := db.Query(sqlselect)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Select Error, countDateCompany: "+sql.ErrNoRows.Error()+"\n")
		log.Fatal(err)

	}
	defer res.Close()
	var n, count int

	for res.Next() {
		err := res.Scan(&n)
		if err != nil {
			logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Scan Error, countDateCompany: "+sql.ErrConnDone.Error()+"\n")
			log.Fatal(err)
		}
		count = n
	}
	return count

}
