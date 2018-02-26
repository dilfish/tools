package tools

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"time"
)

// curl --get --include 'https://wordsapiv1.p.mashape.com/words/bump/also' \
//  -H 'X-Mashape-Key: RTKPpVTEZGmshPQXm2KU5BNQjAI8p1O3uCgjsnmBFWXUtuYjOE' \
//    -H 'Accept: application/json'
// test

const base = "https://wordsapiv1.p.mashape.com"

func doKGet(res, k string) ([]byte, error) {
	rq, err := http.NewRequest("GET", base+res, nil)
	if err != nil {
		return nil, err
	}
	rq.Header.Set("Accept", "application/json")
	rq.Header.Set("X-Mashape-Key", k)
	resp, err := http.DefaultClient.Do(rq)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

type DWRSyllable struct {
	Count int      `json:"count"`
	List  []string `json:"list"`
}

type DWRPronunce struct {
	All string `json:"all"`
}

type DWResult struct {
	Def    string   `json:"definition"`
	PoS    string   `json:"partOfSpeech"`
	Syno   []string `json:"synonyms"`
	To     []string `json:"typeOf"`
	Derive []string `json:"derivation"`
}

type DWRResults struct {
	Result []DWResult
}

type DetailWordResults struct {
	Results  []DWResult  `json:"results"`
	Syllable DWRSyllable `json:"syllables"`
	Pronunce DWRPronunce `json:"pronunciation"`
}

func DetailWord(word, k string) (string, error) {
	rsc := "/words/" + word
	bt, err := doKGet(rsc, k)
	if err != nil {
		return "", err
	}
	return string(bt), nil
}

var Keys = []string{ /* copied from github */ }

func WriteInfo(i, j int, name, ct string) error {
	fn := strconv.Itoa(i) + "/" + strconv.Itoa(j) + "/" + name + ".txt"
	file, err := os.Create(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = file.Write([]byte(ct))
	return err
}

func Loop() error {
	file, err := os.Open("list.txt")
	if err != nil {
		return err
	}
	defer file.Close()
	i := 1
	j := 1
	k := 0
	rd := bufio.NewReader(file)
	for {
		time.Sleep(time.Second)
		line, err := rd.ReadString('\n')
		if err != nil {
			return err
		}
		// drop 0a0d
		if line != "" {
			line = line[:len(line)-2]
		}
		if line == "" {
			continue
		}
		fmt.Println("word is", line)
		info, err := DetailWord(line, Keys[k])
		if err != nil {
			return err
		}
		k = k + 1
		if k == len(Keys) {
			k = 0
		}
		err = WriteInfo(i, j, line, info)
		if err != nil {
			return err
		}
		fmt.Println("proceed:", line, i, j)
		j = j + 1
		if j > 25 {
			j = 1
			i = i + 1
			if i > 20 {
				i = 1
			}
		}
	}
}
