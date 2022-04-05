package main

import (
	"encoding/json"
	"github.com/gookit/color"
	"github.com/zLeki/RoShield/fetchers"
	"io/ioutil"
	"os"
	"strconv"
	"time"
)

type Config struct {
	Warning  bool   `json:"warning"`
	Warnings int    `json:"warnings"`
	Cookie   string `json:"cookie"`
	Groupid  int    `json:"groupid"`
}

func CheckConfig() Config {
	dir, _ := os.ReadDir("./")
	for _, v := range dir {
		if v.Name() == "config.json" {
			color.Notice.Tips("Config file found")
			plan, _ := ioutil.ReadFile("config.json")
			var data Config
			err := json.Unmarshal(plan, &data)
			if err != nil {
				color.Error.Tips("Config file is not valid")
				os.Exit(1)
			}
			return data
		}
	}
	color.Warn.Tips("config.json not found.. creating")
	f, _ := os.Create("config.json")
	defer func(f *os.File) {
		err := f.Close()
		if err != nil {
			color.Error.Tips("Fatal error detected while encoding config:", err)
			os.Exit(3)
		}
	}(f)
	config := Config{
		Warning:  true,
		Warnings: 5,
		Cookie:   "insert your cookie here and group id below",
		Groupid:  0,
	}
	encoder := json.NewEncoder(f)
	encoder.SetIndent("", "  ")
	err := encoder.Encode(config)
	if err != nil {
		color.Error.Tips("Fatal error detected while encoding config:", err)
		os.Exit(3)
	}
	color.Success.Tips("config.json created, please fill out the nessary information")
	time.Sleep(time.Second * 15)
	os.Exit(3)
	return Config{}
}
func main() {
	config := CheckConfig()
	if config.Warning {
		if config.Warnings >= 5 {
			color.Warn.Tips("Un-safe amount of warnings in the config.json file.")
		}
	}
	color.Info.Tips("Starting...")
	color.Info.Tips("Cookie: " + config.Cookie)
	color.Info.Tips("Group ID: " + strconv.Itoa(config.Groupid))
	color.Info.Tips("Warning: " + strconv.FormatBool(config.Warning))
	color.Info.Tips("Warnings:" + strconv.Itoa(config.Warnings))
retry:
	roles, err := fetchers.FetchRoles(fetchers.Config(config))
	if err != nil {
		color.Error.Tips("Ignored error while fetching roles:", err)
		goto retry
	}
	if len(roles.Roles) == 0 {
		color.Error.Tips("No roles found")
		return
	} else {
		color.Success.Tips("Roles found: " + strconv.Itoa(len(roles.Roles)) + " roles")
		for i := range roles.Roles {
			go func() {
				config.MainProcess(roles, i)
			}()
		}

	}
}
func (c Config) MainProcess(role fetchers.Data, i int) {
	for {
		x := role.Roles[i].MemberCount
		fetchers.FetchRoles(fetchers.Config(c))
		time.Sleep(time.Second * 10)
	}

}
