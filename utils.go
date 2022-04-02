package main

import (
	"fmt"
	"time"

	"github.com/fatih/color"
)

func diplayTimestamp() {

	color.Set(color.BgHiYellow)
	color.Set(color.FgBlack)
	defer color.Unset()

	t := time.Now()
	fmt.Println(t.Format("2006-01-02 15:04:05"))
}
