package main

import (
	"github.com/eliofery/go-chix/pkg/config"
	"github.com/eliofery/go-chix/pkg/config/viperr"
	"github.com/eliofery/go-chix/pkg/utils"
)

func main() {
	config := config.MustInit(viperr.New(utils.GetEnv()))
	_ = config
}
