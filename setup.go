package main
import "os"


const AFTERYEAR = "2015"

type Company struct {
	companyID  string
	reportDate string
	ind1       string
	ind2       string
	indres     string
}


var forDate=os.Args[1]



//TODO writing error logs in separated text file
//TODO set argument and its control

var conn string
