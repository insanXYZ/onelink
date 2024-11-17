package main

import (
	_ "github.com/go-sql-driver/mysql"
	"radproject/config"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err.Error())
	}

	// config init

	db := config.NewDb()
	echo := config.NewEcho()
	validator := config.NewValidator()

	cfg := config.BootstrapConfigs{
		Validator: validator,
		Echo:      echo,
		Db:        db,
	}

	cfg.Run()
}
