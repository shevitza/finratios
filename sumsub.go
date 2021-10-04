package main

import (
	"database/sql"
	"fmt"
	"log"
	s "strings"
)

func calc_sumsub(indicators map[string]string, datecompany []string) {
	for j := 0; j < len(datecompany); j++ {
		dateAndCompany := datecompany[j]
		sqlselect_sumsub(indicators, dateAndCompany)

	}
}
func sqlselect_sumsub(indicators map[string]string, dateAndCompany string) {
	dateAndCompany_arr := s.Split(dateAndCompany, ",")
	reportDate := dateAndCompany_arr[1]
	companyID := dateAndCompany_arr[0]
	dateAndCompany = s.Replace(dateAndCompany, ",", " AND reportDate = '", -1)
	dateAndCompany = " AND companyID=" + dateAndCompany + "' "

	for k, v := range indicators {
		arr := formulaRead(v)

		if arr[0] == "sumsub:" {
			//	fmt.Println("Calc for indicatorID: ", k)
			indicatorID := k
			sqlSumSub := ""
			for i := 1; i < len(arr)-1; i++ {
				if arr[i][0:1] == "-" {

					sqlSumSub = sqlSumSub + " SELECT -finValue FROM finval  WHERE IndicatorID=" + s.Replace(arr[i], "-", "", -1) + dateAndCompany + " UNION ALL "
				} else {
					sqlSumSub = sqlSumSub + " SELECT finValue FROM finval WHERE IndicatorID=" + s.Replace(arr[i], "-", "", -1) + dateAndCompany + " UNION ALL "
				}
			}
			i := len(arr) - 1

			if arr[i][0:1] == "-" {

				sqlSumSub = sqlSumSub + " SELECT -finValue FROM finval  WHERE IndicatorID=" + s.Replace(arr[i], "-", "", -1) + dateAndCompany + "; "

			} else {
				sqlSumSub = sqlSumSub + " SELECT finValue FROM finval WHERE IndicatorID=" + s.Replace(arr[i], "-", "", -1) + dateAndCompany + "; "

			}
			//	fmt.Println(sqlSumSub)
			dbvalues := sqlcalc_sumsub(sqlSumSub)
			replace_sumsub_value(companyID, reportDate, indicatorID, dbvalues)

		}

	}
}
func sqlcalc_sumsub(sqlSumSub string) float64 {
	db, err := sql.Open("mysql", CONN)
	res, err := db.Query(sqlSumSub)
	defer db.Close()
	if err != nil {
		log.Fatal(err)
	}
	//Get indicator values for ratio calculation
	var dbvalues float64

	var val float64
	for res.Next() {
		err := res.Scan(&val)
		if err != nil {
			log.Fatal(err)
		}
		dbvalues = dbvalues + float64(val)
	}

	return dbvalues
}

func replace_sumsub_value(companyID, reportDate, indicatorID string, dbvalues float64) {
	dbvalues_string := fmt.Sprintf("%f", dbvalues)
	sqlreplace := "REPLACE into finval values(" + string(companyID) + ", '" + reportDate + "', " + string(indicatorID) + ", " + dbvalues_string + ", '',  current_timestamp());"
	db, err := sql.Open("mysql", CONN)
	defer db.Close()
	_, err = db.Exec(sqlreplace)
	if err != nil {
		log.Fatal(err)
	}
}
