/*
	Copyright (C) 2024  Pagefault Games

	This program is free software: you can redistribute it and/or modify
	it under the terms of the GNU Affero General Public License as published by
	the Free Software Foundation, either version 3 of the License, or
	(at your option) any later version.

	This program is distributed in the hope that it will be useful,
	but WITHOUT ANY WARRANTY; without even the implied warranty of
	MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
	GNU Affero General Public License for more details.

	You should have received a copy of the GNU Affero General Public License
	along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package account

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"net/url"
	"os"
)

func RetrieveDiscordId(code string) (string, error) {
	token, err := http.PostForm("https://discord.com/api/oauth2/token", url.Values{
		"client_id":     {os.Getenv("DISCORD_CLIENT_ID")},
		"client_secret": {os.Getenv("DISCORD_CLIENT_SECRET")},
		"grant_type":    {"authorization_code"},
		"code":          {code},
		"redirect_uri":  {os.Getenv("DISCORD_CALLBACK_URI")},
		"scope":         {"identify"},
	})

	if err != nil {
		log.Println("error getting token:", err)
		return "", err
	}
	log.Println("token: ", token)
	// extract access_token from token
	type TokenResponse struct {
		AccessToken string `json:"access_token"`
	}

	var tokenResponse TokenResponse
	err = json.NewDecoder(token.Body).Decode(&tokenResponse)
	if err != nil {
		return "", err
	}
	access_token := tokenResponse.AccessToken
	log.Printf("access_token: %s", access_token)

	if access_token == "" {
		err = errors.New("access token is empty")
		return "", err
	}

	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://discord.com/api/users/@me", nil)
	if err != nil {
		log.Println("error creating request:", err)
		return "", err
	}
	req.Header.Set("Authorization", "Bearer "+access_token)
	resp, err := client.Do(req)
	if err != nil {
		log.Println("error getting user info:", err)
		return "", err
	}
	defer resp.Body.Close()

	type User struct {
		Id string `json:"id"`
	}
	var user User
	err = json.NewDecoder(resp.Body).Decode(&user)

	log.Println("user", user.Id)

	return user.Id, nil
}
