// Copyright 2018 Sean.ZH

package clients

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"sort"
	"encoding/base64"
	"strings"
	"crypto/hmac"
	"crypto/sha1"
	"time"
	"math/rand"
	"bytes"
)

// AliyunClient defines a client
type AliyunClient struct {
	KeyId string
	KeySecret string
}

// NewAliyunClient returns a new client
func NewAliyunClient(id, secret string) *AliyunClient {
	return &AliyunClient{KeyId:id, KeySecret: secret}
}

// AliKeyValuePair to sort arguments
type AliKeyValuePair struct {
	Key   string
	Value string
}


// ByKey implements sort.Sort
type ByKey []AliKeyValuePair

func (a ByKey) Len() int           { return len(a) }
func (a ByKey) Swap(i, j int)      { a[i], a[j] = a[j], a[i] }
func (a ByKey) Less(i, j int) bool { return a[i].Key < a[j].Key }


func generateSignature(method string, params map[string]string, accessKeySecret string) (signature string) {
	pairs := make([]AliKeyValuePair, 0)
	for k, v := range params {
		pairs = append(pairs, AliKeyValuePair{
			Key:   k,
			Value: v,
		})
	}
	sort.Sort(ByKey(pairs))
	urlParams := ""
	for _, item := range pairs {
		if len(urlParams) > 0 {
			urlParams += "&"
		}
		urlParams += item.Key + "=" + strings.Replace(url.QueryEscape(item.Value), "+", "%20", -1)
	}
	encodedUrlParams := url.QueryEscape(urlParams)
	StringToSign := method + "&" + url.QueryEscape("/") + "&" + encodedUrlParams
	// fmt.Println("string to sign is", StringToSign)
	hmacObj := hmac.New(sha1.New, []byte(accessKeySecret + "&"))
	hmacObj.Write([]byte(StringToSign))
	signature = base64.StdEncoding.EncodeToString(hmacObj.Sum(nil))
	// fmt.Println("signature is", signature)
	return
}

func generateUrlParam(method string, params map[string]string, accessKeySecret string) (ret string) {
	for k, v := range params {
		if len(ret) > 0 {
			ret += "&"
		}
		ret += k + "=" + url.QueryEscape(v)
	}
	ret += "&Signature=" + url.QueryEscape(generateSignature(method, params, accessKeySecret))
	return
}


type AliResponse struct {
	RequestId string `json:"RequestId"` // always return
	RecordId string `json:"RecordId"` // returned when ok
	HostId string `json:"HostId"` // returned when error
	Code string `json:"Code"` // returned when error
	Message string `json:"Message"` // returned when error
}

func (ali *AliyunClient) ModifyRecord(subDomain, domain, value, recordId string) error {
	ret := make(map[string]string)
	ret["Format"] = "JSON"
	ret["Version"] = "2015-01-09"
	ret["AccessKeyId"] = ali.KeyId
	ret["SignatureMethod"] = "HMAC-SHA1"
	ret["Timestamp"] = time.Now().UTC().Format("2006-01-02T15:04:05Z")
	ret["SignatureVersion"] = "1.0"
	rand.Seed(time.Now().UnixNano())
	ret["SignatureNonce"] = fmt.Sprintf("%v", rand.Int())
	ret["Action"] = "UpdateDomainRecord"
	ret["DomainName"] = domain
	ret["RR"] = subDomain
	ret["Type"] = "A"
	ret["Value"] = value
	ret["RecordId"] = recordId

	body := []byte(generateUrlParam(http.MethodPost, ret, ali.KeySecret))
	ct := "application/x-www-form-urlencoded"
	buf:= bytes.NewBuffer(body)
	url := "https://alidns.aliyuncs.com"
	resp, err := http.Post(url, ct, buf)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	bt, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	var aliResp AliResponse
	err = json.Unmarshal(bt, &aliResp)
	if err != nil {
		return err
	}
	if aliResp.RecordId == "" {
		return errors.New(string(bt))
	}
	return nil
}
