package tools

import (
	"context"
	"database/sql"
	"errors"
	"github.com/go-sql-driver/mysql"
	"strconv"
	"strings"
	"time"
)

type DBConfig struct {
	User   string `json:"user"`
	Pass   string `json:"pass"`
	Host   string `json:"host"`
	Port   int    `json:"port"`
	DBName string `json:"db"`
}

func initDB(conf *DBConfig) (*sql.DB, error) {
	dsn := conf.User + ":" + conf.Pass + "@tcp"
	dsn = dsn + "(" + conf.Host + ":"
	dsn = dsn + strconv.Itoa(conf.Port) + ")"
	dsn = dsn + "/" + conf.DBName
	return sql.Open("mysql", dsn)
}

func InitDB(conf *DBConfig) (*sql.DB, error) {
	db, err := initDB(conf)
	if err != nil {
		return nil, err
	}
	return db, db.Ping()
}

type MultiDB struct {
	conf []DBConfig
	db   *sql.DB
	idx  int
	num  int
}

var cancel context.CancelFunc
var ctx context.Context

var ErrEmptyConf = errors.New("empty config")
var ErrAllDead = errors.New("all connections are dead")

func InitMultiDB(conf []DBConfig) (*MultiDB, error) {
	if len(conf) == 0 {
		return nil, ErrEmptyConf
	}
	var m MultiDB
	m.num = len(conf)
	m.conf = conf
	for idx, c := range conf {
		db, err := initDB(&c)
		if err != nil {
			return nil, err
		}
		err = db.Ping()
		if err == nil {
			m.db = db
			m.idx = idx
			ctx, cancel = context.WithCancel(context.Background())
			return &m, nil
		}
		db.Close()
	}
	return nil, ErrAllDead
}

func (m *MultiDB) Close() error {
	if m.db != nil {
		return m.db.Close()
	}
	return nil
}

func (m *MultiDB) switchDB() {
	// avoid switch infinitely
	time.Sleep(time.Second * 5)
	if cancel != nil {
		cancel()
	}
	var err error
	ctx, cancel = context.WithCancel(context.Background())
	m.idx = m.idx + 1
	if m.idx == m.num-1 {
		m.idx = 0
	}
	// we should check conf format at init
	m.db, err = initDB(&m.conf[m.idx])
	if err != nil {
		m.switchDB()
	}
}

func checkConnErr(err error) bool {
	switch err {
	case mysql.ErrMalformPkt:
		return true
	case mysql.ErrInvalidConn:
		return true
		// timeout
	}
	if strings.Index(err.Error(), "connection refused") >= 0 {
		return true
	}
	return false
}

func (m *MultiDB) Query(q string) (*sql.Rows, error) {
	r, err := m.db.QueryContext(ctx, q)
	if checkConnErr(err) == true {
		m.switchDB()
		m.Query(q)
	}
	return r, err
}

// func (m *MultiDB) Exec (q string) (*sql.Rows, error) {} ...
