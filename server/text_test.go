package main

import (
	"fmt"
	"testing"
)

func TestSecondsToDayStr(t *testing.T) {
	var d, h, m, s int64 = 100, 3, 32, 18
	seconds := d*24*60*60 + h*3600 + m*60 + s
	timeStr := SecondsToDayStr(seconds)
	if timeStr != "100d 03h 32m 18s" {
		t.Fatal("error", timeStr)
	} else {
		fmt.Println(timeStr)
	}
}
