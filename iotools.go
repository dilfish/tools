package tools

import (
	"io/ioutil"
	"net/http"
	"net/url"
	"os"
	"io"
	"bufio"
)

func ReadFile(fn string) ([]byte, error) {
	file, err := os.Open(fn)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	return ioutil.ReadAll(file)
}

func DoPost(url string, v *url.Values) ([]byte, error) {
	resp, err := http.PostForm(url, *v)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}

func DoGet(url string) ([]byte, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return ioutil.ReadAll(resp.Body)
}


type LineFunc func (line string) error


func ReadLine(fn string, lf LineFunc) error {
	file, err := os.Open(fn)
	if err != nil {
		return err
	}
	defer file.Close()
	rd := bufio.NewReader(file)
	for {
		line, err := rd.ReadString('\n')
		if err != nil && err != io.EOF {
			return err
		}
		if err == io.EOF {
			break
		}
		if line == "" {
			continue
		}
		line = line[:len(line) - 1]
		err = lf(line)
		if err != nil {
			return err
		}
	}
	return nil
}
