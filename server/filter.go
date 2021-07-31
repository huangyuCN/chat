package main

type trieNode struct {
	isEndOfWord bool
	children map[rune]*trieNode
}

func newTrieNode() *trieNode {
	return &trieNode{
		isEndOfWord: false,
		children:    make(map[rune]*trieNode, 26),
	}
}

type matchIndex struct {
	start int // start index
	end   int // end index
}

func newMatchIndex(start, end int) *matchIndex {
	return &matchIndex{
		start: start,
		end:   end,
	}
}

type DFAUtil struct {
	root *trieNode
}

func (dfa *DFAUtil) insertWord(word []rune) {
	currNode := dfa.root
	for _, c := range word {
		if childNode, exist := currNode.children[c]; !exist {
			childNode = newTrieNode()
			currNode.children[c] = childNode
			currNode = childNode
		} else {
			currNode = childNode
		}
	}
	currNode.isEndOfWord = true
}

func (dfa *DFAUtil) startsWith(prefix []rune) bool {
	currNode := dfa.root
	for _, c := range prefix {
		if childNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = childNode
		}
	}
	return true
}

func (dfa *DFAUtil) searchWord(word []rune) bool {
	currNode := dfa.root
	for _, c := range word {
		if childNode, exist := currNode.children[c]; !exist {
			return false
		} else {
			currNode = childNode
		}
	}
	return currNode.isEndOfWord
}

func (dfa *DFAUtil) searchSentence(sentence string) (matchIndexList []*matchIndex) {
	start, end := 0, 1
	sentenceRuneList := []rune(sentence)
	startsWith := false
	for end <= len(sentenceRuneList) {
		if dfa.startsWith(sentenceRuneList[start:end]) {
			startsWith = true
			end += 1
		} else {
			if startsWith == true {
				for index := end - 1; index > start; index-- {
					if dfa.searchWord(sentenceRuneList[start:index]) {
						matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
						break
					}
				}
			}
			start, end = end-1, end+1
			startsWith = false
		}
	}
	if startsWith {
		for index := end - 1; index > start; index-- {
			if dfa.searchWord(sentenceRuneList[start:index]) {
				matchIndexList = append(matchIndexList, newMatchIndex(start, index-1))
				break
			}
		}
	}
	return
}

// HandleWord 查找字符串中的敏感词，并用特殊字符代替
func (dfa *DFAUtil) HandleWord(sentence string, replaceCh rune) string {
	matchIndexList := dfa.searchSentence(sentence)
	if len(matchIndexList) == 0 {
		return sentence
	}
	sentenceList := []rune(sentence)
	for _, matchIndexObj := range matchIndexList {
		for index := matchIndexObj.start; index <= matchIndexObj.end; index++ {
			sentenceList[index] = replaceCh
		}
	}
	return string(sentenceList)
}

func NewDFAUtil(wordList []string) *DFAUtil {
	dfa := &DFAUtil{
		root: newTrieNode(),
	}
	for _, word := range wordList {
		wordRuneList := []rune(word)
		if len(wordRuneList) > 0 {
			dfa.insertWord(wordRuneList)
		}
	}
	return dfa
}