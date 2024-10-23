package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-sql-driver/mysql"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
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
