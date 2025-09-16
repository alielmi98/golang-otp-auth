package main

import (
	"log"

	"github.com/alielmi98/golang-otp-auth/internal/user/constants"
	"github.com/alielmi98/golang-otp-auth/migrations"
	"github.com/alielmi98/golang-otp-auth/pkg/cache"
	"github.com/alielmi98/golang-otp-auth/pkg/config"
	"github.com/alielmi98/golang-otp-auth/pkg/db"
)

func main() {

	cfg := config.GetConfig()

	err := cache.InitRedis(cfg)
	defer cache.CloseRedis()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Redis, constants.Startup, err.Error())
	}

	err = db.InitDb(cfg)
	defer db.CloseDb()
	if err != nil {
		log.Fatalf("caller:%s  Level:%s  Msg:%s", constants.Postgres, constants.Startup, err.Error())
	}

	migrations.Up1()

}
