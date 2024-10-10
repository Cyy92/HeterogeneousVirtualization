package function

import (
	"bytes"
	"io/ioutil"

	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	"github.com/Cyy92/HeterogeneousVirtualization/faas-cli/config"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
)

var (
	authURL      string
	clientID     string
	clientSecret string
)

type ClientCredentialsReq struct {
	ClientID     string `json:"client_id"`
	ClientSecret string `json:"client_secret"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
	Scope       string `json:"scope"`
	TokenType   string `json:"token_type"`
}

func init() {
	authCmd.Flags().StringVar(&authURL, "auth-url", config.DefaultOAuth2Server, "OAuth2 Authorize URL i.e. Openfx-oauth2 git")
	authCmd.Flags().StringVarP(&clientID, "client-id", "", "", "OAuth2 client_id")
	authCmd.Flags().StringVarP(&clientSecret, "client-secret", "", "", "OAuth2 client_secret, for use with client_credentials grant")
	authCmd.MarkFlagRequired("client-id")
}

var authCmd = &cobra.Command{
	Use:   `login --auth-url `,
	Short: "Get a Accesstoken for your OpenFx Oauth2 Server",
	Long: `
	Get a Accesstoken for Authentication procedure to access Openfx
	`,
	Example: `faas-cli function login --client-id <client-id> --client-secret <client-secret> 
    `,
	PreRunE: preRunAuth,
	RunE: func(cmd *cobra.Command, args []string) error {
		if err := runAuth(); err != nil {
			return err
		}
		return nil
	},
}

func preRunAuth(cmd *cobra.Command, args []string) error {
	return checkValues(authURL, clientID, clientSecret)
}

func checkValues(authURL, clientID, clientSecret string) error {
	if len(authURL) == 0 {
		return fmt.Errorf("--auth-url is required and must be a valid Openfx OAuth2 Server")
	}

	u, uErr := url.Parse(authURL)
	if uErr != nil {
		return fmt.Errorf("--auth-url is an invalid URL: %s", uErr.Error())
	}

	if u.Scheme != "http" && u.Scheme != "https" {
		return fmt.Errorf("--auth-url is an invalid URL: %s", u.String())
	}

	if len(clientID) == 0 {
		return fmt.Errorf("--client-id is required")
	}

	if len(clientSecret) == 0 {
		return fmt.Errorf("--clientSecret is required")
	}
	return nil
}

// auth server 키고 handler 처리(state 받고 login 에 자동으로 url 이 전송되게끔)도 해야함
func runAuth() error {

	body := ClientCredentialsReq{
		ClientID:     clientID,
		ClientSecret: clientSecret,
	}

	bodyBytes, marshalErr := json.Marshal(body)
	if marshalErr != nil {
		return errors.Wrapf(marshalErr, "unable to unmarshal %s", string(bodyBytes))
	}

	resp1, err := http.Post(authURL+"/credentials", "application/json", bytes.NewBuffer(bodyBytes))
	if err != nil {
		return fmt.Errorf("error making POST request: %s", err)
	}
	defer resp1.Body.Close()

	params := url.Values{}
	params.Add("grant_type", "client_credentials")
	params.Add("client_id", clientID)
	params.Add("client_secret", clientSecret)
	params.Add("scope", "read")

	apiUrl := fmt.Sprintf("http://10.0.2.101:9096/token?%s", params.Encode())

	resp2, err := http.Get(apiUrl)
	if err != nil {
		return fmt.Errorf("error making GET request: %s", err)
	}
	defer resp2.Body.Close()

	resp2body, err := ioutil.ReadAll(resp2.Body)
	if err != nil {
		return fmt.Errorf("error reading response body: %s", err)
	}

	// JSON 응답을 구조체로 언마샬링
	var tokenResponse TokenResponse
	err = json.Unmarshal(resp2body, &tokenResponse)
	if err != nil {
		return fmt.Errorf("error unmarshalling JSON: %s", err)
	}

	fmt.Println("Access Token:", tokenResponse.AccessToken)
	/*
		// url parsing 필요

		relativeUrl := "/token"
		u, err := url.Parse(relativeUrl)
		if err != nil {
			log.Fatal(err)
		}

		queryString := u.Query()

		queryString.Set("grant_type", "client_credentials")
		queryString.Set("client_id", clientID)
		queryString.Set("client_secret", clientSecret)
		queryString.Set("scope", "read")

		u.RawQuery = queryString.Encode()

		base, err := url.Parse(authURL)
		if err != nil {
			log.Fatal(err)
		}

		req, _ := http.NewRequest(http.MethodGet, base.ResolveReference(u).String(), buf)

		req.Header.Set("Content-Type", "application/json")

		res, err := http.DefaultClient.Do(req)

		if err != nil {
			return errors.Wrap(err, fmt.Sprintf("cannot POST to %s", authURL))
		}
		if res.Body != nil {
			defer res.Body.Close()

			tokenData, _ := ioutil.ReadAll(res.Body)

			if res.StatusCode != http.StatusOK {
				// 에러 메세지 업데이트 필요
				return fmt.Errorf("[Information Error] The client information is incorrect.  %s")
			}
			token := token_data{}
			tokenErr := json.Unmarshal(tokenData, &token)
			if tokenErr != nil {
				return errors.Wrapf(tokenErr, "unable to unmarshal token: %s", string(tokenData))
			}
			config.UpdateAuthConfig(clientID, clientSecret, token.Access_token)

			log.Println("successfully completed the certification.\n")
		}
	*/
	return nil
}
