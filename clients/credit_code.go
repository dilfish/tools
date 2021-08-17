// Copyright 2018 Sean.ZH

package clients

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"os"
	"time"

	dio "github.com/dilfish/tools/io"
)

var resultFile *os.File
var counter int

// ResultData is response json
type ResultData struct {
	Ent  string `json:"entname"`
	Code string `json:"creditCode"`
	// ...
}

// Result is response json
type Result struct {
	Msg    string       `json:"message"`
	Status int          `json:"status"`
	Result []ResultData `json:"results"`
	// ...
}

// Go is to go download
func Go(name string) error {
	time.Sleep(time.Millisecond * 500)

	uname := url.QueryEscape(name)
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr}

	uri := "https://www.creditchina.gov.cn/api/public_search/getCreditCodeFacades?keyword="
	end := "&filterManageDept=0&filterOrgan=0&filterDivisionCode=0&page=1&pageSize=10&_=1534765935044"
	uri = uri + uname + end
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return err
	}
	ua := "Mozilla/5.0 (Macintosh; Intel Mac OS X 10_11_6) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/68.0.3440.106 Safari/537.36"
	req.Header.Add("Accept", "application/json, text/javascript, */*; q=0.01")
	req.Header.Add("Accept-Encoding", "gzip, deflate, br")
	req.Header.Add("Accept-Language", "en,zh-CN;q=0.9,zh;q=0.8,zh-TW;q=0.7,ja;q=0.6")
	req.Header.Add("Cache-Control", "no-cache")
	req.Header.Add("Connection", "keep-alive")
	req.Header.Add("Host", "www.creditchina.gov.cn")
	req.Header.Add("Pragma", "no-cache")
	req.Header.Add("Sec-Metadata", "cause=forced, destination=\"\", target=subresource, site=same-origin")
	req.Header.Add("User-Agent", ua)
	req.Header.Add("X-Requested-With", "XMLHttpRequest")

	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var ret Result
	err = json.Unmarshal(bt, &ret)
	if err != nil {
		return err
	}
	if len(ret.Result) != 1 {
		return errors.New(ret.Msg)
	}
	io.WriteString(resultFile, name+"\t"+ret.Result[0].Code+"\n")
	return nil
}

func readFile(line string) error {
	counter = counter + 1
	if counter%10 == 0 {
		fmt.Println("Finishied:", counter)
	}
	err := Go(line)
	if err != nil {
		io.WriteString(resultFile, line+"\t"+"Error:"+err.Error()+"\n")
	}
	return nil
}

// ReadInput read files
func ReadInput() error {
	return dio.ReadLine("./list.txt", readFile)
}

// GetCreditCode get code in a file
func GetCreditCode() {
	var err error
	resultFile, err = os.Create("./result.txt")
	if err != nil {
		panic("could not create result file")
	}
	defer resultFile.Close()
	err = ReadInput()
	if err != nil {
		panic(err)
	}
}
