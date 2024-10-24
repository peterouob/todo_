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
	dsn, err := url.Parse(os.Getenv("DSN"))
	if err != nil {
		log.Printf("error in parse url:%s", err.Error())
	}
	dsn.RawQuery = os.Getenv("MODE")
	db, err := gorm.Open(mysql2.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		log.Printf("error in connect mysql:%s", err.Error())
	}
	DB = db
	log.Println("connect mysql ...")
	if err := DB.AutoMigrate(&model.User{}); err != nil {
		log.Panicf("error to migrate model.User:%s", err.Error())
	}
}

func initMysql() {
	rootCertPool := x509.NewCertPool()
	pem, _ := os.ReadFile(os.Getenv("MYSQL_PEM_KEY"))
	rootCertPool.AppendCertsFromPEM(pem)
	mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	})
}
