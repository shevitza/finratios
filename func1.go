package main

import (
	"database/sql"
	"fmt"
	"log"
	s "strings"

	"time"

	_ "github.com/go-sql-driver/mysql"
)

func calcall_func1(indicators map[string]string, datecompany []string, numrows int) {
	var c Company
	for k, v := range indicators {
		arr := formulaRead(v)
		if arr[0] == "func1:" {
			for i := 0; i < numrows; i++ {

				c.indres = k
				c.ind1 = arr[1]
				c.ind2 = arr[2]
				temp := s.Split(datecompany[i], ",")

				c.companyID = temp[0]
				c.reportDate = temp[1]
				//	fmt.Println(c)
				sqlcalc_func1(&c)

			}
		}

	}

}

func formulaRead(f string) []string {
	temp := s.Replace(f, "+[", "[+", -1)
	temp = s.Replace(temp, "-[", "[-", -1)
	temp = s.Replace(temp, "]", "", -1)
	temp = s.Replace(temp, "[", ",", -1)
	temp = s.Replace(temp, "/", "", -1)
	temp = s.Replace(temp, "*", "", -1)
	temp = s.Replace(temp, "+", "", -1)
	//temp = s.Replace(temp, "-", "", -1)
	temp = s.Replace(temp, "(", "", -1)
	temp = s.Replace(temp, ")", "", -1)
	temp = s.Replace(temp, "avr", "", -1)

	arr := s.Split(temp, ",")
	return arr
}

func sqlselect_func1(c *Company) string {
	str1 := " SELECT indicatorID, finvalue FROM sqlt.finval where companyID=" + c.companyID + " and reportDate='" + c.reportDate + "' and indicatorID="
	str2 := " union all "
	sqlselect := str1 + c.ind1 + str2 + str1 + c.ind2 + ";"
	return sqlselect
}

func sqlreplace_func1(c *Company, ratio float64) string {
	ratiostr := fmt.Sprintf("%f", ratio)
	x0 := "REPLACE into finval values(" + c.companyID + ", '" + c.reportDate + "', " + c.indres + ", " + ratiostr + ", '',  current_timestamp());"
	return x0

}

func ratiocalc_func1(c *Company, dbvalues map[string]float64) float64 {
	var ratio float64
	if dbvalues[c.ind2] != 0 {
		ratio = dbvalues[c.ind1] / dbvalues[c.ind2]
	} else {
		ratio = 0.
	}

	return ratio
}

func sqlcalc_func1(c *Company) {
	db, err := sql.Open("mysql", conn)

	if err != nil {
		log.Fatal(err)
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Connection Error, sqlcalc_func1: "+sql.ErrConnDone.Error()+"\n")
	}
	defer db.Close()
	sqlselect := sqlselect_func1(c)
	res, err := db.Query(sqlselect)
	if err != nil {
		log.Fatal(err)
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Connection Error, sqlcalc_func1: "+sql.ErrNoRows.Error()+"\n")

	}
	defer res.Close()

	//Get indicator values for ratio calculation
	dbvalues := make(map[string]float64)
	var key string
	var val float64
	for res.Next() {
		err := res.Scan(&key, &val)
		if err != nil {
			log.Fatal(err)
		}
		dbvalues[key] = val
	}
	//Calucate the ratio
	ratio := ratiocalc_func1(c, dbvalues)

	//Replace ratio value in the database
	sqlreplace := sqlreplace_func1(c, ratio)
	_, err = db.Exec(sqlreplace)

	if err != nil {
		logerror("finratios.log", time.Now().Format("2006.01.02 15:04:05")+"  Replace (Exec) Error, sqlcalc_func1: "+sql.ErrConnDone.Error()+"\n")
		log.Fatal(err)

	}

}
