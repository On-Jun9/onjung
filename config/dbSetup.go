package config

import (
	"errors"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB
var Driver DBDriver

// DBDriver는 데이터베이스 드라이버 상수를 정의합니다.
type DBDriver string

const (
	// PostgreSQL 드라이버
	PostgreSQL DBDriver = "postgres"
	// MySQL 드라이버
	MySQL DBDriver = "mysql"
	// Oracle 드라이버
	Oracle DBDriver = "oracle"
)

// DBConnector 는 데이터베이스 연결 인터페이스입니다.
type DBConnector interface {
	Connect() (*gorm.DB, error)
}

// PostgreSQLConnector 는 PostgreSQL 데이터베이스 연결 구조체입니다.
type PostgreSQLConnector struct{}

// Connect 은 PostgreSQL 데이터베이스에 연결합니다.
func (p *PostgreSQLConnector) Connect() (*gorm.DB, error) {
	Driver = PostgreSQL

	host := RuntimeConf.Datasource.Host
	user := RuntimeConf.Datasource.User
	password := RuntimeConf.Datasource.Password
	dbname := RuntimeConf.Datasource.Name
	port := RuntimeConf.Datasource.Port
	sslMode := RuntimeConf.Datasource.SslMode

	dsn := "host=" + host + " user=" + user + " password=" + password + " dbname=" + dbname + " port=" + port + " sslmode=" + sslMode + " TimeZone=Asia/Seoul"
	return gorm.Open(postgres.Open(dsn), &gorm.Config{
		Logger: logger.Default.LogMode(logger.LogLevel(RuntimeConf.Server.DBLogLevel)),
	})
}

// MySQLConnector 는 MySQL 데이터베이스 연결 구조체입니다.
type MySQLConnector struct{}

// Connect 은 MySQL 데이터베이스에 연결합니다.
func (m *MySQLConnector) Connect() (*gorm.DB, error) {
	Driver = MySQL

	// MySQL에 대한 연결 설정을 추가하세요.
	// 예: dsn := "user:password@tcp(host:port)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
	// return gorm.Open(mysql.Open(dsn), &gorm.Config{ ... })
	return nil, nil
}

// OracleConnector 는 Oracle 데이터베이스 연결 구조체입니다.
type OracleConnector struct{}

// Connect 은 Oracle 데이터베이스에 연결합니다.
func (o *OracleConnector) Connect() (*gorm.DB, error) {
	Driver = Oracle

	// Oracle에 대한 연결 설정을 추가하세요.
	// return gorm.Open(oracle.Open(dsn), &gorm.Config{ ... })
	return nil, nil
}

// ConnectDatabase 는 설정 파일에서 정보를 읽어와 DB에 연결하는 함수입니다.
func ConnectDatabase() error {
	driver := RuntimeConf.Datasource.Driver

	var connector DBConnector
	switch DBDriver(driver) {
	case PostgreSQL:
		connector = &PostgreSQLConnector{}
	case MySQL:
		connector = &MySQLConnector{}
	case Oracle:
		connector = &OracleConnector{}
	default:
		return errors.New("unsupported driver")
	}

	gormdb, err := connector.Connect()
	if err != nil {
		return err
	}

	DB = gormdb
	return nil
}
