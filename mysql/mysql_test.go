package mysql

import (
	"crypto/tls"
	"crypto/x509"
	"github.com/go-sql-driver/mysql"
	"github.com/joho/godotenv"
	mysql2 "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"log"
	"net/url"
	"os"
	"testing"
)

func init() {
	if err := godotenv.Load("test.env"); err != nil {
		log.Panicf("error to load env file ... :%s", err.Error())
	}
}

func Test_Connect(t *testing.T) {
	testInitMySQLTLSConfig(t)
	dsn, err := url.Parse(os.Getenv("DSN"))
	if err != nil {
		t.Errorf("URL parse error:%s", err.Error())
	}
	dsn.RawQuery = os.Getenv("MODE")
	if dsn.String() == "" {
		t.Errorf("not get the dsn value")
	}
	db, err := gorm.Open(mysql2.Open(dsn.String()), &gorm.Config{})
	if err != nil {
		t.Errorf("error to connect mysql :%s", err.Error())
	}
	DB = db
}

func testInitMySQLTLSConfig(t *testing.T) {
	rootCertPool := x509.NewCertPool()

	pem, err := os.ReadFile(os.Getenv("PEM_PATH"))
	if err != nil {
		t.Errorf("cannot load pem: %v", err)
	}
	if ok := rootCertPool.AppendCertsFromPEM(pem); !ok {
		t.Errorf("cannot append pem")
	}

	if err = mysql.RegisterTLSConfig("custom", &tls.Config{
		RootCAs:            rootCertPool,
		InsecureSkipVerify: true,
	}); err != nil {
		t.Errorf("cannnot register tls")
	}
}
