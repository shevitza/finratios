package main

import (
	"database/sql"
	"fmt"
	"log"
	s "strings"
	"time"
)

func calc_adj(indicators map[string]string, datecompany []string, numrows int) {
	var c Company
	for k, v := range indicators {
		arr := formulaRead(v)
		if arr[0] == "adj:" {
			for i := 0; i < numrows; i++ {
				c.indres = k
				c.ind1 = arr[1]
				temp := s.Split(datecompany[i], ",")

				c.companyID = temp[0]
				c.reportDate = temp[1]
				c.ind2 = "-"
				//	if s.HasSuffix(temp[1], "-12-31") {Nothing}

				if s.HasSuffix(temp[1], "-03-31") {
					setTheSame(c)
				}
				//		if s.HasSuffix(temp[1], "-06-30") { Nothing}

				if s.HasSuffix(temp[1], "-09-30") {
					// calc adjusted values for Q2 and Q3 so: Q3:=Q3-Q2-Q1, Q2:=Q2-Q1
					adjCalc(c)
				}

			}

		}

	}

}

func setTheSame(c Company) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Connection  Error, setTheSame: "+"\n")
		log.Fatal(err)
	}
	defer db.Close()

	sqlstring := "REPLACE INTO finval (companyID, reportDate, indicatorID, finValue) select companyID, reportDate, 67, finvalue from finval"
	sqlstring = sqlstring + " where companyID=" + c.companyID + " and reportDate='" + c.reportDate + "' and indicatorID=" + c.ind1 + "; "
	_, err = db.Exec(sqlstring)

	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Replace  Error, setTheSame: "+"\n")
		log.Fatal(err)

	}

}

func adjCalc(c Company) {
	temp := c.reportDate
	year := temp[0:4]

	db, err := sql.Open("mysql", conn)

	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()
	sqlstring := "SELECT reportDate, finValue FROM finval WHERE companyID=" + c.companyID + " AND indicatorID=" + c.ind1 + " AND year(reportDate)=" + year + ";"

	res, err := db.Query(sqlstring)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Select Error, adjCalc: "+"\n")
		log.Fatal(err)
	}

	dbvalues := make(map[string]float64)
	var key string
	var val float64
	for res.Next() {
		err := res.Scan(&key, &val)
		if err != nil {
			logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Scan Error, adjCalc: "+"\n")
			log.Fatal(err)

		}
		dbvalues[key] = val
	}
	defer res.Close()

	dbvalues["-06-30"] = dbvalues[year+"-06-30"] - dbvalues[year+"-03-31"]
	dbvalues["-09-30"] = dbvalues[year+"-09-30"] - dbvalues[year+"-06-30"]
	dbvalues["-12-31"] = dbvalues[year+"-12-31"] - dbvalues[year+"-09-30"]

	sqlstring = "REPLACE INTO finval (companyID, reportDate, indicatorID, finValue) VALUES"

	sqlexec := sqlstring + "(" + c.companyID + ", '" + year + "-06-30" + "', " + c.indres + ", " + fmt.Sprintf("%f", dbvalues["-06-30"]) + "), "
	sqlexec = sqlexec + "(" + c.companyID + ", '" + year + "-09-30" + "', " + c.indres + ", " + fmt.Sprintf("%f", dbvalues["-09-30"]) + "), "
	sqlexec = sqlexec + "(" + c.companyID + ", '" + year + "-12-31" + "', " + c.indres + ", " + fmt.Sprintf("%f", dbvalues["-12-31"]) + "); "

	_, err = db.Exec(sqlexec)
	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Replace Error, adjCalc: "+"\n")
		log.Fatal(err)
	}

}
