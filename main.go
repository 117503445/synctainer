package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/117503445/goutils"
	"github.com/rs/zerolog/log"
	"github.com/simonnilsson/ask"
	"gopkg.in/yaml.v3"
)

func main() {
	goutils.InitZeroLog()

	if err := goutils.CMD(funcDir, "s", "invoke", "--invocation-type", "Async"); err != nil {
		panic(err)
	}
}