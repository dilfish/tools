package net

import (
	dio "github.com/dilfish/tools/io"
	"golang.org/x/exp/slog"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"sync"
	"time"
)

type UploaderService struct {
	MaxSize     int64
	MaxMem      int64
	Curr        int64
	BasePath    string
	BaseURL     string
	JumpBackURL string
	NameLen     int
	Expire      time.Duration
	Lock        sync.Mutex
	Map         map[string]time.Time
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
		slog.Error("create file name", err, "name", name)
		return 0, "", err
	}
	defer file.Close()
	u.Lock.Lock()
	defer u.Lock.Unlock()
	u.Map[fn] = time.Now().Add(u.Expire)
	slog.Info("upload file", "file name", fn, "path", u.Map[fn])
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
		slog.Info("too many write", "curr", u.Curr, "size", header.Size, "max size", u.MaxSize)
		io.WriteString(w, "Too many write")
		return
	}
	defer file.Close()
	n, name, err := u.WriteFile(header.Filename, file)
	if err != nil {
		slog.Error("write file", err)
		io.WriteString(w, "write file error"+err.Error())
		return
	}
	u.Curr = u.Curr + n
	show := "<html lang=\"zh-cmn-Hans\"><head><meta charset=\"UTF-8\"></head><h1>上传成功！，你可以访问这里看一看:<a href=\"" + u.BaseURL + name + "\">File</a></h1>"
	if u.JumpBackURL != "" {
		show = show + "<h1>或者你也可以再次返回<a href=\"" + u.JumpBackURL + "\">上传页面</a></h1>"
	}
	show = show + "</html>"
	io.WriteString(w, show)
	return
}

func NewUploadService(xl *slog.Logger, baseURL, basePath, jump string, maxSize int64, expire time.Duration, nameLen int) *UploaderService {
	var u UploaderService
	u.MaxSize = maxSize
	u.MaxMem = maxSize
	u.BasePath = basePath
	u.BaseURL = baseURL
	u.JumpBackURL = jump
	u.NameLen = nameLen
	if expire < time.Minute {
		expire = time.Minute
	}
	u.Expire = expire
	xl.Info("u.Expire", "expire", expire)
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
			if v.Before(time.Now()) {
				tbd = append(tbd, k)
			}
		}
		u.Lock.Unlock()
		for _, tb := range tbd {
			os.Remove(tb)
			slog.Info("uploader service remove", "file name", tb)
			delete(u.Map, tb)
		}
	}
}
