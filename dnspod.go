package tools

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
)

const Token = "1111111111111111111111111"

var ErrBadStatus = errors.New("status is not 1")

type Status struct {
	Code string `json:"code"`
	// ...
}

func DNSPodRecordModify(domain, sub, rid, nip string) error {
	type RecordModifyStruct struct {
		Status `json:"status"`
	}
	v := url.Values{
		"domain":         {domain},
		"record_id":      {rid},
		"sub_domain":     {sub},
		"record_type":    {"A"},
		"record_line":    {"默认"},
		"record_line_id": {"0"},
		"value":          {nip},
		"ttl":            {"600"},
		"format":         {"json"},
		"login_token":    {Token},
	}
	u := "https://dnsapi.cn/Record.Modify"
	var dpr RecordModifyStruct
	err := SendPost(u, &v, &dpr)
	if err != nil {
		return err
	}
	if dpr.Status.Code != "1" {
		fmt.Println("status code is", dpr)
		return ErrBadStatus
	}
	return nil
}

func DNSPodRecordList(domain, sub string) (string, error) {
	type RecordStruct struct {
		Id string `json:"id"`
	}
	type RecordListStruct struct {
		Status  `json:"status"`
		Records []RecordStruct `json:"records"`
	}
	v := url.Values{
		"domain":      {domain},
		"sub_domain":  {sub},
		"login_token": {Token},
		"format":      {"json"},
	}
	u := "https://dnsapi.cn/Record.List"
	var dpr RecordListStruct
	err := SendPost(u, &v, &dpr)
	if err != nil {
		return "", err
	}
	if dpr.Status.Code != "1" {
		fmt.Println(dpr.Status.Code)
		return "", ErrBadStatus
	}
	return dpr.Records[0].Id, nil
}

func SendPost(u string, v *url.Values, ret interface{}) error {
	resp, err := http.PostForm(u, *v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, &ret)
}

func ModifyRecord(sub, domain, nip string) error {
	rid, err := DNSPodRecordList(domain, sub)
	if err != nil {
		return err
	}
	return DNSPodRecordModify(domain, sub, rid, nip)
}

// Call ModifyRecord(sub, domain, nip)
