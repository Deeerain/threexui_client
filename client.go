package threexuiclient

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"
)

type XUIInboundClient struct {
	ID    string `json:"id"`
	Email string `json:"email"`
	Flow  string `json:"flow"`
	TgID  int    `json:"tgId"`
	SubID string `json:"subId"`
}

type Client struct {
	host       string
	port       int
	secretPath string
	client     http.Client
}

func CreateClient(host string, port int, secretPath *string) *Client {
	var computedString string
	if secretPath == nil {
		computedString = "/secret"
	}

	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Fatalln(err)
	}

	return &Client{
		host:       host,
		port:       port,
		secretPath: computedString,
		client: http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		},
	}
}

func (s *Client) Login(creds *LoginRequest) error {
	url := s.makeUrl("login")
	var responseBody XUIResponse[XUIInbound]

	resp, err := s.doRequest("POST", url.String(), creds)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return fmt.Errorf("decode err: %w", err)
	}

	if !responseBody.Success {
		return fmt.Errorf("Login error: %s", responseBody.Msg)
	}

	return nil
}

func (s *Client) Inbounds() ([]XUIInbound, error) {
	url := s.makeUrl("panel", "api", "inbounds", "list")

	resp, err := s.doRequest("GET", url.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	var respBody XUIResponse[[]XUIInbound]

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}

	return respBody.Obj, nil
}

func (s *Client) Inbound(id int) (*XUIInbound, error) {
	url := s.makeUrl("panel", "api", "inbounds", "get", string(id))

	resp, err := s.doRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	var responseBody XUIResponse[*XUIInbound]

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("encoder error: %w", err)
	}

	return responseBody.Obj, nil
}

func (s *Client) AddClientToInbound(inboundId int, settings XUIInboundSettings) error {
	url := s.makeUrl("panel", "api", "inbounds", "addClient")

	settingsString, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	payload := struct {
		Id       int    `json:"id"`
		Settings string `json:"settings"`
	}{
		Id:       inboundId,
		Settings: string(settingsString),
	}

	resp, err := s.doRequest("POST", url.String(), payload)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	var result XUIResponse[any]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("api error: %v (status_code: %v)", result.Msg, resp.StatusCode)
	}

	return nil
}

func (s *Client) makeUrl(elem ...string) *url.URL {
	url := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%v", s.host, s.port),
	}

	url = url.JoinPath(s.secretPath)
	url = url.JoinPath(elem...)

	return url
}

func (s *Client) doRequest(method string, url string, body any) (*http.Response, error) {
	var bodyBuffer *bytes.Buffer

	if body != nil {
		bodyBuffer = bytes.NewBuffer(nil)
		json.NewEncoder(bodyBuffer).Encode(body)
	}

	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	req.Header.Add("Content-Type", "application/json")

	resp, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("response error:  %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %v", resp.StatusCode)
	}

	return resp, nil
}
