package net

import (
	"net/http"
	"strings"
)

const BlockHTML = `
<!DOCTYPE html>
<html>
<head>
<title>Wechat browser is blocked by ooxx team!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Wechat browser is blocked by ooxx team !</h1>
<h1>不支持微信内部访问，请使用 Chrome/Firefox/Safari/Edge/cURL/NetCat</h1>
<h1>强者可以不了解弱者，但是弱者必须了解强者才能保护自己 -- <我们无处安放的青春></h1>
<p><em>Thank you for visiting OOXX.DEV.</em></p>
<p><em>Powered by ooxx proactive blocking system®©</em></p>
</body>
</html>
`

func CheckBlocked(r *http.Request) bool {
	tc := r.Header["X-Requested-With"]
	if len(tc) == 1 && tc[0] == "com.tencent.mm" {
		return true
	}
	tc = r.Header["User-Agent"]
	if len(tc) == 1 && strings.Index(tc[0], "MicroMessenger") >= 0 {
		return true
	}
	return false
}
