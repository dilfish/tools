package net

import (
	"net/http"
	"strings"
)

const BlockHTML = `
<!DOCTYPE html>
<html>
<head>
<title>微信被DILFISH团队屏蔽!</title>
<style>
    body {
        width: 35em;
        margin: 0 auto;
        font-family: Tahoma, Verdana, Arial, sans-serif;
    }
</style>
</head>
<body>
<h1>Wechat browser is blocked by x.dilfish.icu team !</h1>
<h1>Please use desktop or mobile browsers instead</1>
<h1>不支持微信内部访问，请使用 Chrome, Firefox, Safari, Edge, Opera</h1>
<h1>或者 cURL/NetCat </h1>
<h1>强者可以不了解弱者，但是弱者必须了解强者才能保护自己 -- <我们无处安放的青春></h1>
<p><em>Thank you for visiting x.dilfish.icu.</em></p>
<p><em>Powered by x.dilfish.icu proactive blocking system®©</em></p>
</body>
</html>
`

func CheckBlocked(r *http.Request) bool {
	tc := r.Header["X-Requested-With"]
	if len(tc) == 1 && tc[0] == "com.tencent.mm" {
		return true
	}
	tc = r.Header["User-Agent"]
	if len(tc) == 1 {
        if strings.Index(tc[0], "MicroMessenger") >= 0 {
		    return true
	    }
        if tc[0] == "wx" {
            return true
        }
    }
	return false
}
