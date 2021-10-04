package main

import (
	"fmt"
	"time"
	"os"
	"log"

)
var (
    outfile, _ = os.Create("finratios.log") // log file path
    l      = log.New(outfile, "", 0)
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
