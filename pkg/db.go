package pkg

import (
	"database/sql"
	"fmt"
	"os"
	"sync"

	_ "github.com/lib/pq"
)

var (
	DBEngine *sql.DB
	DBOnce   sync.Once
	DBmu     sync.Mutex
)

var CryptoDBsetting = &DBSetting{
	Host:     "db", // should be changed to run unit test, eg: 192.168.56.127:5432
	DBName:   "crypto",
	UserName: "postgres",
	Password: "postgres",
}

type DBSetting struct {
	UserName string
	Password string
	Host     string
	DBName   string
}

func NewDBEngine(s *DBSetting) (*sql.DB, error) {
	dns := fmt.Sprintf("postgres://%s:%s@%s/%s?sslmode=disable", s.UserName, s.Password, s.Host, s.DBName)
	db, err := sql.Open("postgres", dns)
	if err != nil {
		return nil, err
	}
	return db, nil
}

func SetupDBEngine() error {
	var err error
	DBOnce.Do(func() {
		DBmu.Lock()
		defer DBmu.Unlock()

		if h := os.Getenv("pq_host"); h != "" { // should be changed to run unit test, eg: 192.168.56.127:5432
			CryptoDBsetting.Host = h
		}
		DBEngine, err = NewDBEngine(CryptoDBsetting)
	})
	if err != nil {
		return err
	}
	return nil
}

func SetupTestDBEngine() error {
	var err error
	if h := os.Getenv("pq_host"); h != "" { // should be changed to run unit test, eg: 192.168.56.127:5432
		CryptoDBsetting.Host = h
	}
	DBEngine, err = NewDBEngine(CryptoDBsetting)
	if err != nil {
		return err
	}
	return nil
}
