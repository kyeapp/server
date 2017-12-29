package autocomplete

import (
	"errors"
	"fmt"
	"bufio"
	"os"
	"time"
)

type Trie struct {
	base     string
	isWord   bool
	TrieChar map[byte]*Trie
}

func (t *Trie) Init() {
	t.isWord = false
	t.TrieChar = make(map[byte]*Trie)
}

// adds word into the Trie
func (t *Trie) add(word []byte) {

	//if end of word add new Trie and mark
	if len(word) == 0 {
		t.isWord = true
		return
	}

	letter := word[0]
	var next *Trie
	var ok bool
	if next, ok = t.TrieChar[letter]; !ok {
		newT := new(Trie)
		newT.Init()
		newT.base = t.base + string(letter)
		t.TrieChar[letter] = newT
		next = t.TrieChar[letter]
	}

	next.add(word[1:])
}

// find the top level Trie of possible matchine autocomplete words
func (t *Trie) findRoot(w []byte) (*Trie, error) {
	if len(w) == 0 {
		return t, nil
	}
	nextLetter := w[0]
	nextTrie := t.TrieChar[nextLetter]
	if nextTrie == nil {
		return nil, errors.New("no word possible matches found")
	}
	return nextTrie.findRoot(w[1:])
}



//lists words from the current Trie
func (t *Trie) listWords() (list []string) {
	if t.isWord {
		list = append(list, t.base)
	}

	for _, childTrie := range t.TrieChar {
		list = append(list, childTrie.listWords()...)
	}

	return list
}

//load dictionary into the Trie
func LoadDictionary(root *Trie, filename string) {
	defer timeTrack(time.Now(), "load Dictionary")
	file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		root.add([]byte(scanner.Text()))
	}

	if err := scanner.Err(); err != nil {
		panic(err)
	}
}

// returns a list of all words that are a possible autocomplete match
func (t *Trie) Autocomplete(baseWord string) []string {
	defer timeTrack(time.Now(), "autocomplete")
	start, err := t.findRoot([]byte(baseWord))
	if err != nil {
		return []string{}
	}

	return start.listWords()
}

func timeTrack(start time.Time, name string) {
    elapsed := time.Since(start)
    fmt.Printf("%s took %s\n", name, elapsed)
}

/*
func main() {
	root := new(Trie)
	root.Init()
	loadDictionary(root, "words.txt")

	_ = root.autocomplete("brin")
	//fmt.Println(list)


	fmt.Println()

}
*/
