package tui

import (
	"fmt"
	"strings"

	"github.com/117503445/goutils"
	"github.com/117503445/synctainer/cmd/cli/common"
	"github.com/rs/zerolog/log"
)

func loadOld() map[string]any {
	oldCfg := make(map[string]any)
	if goutils.FileExists(common.HomeFileCfg) {
		err := goutils.ReadToml(common.HomeFileCfg, &oldCfg)
		if err != nil {
			log.Fatal().Err(err).Send()
		}
	}
	log.Debug().Interface("oldCfg", oldCfg).Send()
	return oldCfg
}

func tui(cfg map[string]any) {
	if _, ok := cfg["registry"]; !ok {
		cfg["registry"] = make(map[string]any)
	}

	var host string
	if _, ok := cfg["registry"].(map[string]any)["host"]; ok {
		host = cfg["registry"].(map[string]any)["host"].(string)
	}
	if host == "" {
		println("please input registry host, or press enter to use default:", "registry.cn-hangzhou.aliyuncs.com")
	} else {
		println("please input registry host, or press enter to use current:", host)
	}
	var input string
	fmt.Scanln(&input)
	input = strings.TrimSpace(input)
	if input != "" {
		cfg["registry"].(map[string]any)["host"] = input
	} else if host == "" {
		cfg["registry"].(map[string]any)["host"] = "registry.cn-hangzhou.aliyuncs.com"
	} else {
		cfg["registry"].(map[string]any)["host"] = host
	}

	var user string
	if _, ok := cfg["registry"].(map[string]any)["user"]; ok {
		user = cfg["registry"].(map[string]any)["user"].(string)
	}
	if user == "" {
		for {
			println("please input registry user:")
			fmt.Scanln(&user)
			user = strings.TrimSpace(user)
			if user != "" {
				cfg["registry"].(map[string]any)["user"] = user
				break
			} else {
				fmt.Print("user can not be empty, ")
			}
		}
	} else {
		var input string
		println("please input registry user, or press enter to use current user:", user)
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if input != "" {
			cfg["registry"].(map[string]any)["user"] = input
		}
	}

	var pass string
	if _, ok := cfg["registry"].(map[string]any)["pass"]; ok {
		pass = cfg["registry"].(map[string]any)["pass"].(string)
	}
	if pass == "" {
		for {
			println("please input registry pass:")
			fmt.Scanln(&pass)
			pass = strings.TrimSpace(pass)
			if pass != "" {
				cfg["registry"].(map[string]any)["pass"] = pass
				break
			} else {
				fmt.Print("pass can not be empty, ")
			}
		}
	} else {
		var input string
		println("please input registry pass, or press enter to use current pass:", pass)
		fmt.Scanln(&input)
		input = strings.TrimSpace(input)
		if input != "" {
			cfg["registry"].(map[string]any)["pass"] = input
		}
	}
}

func writeNew(new map[string]any) {
	log.Debug().Str("fileCfg", common.HomeFileCfg).Interface("cfg", new).Send()
	err := goutils.WriteToml(common.HomeFileCfg, new)
	if err != nil {
		log.Fatal().Err(err).Send()
	}
}

func EditCfg() {
	log.Debug().Msg("EditCfg")
	cfg := loadOld()

	tui(cfg)

	writeNew(cfg)
}
