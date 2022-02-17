package net

import (
	"io"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"

	dio "github.com/dilfish/tools/io"
)

type UploaderService struct {
	MaxSize  int64
	MaxMem   int64
	Curr     int64
	BasePath string
	BaseURL  string
	NameLen  int
	Expire   time.Duration
	Lock     sync.Mutex
	Map      map[string]time.Time
}

// WriteFile write reader into file
func (u *UploaderService) WriteFile(name string, rc io.Reader) (int64, string, error) {
	ext := filepath.Ext(name)
	if ext == "" {
		ext = ".noext"
	}
	name = dio.RandStr(u.NameLen) + ext
	fn := u.BasePath + "/" + name
	file, err := os.Create(fn)
	if err != nil {
		log.Println("create file name error:", name, err)
		return 0, "", err
	}
	defer file.Close()
	u.Lock.Lock()
	defer u.Lock.Unlock()
	u.Map[fn] = time.Now().Add(u.Expire)
	n, err := io.Copy(file, rc)
	return n, name, err
}

// Handler return page if get
// and write file into disk if post
func (u *UploaderService) Handler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		io.WriteString(w, "Not Supported")
		return
	}
	err := r.ParseMultipartForm(u.MaxMem)
	if err != nil {
		io.WriteString(w, "Read File Error:"+err.Error())
		return
	}

	file, header, err := r.FormFile("file")
	if err != nil {
		io.WriteString(w, "Read File error:"+err.Error())
		return
	}
	if u.Curr+header.Size > u.MaxSize {
		log.Println("too many write:", u.Curr, header.Size, u.MaxSize)
		io.WriteString(w, "Too many write")
		return
	}
	defer file.Close()
	n, name, err := u.WriteFile(header.Filename, file)
	if err != nil {
		log.Println("write file error:", err)
		io.WriteString(w, "write file error"+err.Error())
		return
	}
	u.Curr = u.Curr + n
	io.WriteString(w, "<html lang=\"zh-cmn-Hans\"><head><meta charset=\"UTF-8\"></head><h1>上传成功！，你可以访问这里看一看:<a href=\""+u.BaseURL+name+"\">File</a></h1></html>")
	return
}

func NewUploadService(baseURL, basePath string, maxSize int64, expire time.Duration, nameLen int) *UploaderService {
	var u UploaderService
	u.MaxSize = maxSize
	u.MaxMem = maxSize
	u.BasePath = basePath
	u.BaseURL = baseURL
	u.NameLen = nameLen
	if expire < time.Minute {
		expire = time.Minute
	}
	u.Expire = expire
	if u.NameLen < 1 {
		u.NameLen = 10
	}
	u.Map = make(map[string]time.Time)
	go u.Patrol()
	return &u
}

func (u *UploaderService) Patrol() {
	for {
		time.Sleep(time.Minute)
		tbd := []string{}
		u.Lock.Lock()
		for k, v := range u.Map {
			if v.After(time.Now()) {
				tbd = append(tbd, k)
			}
		}
		for _, tb := range tbd {
			os.Remove(tb)
			delete(u.Map, tb)
		}
		u.Lock.Unlock()
	}
}
