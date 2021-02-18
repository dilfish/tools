// Copyright 2018 Sean.ZH

package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
)

// DnspodClient defines a client
type DnspodClient struct {
	Token string
}

// NewDnspodClient returns a new client
func NewDnspodClient(token string) *DnspodClient {
	return &DnspodClient{Token: token}
}

// ErrBadStatus is a error status for dnspod
var ErrBadStatus = errors.New("status is not 1")

// Status we just need to known it's status
type Status struct {
	Code string `json:"code"`
	// ...
}

// DNSPodRecordModify modify a record
func dNSPodRecordModify(token, domain, sub, rid, nip string) error {
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
		"login_token":    {token},
	}
	u := "https://dnsapi.cn/Record.Modify"
	var dpr RecordModifyStruct
	err := sendPost(u, &v, &dpr)
	if err != nil {
		return err
	}
	if dpr.Status.Code != "1" {
		fmt.Println("status code is", dpr)
		return ErrBadStatus
	}
	return nil
}

// DNSPodRecordList read all records
func dNSPodRecordList(token, domain, sub string) (string, error) {
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
		"login_token": {token},
		"format":      {"json"},
	}
	u := "https://dnsapi.cn/Record.List"
	var dpr RecordListStruct
	err := sendPost(u, &v, &dpr)
	if err != nil {
		return "", err
	}
	if dpr.Status.Code != "1" {
		fmt.Println(dpr.Status.Code)
		return "", ErrBadStatus
	}
	return dpr.Records[0].Id, nil
}

// SendPost send post to api
func sendPost(u string, v *url.Values, ret interface{}) error {
	resp, err := http.PostForm(u, *v)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	return json.Unmarshal(bt, &ret)
}

// ModifyRecord is a demo
func (d *DnspodClient) ModifyRecord(sub, domain, nip string) error {
	rid, err := dNSPodRecordList(d.Token, domain, sub)
	if err != nil {
		return err
	}
	return dNSPodRecordModify(d.Token, domain, sub, rid, nip)
}

// Call ModifyRecord(sub, domain, nip)
