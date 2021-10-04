package main

const CONN = "tera:111111@tcp(127.0.0.1:3306)/sqlt"
const AFTERYEAR = "2015"

type Company struct {
	companyID  string
	reportDate string
	ind1       string
	ind2       string
	indres     string
}

var forDate = "2019-06-30"

//TODO connection string in separated file
//TODO writing error logs in separated text file
//TODO set argument and its control
