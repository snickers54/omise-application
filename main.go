package main

import (
	"github.com/snickers54/omise-application/ROTreader"
)

func main() {
	ROTreader.NewRot128Reader("./data/fng.1000.csv.rot128")
	for s.Scan() {
		s.Text()
	}
}
