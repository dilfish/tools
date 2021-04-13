package net

import (
	"crypto/sha512"
	"fmt"
	"log"
	"net/http"
	"time"
)

type CookieManager struct {
	CookiePass string
}

func NewCookieManager(pass string) *CookieManager {
	var c CookieManager
	c.CookiePass = pass
	return &c
}
func (cm *CookieManager) CheckCookie(r *http.Request, key, cookieKey string) string {
	c, err := r.Cookie(key)
	if err != nil {
		if err != http.ErrNoCookie {
			log.Println("get cookie error:", key, err)
		}
		return ""
	}
	should := cm.calcCookie(cm.CookiePass, cookieKey)
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
	c.Expires = time.Now().Add(time.Hour * 24)
	c.Value = cm.calcCookie(cm.CookiePass, cookieKey)
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

func (cm *CookieManager) calcCookie(val, cookieKey string) string {
	bt := sha512.Sum512([]byte(cookieKey + val))
	return fmt.Sprintf("%x", bt)
}
