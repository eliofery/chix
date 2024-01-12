package main

import (
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/godotenv"
	"github.com/eliofery/go-chix/pkg/utils"
)

func main() {
	config := config.MustInit(godotenv.New(utils.GetEnv()))
	_ = config
}
