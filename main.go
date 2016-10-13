package main

import (
	"fmt"

	pmc "./pmcounter"
)

func main() {
	PMCounters := pmc.GetPMCountersForDevice("1234")
	fmt.Println(PMCounters[0])
}
