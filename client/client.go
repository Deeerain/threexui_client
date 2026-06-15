package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"time"

	"github.com/Deeerain/threexui_client/model"
)

var (
	defaultHost = "localhost"
)

type Scheme string

const (
	HTTPScheme  Scheme = "http"
	HTTPSScheme Scheme = "https"
)

type ClientOptions struct {
	Host       *string
	BasePath   *string
	Token      string
	httpClient *http.Client
	Scheme     *Scheme
}

func NewDefaultClientOptions(token string) ClientOptions {
	return ClientOptions{
		Token: token,
	}
}

type Client struct {
	options ClientOptions
}

func CreateClient(options ClientOptions) *Client {
	if options.Host == nil {
		options.Host = &defaultHost
	}

	if options.httpClient == nil {
		jar, err := cookiejar.New(nil)
		if err != nil {
			log.Fatalln(err)
		}

		options.httpClient = &http.Client{
			Jar:     jar,
			Timeout: 30 * time.Second,
		}
	}

	return &Client{
		options: options,
	}
}

func (s *Client) Inbounds() ([]model.XUIInbound, error) {
	url := s.makeUrl("panel", "api", "inbounds", "list")

	resp, err := s.doRequest("GET", url.String(), nil)

	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	var respBody model.XUIResponse[[]model.XUIInbound]

	if err := json.NewDecoder(resp.Body).Decode(&respBody); err != nil {
		return nil, fmt.Errorf("decode error: %s", err)
	}

	if respBody.Obj == nil {
		return nil, fmt.Errorf("failed to get 'obj' from response: %v", respBody)
	}

	return respBody.Obj, nil
}

func (s *Client) Inbound(id int) (*model.XUIInbound, error) {
	url := s.makeUrl("panel", "api", "inbounds", "get", fmt.Sprint(id))

	resp, err := s.doRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	var responseBody model.XUIResponse[*model.XUIInbound]

	if err := json.NewDecoder(resp.Body).Decode(&responseBody); err != nil {
		return nil, fmt.Errorf("encoder error: %w", err)
	}

	return responseBody.Obj, nil
}

func (s *Client) AddClientToInbound(inboundId int, settings model.XUIInboundSettings) error {
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

	var result model.XUIResponse[any]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("api error: %v (status_code: %v)", result.Msg, resp.StatusCode)
	}

	return nil
}

func (s *Client) UpdateClient(clientID string, inboundID int, settings model.XUIInboundSettings) error {
	url := s.makeUrl("panel", "api", "inbounds", "updateClient", clientID)

	settingsString, err := json.Marshal(settings)
	if err != nil {
		return fmt.Errorf("marshal error: %w", err)
	}

	payload := struct {
		Id       int    `json:"id"`
		Settings string `json:"settings"`
	}{
		Id:       inboundID,
		Settings: string(settingsString),
	}
	resp, err := s.doRequest("POST", url.String(), payload)
	if err != nil {
		return fmt.Errorf("request error: %w", err)
	}

	var result model.XUIResponse[any]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return fmt.Errorf("decode error: %w", err)
	}

	if !result.Success {
		return fmt.Errorf("api error: %v (status_code: %v)", result.Msg, resp.StatusCode)
	}

	return nil
}

func (s *Client) ApiTokens() ([]model.TokenInfo, error) {
	url := s.makeUrl("panel", "settings", "apiTokens")

	if resp, err := s.doRequest("POST", url.String(), nil); err == nil {
		var result model.XUIResponse[[]model.TokenInfo]
		if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
			return nil, fmt.Errorf("decode error: %w", err)
		}

		return result.Obj, nil
	} else {
		return nil, fmt.Errorf("request error: %w", err)
	}
}

func (s *Client) GetClientList() ([]model.XUIInboundClient, error) {
	url := s.makeUrl("panel", "api", "clients", "list")

	resp, err := s.doRequest("GET", url.String(), nil)
	if err != nil {
		return nil, fmt.Errorf("Failed do request: %w", err)
	}
	defer resp.Body.Close()

	var result model.XUIResponse[[]model.XUIInboundClient]
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, fmt.Errorf("Failed to decode response: %w", err)
	}

	return result.Obj, nil
}

func (s *Client) GetClientByTelegramId(id int64) ([]model.XUIInboundClient, error) {
	clients, err := s.GetClientList()
	if err != nil {
		return nil, err
	}

	var result []model.XUIInboundClient

	for _, client := range clients {
		if client.TgID == int(id) {
			result = append(result, client)
		}
	}

	return result, err
}

func (s *Client) Login(username string, passwword string) error {
	var creds struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	creds.Username = username
	creds.Password = passwword

	url := s.makeUrl("login")

	_, err := s.doRequest("post", url.String(), creds)
	if err != nil {
		return fmt.Errorf("failed to login to panel: %w", err)
	}

	return nil
}

func (s *Client) httpClient() *http.Client {
	return s.options.httpClient
}

func (s *Client) makeUrl(elem ...string) *url.URL {
	scheme := HTTPScheme

	if s.options.Scheme != nil && *s.options.Scheme != "" {
		scheme = *s.options.Scheme
	}

	url := &url.URL{
		Scheme: string(scheme),
		Host:   *s.options.Host,
	}

	if s.options.BasePath != nil {
		url = url.JoinPath(*s.options.BasePath)
	}

	url = url.JoinPath(elem...)

	return url
}

func (s *Client) doRequest(method string, url string, body any) (*http.Response, error) {
	var bodyBuffer io.Reader

	log.Printf("Request: [%s] %s", method, url)

	if body != nil {
		buf := &bytes.Buffer{}
		if err := json.NewEncoder(buf).Encode(body); err != nil {
			return nil, fmt.Errorf("failed to encode request body: %w", err)
		}
		bodyBuffer = buf
	}

	req, err := http.NewRequest(method, url, bodyBuffer)
	if err != nil {
		return nil, fmt.Errorf("request error: %w", err)
	}

	log.Printf("Requst: %s", req.Body)

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Requested-With", "XMLHttpRequest")
	req.Header.Set("User-Agent", "threexui-client/1.0")
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", s.options.Token))

	resp, err := s.httpClient().Do(req)
	if err != nil {
		return resp, fmt.Errorf("response error:  %w", err)
	}

	if resp.StatusCode != http.StatusOK {
		bodyBytes, _ := io.ReadAll(resp.Body)
		resp.Body.Close()

		resp.Body = io.NopCloser(bytes.NewBuffer(bodyBytes))
		return resp, fmt.Errorf("unexpected status code: %d, body: %v", resp.StatusCode, string(bodyBytes))
	}

	return resp, nil
}
