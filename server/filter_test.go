package main

import (
	"fmt"
	"testing"
)

func TestHandleWord(t *testing.T) {
	sensitiveList := LoadSensitiveWords()
	input := "hellboy wankycd dsviagra"

	util := NewDFAUtil(sensitiveList)
	newInput := util.HandleWord(input, '*')
	expected := "****boy *****cd ds******"
	if newInput != expected {
		t.Errorf("Expected %s, but got %s", expected, newInput)
	} else {
		fmt.Println("newInput", newInput)
	}
}
