package main

import (
	"fmt"

	pmc "./pmcounter"
)

func main() {
	PMCounters := pmc.GetPMCounters()
	fmt.Println(PMCounters)
}
