package lib

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/pkg/errors"
)

func GetMastodonTokens(code string) (*MastodonAuthResponse, error) {
	data := url.Values{
		"code":          {code},
		"grant_type":    {"authorization_code"},
		"client_id":     {os.Getenv("VITE_MASTODON_CLIENT_ID")},
		"client_secret": {os.Getenv("MASTODON_CLIENT_SECRET")},
		"redirect_uri":  {os.Getenv("VITE_REDIRECT_URI")},
		"scope":         {"read write follow"},
	}

	log.Println(data)

	req, err := http.NewRequest("POST", "https://fosstodon.org/oauth/token", strings.NewReader(data.Encode()))
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) http.NewRequest")
	}

	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) ioutil.ReadAll")
	}

	var authResp MastodonAuthResponse
	err = json.Unmarshal([]byte(bodyText), &authResp)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonTokens) json.Unmarshal")
	}
	return &authResp, nil
}

type MastodonAuthResponse struct {
	AccessToken string `json:"access_token"`
	CreatedAt   int    `json:"created_at"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func GetMastodonUserDetails(token string) (*MastodonGetUserResponse, error) {
	req, err := http.NewRequest("GET", "https://fosstodon.org/api/v1/accounts/verify_credentials", nil)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) http.NewRequest")
	}
	req.Header.Add("Authorization", fmt.Sprintf("Bearer %v", token))
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) client.Do")
	}
	defer resp.Body.Close()
	bodyText, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) ioutil.ReadAll")
	}

	var response MastodonGetUserResponse
	err = json.Unmarshal([]byte(bodyText), &response)
	if err != nil {
		return nil, errors.Wrap(err, "(GetMastodonUserDetails) json.Unmashal")
	}
	return &response, nil
}

type MastodonGetUserResponse struct {
	ID          string `json:"id"`
	Username    string `json:"username"`
	Acct        string `json:"acct"`
	DisplayName string `json:"display_name"`
	Locked      bool   `json:"locked"`
	Bot         bool   `json:"bot"`
	// CreatedAt      time.Time `json:"created_at"`
	Note           string `json:"note"`
	URL            string `json:"url"`
	Avatar         string `json:"avatar"`
	AvatarStatic   string `json:"avatar_static"`
	Header         string `json:"header"`
	HeaderStatic   string `json:"header_static"`
	FollowersCount int    `json:"followers_count"`
	FollowingCount int    `json:"following_count"`
	StatusesCount  int    `json:"statuses_count"`
	// LastStatusAt   time.Time `json:"last_status_at"`
	Source struct {
		Privacy   string `json:"privacy"`
		Sensitive bool   `json:"sensitive"`
		Language  string `json:"language"`
		Note      string `json:"note"`
		Fields    []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
			// VerifiedAt time.Time `json:"verified_at"`
		} `json:"fields"`
		FollowRequestsCount int `json:"follow_requests_count"`
	} `json:"source"`
	Emojis []struct {
		Shortcode       string `json:"shortcode"`
		URL             string `json:"url"`
		StaticURL       string `json:"static_url"`
		VisibleInPicker bool   `json:"visible_in_picker"`
	} `json:"emojis"`
	Fields []struct {
		Name       string    `json:"name"`
		Value      string    `json:"value"`
		VerifiedAt time.Time `json:"verified_at"`
	} `json:"fields"`
}
