package main

import (
	"fmt"
	"testing"
)

func TestIsGm(t *testing.T) {
	command1 := "/help"
	command2 := "help"
	if !IsGm(command1) {
		t.Fatal("error1")
	}
	if IsGm(command2) {
		t.Fatal("error2")
	}
}

func TestLoadSensitiveWords(t *testing.T) {
	strList := LoadSensitiveWords()
	if len(strList) != 451 {
		t.Fatal("load error")
	}
}

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
