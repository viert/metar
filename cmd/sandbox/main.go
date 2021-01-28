package main

import (
	"fmt"
	"time"

	"github.com/viert/metar"
)

func main() {
	m := metar.New(300 * time.Second)
	time.Sleep(3 * time.Second)
	fmt.Println(m.GetAirportData("EDDF"))
}
