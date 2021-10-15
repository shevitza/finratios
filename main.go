package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/joho/godotenv"
)

func main() {


	start := time.Now()

	stepByStep("first_step")
	t := time.Now()
	elapsed := t.Sub(start)
	fmt.Println("Elapsed time for first step: ", elapsed)

	stepByStep("second_step")
	t = time.Now()
	elapsed = t.Sub(start)
	fmt.Println("Total elapsed time: ", elapsed)

	fmt.Println("Calculation ended!")
}

func stepByStep(step string) {
	indicators := readCalcIndicators(step)
	fmt.Println("Calculated indicators in ", step, ":")
	fmt.Println(indicators)
	numrows := countDateCompany()
	datecompany := readDateCompany(numrows)
	calcall_func1(indicators, datecompany, numrows)
	calc_adj(indicators, datecompany, numrows)
	calc_sumsub(indicators, datecompany)
}

func logerror(logfile string, errorstring string) {

	// If the file doesn't exist, create it, or append to the file
	f, err := os.OpenFile(logfile, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		log.Fatal(err)
	}
	if _, err := f.Write([]byte(errorstring)); err != nil {
		log.Fatal(err)
	}
	if err := f.Close(); err != nil {
		log.Fatal(err)
	}

}

func init() {
	fmt.Println(time.Now().Format("2006.01.02 15:04:05"))
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatal("Error loading .env file")
	}

	conn = os.Getenv("DBUSER") + ":" + os.Getenv("DBPASS") + "@tcp(" + os.Getenv("HOST") + ":" + os.Getenv("PORT") + ")/" + os.Getenv("DBNAME")

}
