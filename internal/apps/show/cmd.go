package show

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli/v2"
)

const (
	authUrl   = "https://auth.docker.io/token?service=registry.docker.io&scope=repository:ratelimitpreview/test:pull"
	remainUrl = "https://registry-1.docker.io/v2/ratelimitpreview/test/manifests/latest"
)

var (
	Cmd = &cli.Command{
		Name:  "show",
		Usage: "display docker hub rate limit status",
		Action: func(c *cli.Context) error {
			if err := run(); err != nil{
				return err
			}
			return nil
		},
	}
)

type DockerHubGetResp struct {
	Token       string `json:"token"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in,omitempty"`
	IssuedAt    string `json:"issued_at,omitempty"`
}

type DockerHubLimitParameters struct {
	PullLimit          string `json:"pull_limit"`
	PullLimitInterval  string `json:"pull_limit_interval"`
	PullLimitRemaining string
}

func (d *DockerHubGetResp) GetToken() string {
	if d != nil {
		return d.Token
	}
	return ""
}

func (d *DockerHubGetResp) GetAccessToken() string {
	if d != nil {
		return d.AccessToken
	}
	return ""
}

func (d *DockerHubGetResp) GetExpiresIn() int {
	if d != nil {
		return d.ExpiresIn
	}
	return 0
}

func (d *DockerHubGetResp) GetIssuedAt() string {
	if d != nil {
		return d.IssuedAt
	}
	return ""
}

func httpGet(url, token string) (*http.Response, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}
	if token != "" {
		req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", token))
	}

	resp, err := client.Do(req)
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}

	code := resp.StatusCode
	if code > 299 {
		err = fmt.Errorf("HTTP Status Code is %v from %#v", code, url)
		return nil, err
	}

	return resp, nil
}

func GetDockerHubToken() (*DockerHubGetResp, error) {
	log.Debugf("Trying to authenticate to %s", authUrl)

	resp, err := httpGet(authUrl, "")
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}

	defer resp.Body.Close()
	bodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}
	log.Debug(string(bodyBytes))

	var respObj DockerHubGetResp
	json.Unmarshal(bodyBytes, &respObj)
	log.Debugf("%#v", respObj)

	return &respObj, nil
}

func GetDockerHubLimit(token string) (*DockerHubLimitParameters, error) {
	log.Debugf("Trying to access to %s", remainUrl)

	resp, err := httpGet(remainUrl, token)
	if err != nil {
		log.Debugf("%#v", err)
		return nil, err
	}
	log.Debugf("%#v", resp.Header)

	remained := resp.Header["Ratelimit-Remaining"]
	maxAndInter := resp.Header["Ratelimit-Limit"]

	log.Debug(remained)
	log.Debug(maxAndInter)

	//Fortmat
	remained = strings.Split(remained[0], ";w=")
	maxAndInter = strings.Split(maxAndInter[0], ";w=")

	return &DockerHubLimitParameters{
		PullLimit:          maxAndInter[0],
		PullLimitInterval:  maxAndInter[1],
		PullLimitRemaining: remained[0],
	}, nil
}

func run() error {
	log.Debug("Start to parse rate limit on Docker hub")

	// Get Docker hub token
	resp, err := GetDockerHubToken()
	if err != nil {
		log.Debugf("%#v", err)
		return err
	}
	token := resp.GetToken()
	if resp.GetToken() == "" {
		err = fmt.Errorf("Could not get token of Docker hub")
		return err
	}

	// Get limit status
	status, err := GetDockerHubLimit(token)
	if err != nil {
		log.Debugf("%#v", err)
		return err
	}
	fmt.Printf("Maximum Limit: %s\n", status.PullLimit)
	fmt.Printf("Limit Interval: %s\n", status.PullLimitInterval)
	fmt.Printf("Remained: %s\n", status.PullLimitRemaining)
	return nil
}
