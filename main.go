package main
import (
	"encoding/json"
	"github.com/gookit/color"
	"os"
	"time"
)
type Config struct {
	Warning  bool   `json:"warning"`
	Warnings int    `json:"warnings"`
	Cookie   string `json:"cookie"`
	Groupid  int    `json:"groupid"`
}
func init() {
	dir, _ := os.ReadDir("./")
	for _,v := range dir {
		if v.Name() == "config.json" {
			f, _ := os.Open("config.json")
			defer f.Close()
			var config Config
			decoder := json.NewDecoder(f).Decode(&config)
			if decoder != nil {
				color.Error.Tips("Error:", decoder)
				return
			}
			return
		}
	}
	color.Warn.Tips("config.json not found.. creating")
	f, _ := os.Create("config.json")
	defer f.Close()
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
		color.Error.Tips("Fatal error detected while encoding config:", err);os.Exit(3)
	}
	color.Success.Tips("config.json created, please fill out the nessary information")
	time.Sleep(time.Second * 15)
	os.Exit(3)
}
func main() {
	config := Config{}
	if config.Warning {
		if config.Warnings >= 5 {
			color.Warn.Tips("Un-safe amount of warnings in the config.json file.")
		}
	}

}

