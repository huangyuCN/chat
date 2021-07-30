package main

import (
	"fmt"
	"strings"
	"testing"
	"time"
)

func TestPopular_Count(t *testing.T) {
	popular := NewPopular()
	words := " hello hello   hello\n you and me hello you"
	strList := strings.Split(words, " ")
	fmt.Println("strList:", strList)
	popular.text <- words
	popular.text <- words
	popular.text <- words
	time.Sleep(100000000)
	str := popular.Count(60)
	fmt.Println("str:", str)
	if str != "hello" {
		t.Fatal("error")
	}

}
