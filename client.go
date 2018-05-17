package main

import (
	"os"
	"fmt"
)

func in_array(str string, strArr []string) bool {
	for _,val := range strArr {
		if str == val {
			return true
		}
	}
	return false
}

func main()  {
	argswithPro := os.Args
	commands := [] string{"status","start","stop","restart"}
	argswithoutPro := argswithPro[1:]
	if !in_array(argswithoutPro[0], commands) {

	}
}