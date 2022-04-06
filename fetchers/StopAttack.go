package fetchers

import (
	"encoding/json"
	"github.com/gookit/color"
	"github.com/zLeki/Goblox/account"
	csrf2 "github.com/zLeki/Goblox/csrf"
	"net/http"
	"strconv"
)

type Audit struct {
	Data []struct {
		Actor struct {
			User struct {
				UserID   int    `json:"userId"`
				Username string `json:"username"`
			} `json:"user"`
		} `json:"actor"`
		ActionType string `json:"actionType"`
	} `json:"data"`
}

func StopAttack(config Config) (bool, error) {
	req, _ := http.NewRequest("GET", "https://groups.roblox.com/v1/groups/"+strconv.Itoa(config.Groupid)+"/audit-log?cursor=&limit=50&sortOrder=Asc", nil)
	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: config.Cookie})
	resp, _ := http.DefaultClient.Do(req)
	var audit Audit
	err := json.NewDecoder(resp.Body).Decode(&audit)
	if err != nil {
		color.Error.Tips("Fatal error detected while encoding config:", err)
		return false, err
	}
	for _, v := range audit.Data {
		if v.ActionType == "Remove Member" || v.ActionType == "Change Rank" {
			reqwest, _ := http.NewRequest("DELETE", "https://groups.roblox.com/v1/groups/"+strconv.Itoa(config.Groupid)+"/users/"+strconv.Itoa(v.Actor.User.UserID), nil)
			reqwest.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: config.Cookie})
			csrf, _ := csrf2.GetCSRF(account.Validate(config.Cookie))
			reqwest.Header.Set("X-CSRF-TOKEN", csrf)
			resp, _ := http.DefaultClient.Do(reqwest)
			if resp.StatusCode == 200 {
				color.Success.Tips("Attack prevented:", v.Actor.User.Username)
				color.Notice.Tips("Created by Leki#6796")
				return true, nil
			} else {
				color.Error.Tips("Failed to remove user:", v.Actor.User.Username+" rip lmao")
				return true, nil
			}
		}
	}
	return true, nil

}
