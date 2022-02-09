package net

import (
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type CookieManager struct {
	CookiePass      string
	ExpireIntervalS int64
}

func NewCookieManager(pass, prefix string, exp int64) *CookieManager {
	var c CookieManager
	c.CookiePass = pass
	if exp < 60 {
		exp = 60
	}
	c.ExpireIntervalS = exp
	return &c
}

func SplitTimestamp(cookie string) int64 {
	arr := strings.Split(cookie, "-")
	if len(arr) != 2 {
		return -1
	}
	num, err := strconv.ParseInt(arr[0], 10, 64)
	if err != nil {
		return -2
	}
	return num
}

func (cm *CookieManager) CheckCookie(r *http.Request, key, cookieKey string) string {
	c, err := r.Cookie(key)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println("get cookie error:", key, err)
		}
		return ""
	}
	ts := SplitTimestamp(c.Value)
	if ts <= 0 || ts < time.Now().Unix() {
		return ""
	}
	should := cm.calcCookie(cm.CookiePass, cookieKey, ts)
	if should != c.Value {
		log.Println("cookie should be", should, "but we got", c.Value)
		return ""
	}
	return c.Value
}

func (cm *CookieManager) SetCookie(w http.ResponseWriter, key, cookieKey string) {
	var c http.Cookie
	c.Name = key
	c.Path = "/"
	c.Expires = time.Now().Add(time.Second * time.Duration(cm.ExpireIntervalS))
	c.Value = cm.calcCookie(cm.CookiePass, cookieKey, time.Now().Unix() + cm.ExpireIntervalS)
	http.SetCookie(w, &c)
}

func (cm *CookieManager) ClrCookie(w http.ResponseWriter, key string) {
	var c http.Cookie
	c.Name = key
	c.Value = ""
	c.Path = "/"
	c.Expires = time.Now()
	http.SetCookie(w, &c)
}

func (cm *CookieManager) calcCookie(val, cookieKey string, unix int64) string {
	ts := strconv.FormatInt(unix, 10)
	bt := sha512.Sum512([]byte(cookieKey + val + ts))
	return ts + "-" + fmt.Sprintf("%x", bt)
}
