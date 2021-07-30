package main

import (
	"bufio"
	"fmt"
	"io"
	"os"
	"runtime"
	"strconv"
	"strings"
)

// IsGm 判断这条消息是否为GM命令
func IsGm(text string) bool {
	return strings.HasPrefix(text, "/")
}

// Filter 过滤脏字
func Filter(text string) string {
	util := NewDFAUtil(SensitiveWords)
	text = util.HandleWord(text, '*')
	return text
}

//LoadSensitiveWords 加载所有脏字
func LoadSensitiveWords() []string {
	pwd, _ := os.Getwd()
	fileName := pwd + "/server/words"
	sysType := runtime.GOOS
	if sysType == "windows" {
		fileName = pwd + "\\server\\words"
	}
	file, err := os.OpenFile(fileName, os.O_RDWR, 0666)
	if err != nil {
		fmt.Println("Open file error!", err)
		panic(err)
	}
	defer file.Close()
	stat, err := file.Stat()
	if err != nil {
		panic(err)
	}
	var size = stat.Size()
	fmt.Println("file size=", size)
	buf := bufio.NewReader(file)
	words := make([]string, 451)
	for {
		line, err := buf.ReadString('\n')
		line = strings.TrimSpace(line)
		if err != nil {
			if err == io.EOF {
				fmt.Println("File read ok!")
				break
			} else {
				fmt.Println("Read file error!", err)
				panic(err)
			}
		}
		words = append(words, line)
	}
	return words
}

func SecondsToDayStr(seconds int64) string {
	var daySeconds int64 = 24 * 60 * 60
	var hourSeconds int64 = 60 * 60
	d := seconds / daySeconds
	r := seconds % daySeconds
	h := r / hourSeconds
	r = r % hourSeconds
	m := r / 60
	s := r % 60
	var builder strings.Builder
	dstr := strconv.FormatInt(d, 10)
	if d < 10 {
		dstr = "0" + dstr
	}
	builder.WriteString(dstr)
	builder.WriteString("d ")
	hstr := strconv.FormatInt(h, 10)
	if h < 10 {
		hstr = "0" + hstr
	}
	builder.WriteString(hstr + "h ")
	mstr := strconv.FormatInt(m, 10)
	if m < 10 {
		mstr = "0" + mstr
	}
	builder.WriteString(mstr + "m ")
	sstr := strconv.FormatInt(s, 10)
	if s < 10 {
		sstr = "0" + sstr
	}
	builder.WriteString(sstr + "s")
	return builder.String()
}
