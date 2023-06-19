package main

import (
	"context"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

var httpServer *http.Server
var db *gorm.DB

func main() {
	//
	ctx, cancel := context.WithCancel(context.Background())
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM, syscall.SIGABRT)
	go shutdown(cancel, quit)
	//
	var err error
	Logger := logger.Default
	if true {
		Logger = Logger.LogMode(logger.Info)
	}
	user := "root"
	password := "Tgy_#0010"
	url := "127.0.0.1"
	scheme := "chart"
	db, err = gorm.Open(mysql.Open(user+":"+password+"@tcp("+url+")/"+
		scheme+"?charset=utf8"), &gorm.Config{Logger: Logger})
	if err != nil {
		panic(err)
	}
	db.AutoMigrate(&Price{})
	//
	startRPC()
	<-ctx.Done()
	stopRPC()
}

func shutdown(cancel context.CancelFunc, quit <-chan os.Signal) {
	osCall := <-quit
	fmt.Printf("System call: %v, auto trader is shutting down......\n", osCall)
	cancel()
}

func cors() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method
		origin := c.Request.Header.Get("Origin")
		if origin != "" {
			c.Header("Access-Control-Allow-Origin", "*") // 可将将 * 替换为指定的域名
			c.Header("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE, UPDATE")
			c.Header("Access-Control-Allow-Headers", "Origin, X-Requested-With, Content-Type, Accept, Authorization")
			c.Header("Access-Control-Expose-Headers", "Content-Length, Access-Control-Allow-Origin, Access-Control-Allow-Headers, Cache-Control, Content-Language, Content-Type")
			c.Header("Access-Control-Allow-Credentials", "true")
		}
		if method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
		}
		c.Next()
	}
}

func startRPC() {
	router := gin.New()
	router.Use(cors())
	g := router.Group("/api")
	g.POST("/price", getPrice)
	//
	httpServer = &http.Server{
		Addr:    "127.0.0.1:8000",
		Handler: router,
	}
	go func() {
		if err := httpServer.ListenAndServe(); err != nil {
			fmt.Printf("ListenAndServe: %s", err.Error())
		}
	}()
}

func stopRPC() {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := httpServer.Shutdown(ctx); err != nil {
		panic(err)
	}
	fmt.Printf("rpc server has stopped......")
}

type Price struct {
	Id    int64   `gorm:"primaryKey;autoIncrement" json:"-"`
	Time  uint64  `gorm:"type:bigint(20);not null" json:"time"`
	Open  float32 `gorm:"type:varchar(32);not null" json:"open,float32"`
	High  float32 `gorm:"type:varchar(32);not null" json:"high,float32"`
	Low   float32 `gorm:"type:varchar(32);not null" json:"low,float32"`
	Close float32 `gorm:"type:varchar(32);not null" json:"close,float32"`
}

func getPrice(c *gin.Context) {
	prices := make([]*Price, 0)
	res := db.Find(&prices)
	if res.Error != nil {
		c.JSON(500, res.Error)
		return
	}
	c.JSON(200, prices)
}
