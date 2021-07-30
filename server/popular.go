package main

import (
	"fmt"
	"sort"
	"strings"
	"time"
)

type popular struct {
	words map[int64]map[string]int64
	times []int64
	text  chan string
}

func NewPopular() *popular {
	popular := &popular{
		words: make(map[int64]map[string]int64),
		text:  make(chan string, 10),
		times: make([]int64, 60),
	}
	go popular.listen()
	return popular
}

func (p *popular) listen() {
	for {
		select {
		case text, ok := <-p.text:
			if !ok {
				return
			}
			now := time.Now().Unix()
			m, find := p.words[now]
			if !find {
				p.times = append(p.times, now)
				m = make(map[string]int64)
			}
			strList := strings.Split(text, " ")
			for _, s := range strList {
				s = strings.Trim(s, " ")
				s = strings.Trim(s, "\n")
				num, numFind := m[s]
				if numFind {
					m[s] = num + 1
				} else {
					m[s] = 1
				}
			}
			p.words[now] = m
			if len(p.times) > 60 {
				delete(p.words, p.times[0])
				p.times = p.times[1:]
			}
		}
	}
}

func (p *popular) Count(seconds int64) string {
	now := time.Now().Unix()
	min := now - seconds
	m := make(map[string]int64)
	for k, v := range p.words {
		if k > min {
			for k1, v1 := range v {
				oldNum, ok := m[k1]
				if ok {
					m[k1] = oldNum + v1
				} else {
					m[k1] = v1
				}
			}
		}
	}
	if len(m) == 0 {
		return ""
	}
	list := rankByWordCount(m)
	fmt.Println("key:", list[0].Key, "times:", list[0].Value)
	return list[0].Key
}

func rankByWordCount(wordFrequencies map[string]int64) PairList {
	pl := make(PairList, len(wordFrequencies))
	i := 0
	for k, v := range wordFrequencies {
		pl[i] = Pair{k, v}
		i++
	}
	sort.Sort(sort.Reverse(pl))
	return pl
}

type Pair struct {
	Key   string
	Value int64
}

type PairList []Pair

func (p PairList) Len() int           { return len(p) }
func (p PairList) Less(i, j int) bool { return p[i].Value < p[j].Value }
func (p PairList) Swap(i, j int)      { p[i], p[j] = p[j], p[i] }
