package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-sql-driver/mysql"
	"github.com/peterouob/todo_/model"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
)

var DB *gorm.DB

func InitMysql() {
	initMysql()
	dsn, _ := url.Parse(os.Getenv("DSN"))
	dsn.RawQuery = os.Getenv("MODE")
	db, _ := gorm.Open(mysql2.Open(dsn.String()), &gorm.Config{})
	DB = db
	log.Println("connect mysql ...")
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Panicf("error to migrate model.User:%s", err.Error())
	}
}

func initMysql() {
	rootCertPool := x509.NewCertPool()
	pem, _ := os.ReadFile(os.Getenv("PEM_PATH"))
	rootCertPool.AppendCertsFromPEM(pem)
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	})
}
