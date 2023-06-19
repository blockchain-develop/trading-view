package generate_data

import (
	"encoding/json"
	"github.com/shopspring/decimal"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"os"
	"testing"
)

type PriceX struct {
	Id    int64           `gorm:"primaryKey;autoIncrement" json:"-"`
	Time  uint64          `gorm:"type:bigint(20);not null" json:"time"`
	Open  decimal.Decimal `gorm:"type:varchar(32);not null" json:"open"`
	High  decimal.Decimal `gorm:"type:varchar(32);not null" json:"high"`
	Low   decimal.Decimal `gorm:"type:varchar(32);not null" json:"low"`
	Close decimal.Decimal `gorm:"type:varchar(32);not null" json:"close"`
}

func TestData(t *testing.T) {
	data, err := os.ReadFile("./price.json")
	if err != nil {
		panic(err)
	}
	//
	prices := make([]*PriceX, 0)
	err = json.Unmarshal(data, &prices)
	if err != nil {
		panic(err)
	}
	//
	Logger := logger.Default
	if true {
		Logger = Logger.LogMode(logger.Info)
	}
	user := "root"
	password := "Tgy_#0010"
	url := "127.0.0.1"
	scheme := "chart"
	db, err := gorm.Open(mysql.Open(user+":"+password+"@tcp("+url+")/"+
		scheme+"?charset=utf8"), &gorm.Config{Logger: Logger})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&PriceX{})
	//
	db.Save(prices)
}
