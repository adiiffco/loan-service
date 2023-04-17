package mysql

import (
	"fmt"
	"time"

	"github.com/spf13/viper"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

type ORM struct {
	DB *gorm.DB
}

var dbORM ORM

func Initialize() {
	fmt.Println("Initializing mysql")
	connectionString := fmt.Sprintf("%s:%s@tcp(%s:3306)/%s",
		viper.GetString("MYSQL_USER"),
		viper.GetString("MYSQL_PASSWORD"),
		viper.GetString("MYSQL_HOST"),
		viper.GetString("MYSQL_NAME"))

	config := &gorm.Config{}
	if viper.GetString("DEPLOY_ENV") == "prod" {
		config.Logger = logger.Default.LogMode(logger.Silent)
	}

	mysqlDB, err := gorm.Open(mysql.Open(connectionString), config)
	if err != nil {
		panic(fmt.Errorf("error while creating mysql connection: %s", err.Error()))
	}

	dbORM.DB = mysqlDB
	dbORM.SetConnectionPool()
}

func GetDbInstance() *gorm.DB {
	return dbORM.DB
}

func (d *ORM) SetConnectionPool() {
	db, err := d.DB.DB()
	if err != nil {
		panic(fmt.Errorf("error while fetching mysql connection: %s", err.Error()))
	}
	db.SetMaxOpenConns(50)
	db.SetMaxIdleConns(25)
	db.SetConnMaxLifetime(2 * time.Minute)
}
