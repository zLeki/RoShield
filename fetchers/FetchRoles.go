package fetchers

import (
	"encoding/json"
	"net/http"
	"strconv"
)

type Data struct {
	Roles []struct {
		ID          int    `json:"id"`
		Name        string `json:"name"`
		MemberCount int    `json:"memberCount"`
	} `json:"roles"`
}
type Config struct {
	Warning  bool
	Warnings int
	Cookie   string
	Groupid  int
}

var client http.Client

func FetchRoles(c Config) (Data, error) {
	req, err := http.NewRequest("GET", "https://groups.roblox.com/v1/groups/"+strconv.Itoa(c.Groupid)+"/roles", nil)
	if err != nil {
		return Data{}, err
	}
	req.AddCookie(&http.Cookie{Name: ".ROBLOSECURITY", Value: c.Cookie})
	resp, err := client.Do(req)
	if resp.StatusCode != 200 {
		return Data{}, err
	}
	var data Data
	err = json.NewDecoder(resp.Body).Decode(&data)
	return data, nil
}
