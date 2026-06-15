package client

import (
	"os"
	"slices"
	"testing"
)

var client *Client

func init() {
	host := os.Getenv("XUI_HOST")
	token := os.Getenv("XUI_TOKEN")
	basePath := os.Getenv("XUI_BASE_PATH")
	scheme := Scheme(os.Getenv("XUI_SCHEME"))

	client = CreateClient(ClientOptions{
		Host:     &host,
		Token:    token,
		BasePath: &basePath,
		Scheme:   &scheme,
	})
}

func TestClient_GetInbounds(t *testing.T) {
	t.Logf("ENV: %v", os.Environ())

	inbounds, err := client.Inbounds()
	if err != nil {
		t.Errorf("Failed to do request: %s", err)
	}

	t.Logf("Inbounds: %v", inbounds)
}

func TestClient_GetClients(t *testing.T) {
	clients, err := client.GetClientList()
	if err != nil {
		t.Error(err)
	}

	t.Logf("Clients: %v", clients)
}

func TestClient_GetClintsByTgId(t *testing.T) {
	clients, err := client.GetClientList()
	if err != nil {
		t.Error(err)
	}

	for _, cl := range clients {
		if cl.TgID != 0 && cl.ID != -1 {
			tg_clients, err := client.GetClientByTelegramId(int64(cl.TgID))
			if err != nil {
				t.Error(err)
			}

			for _, tgcl := range tg_clients {
				if !slices.Contains(clients, tgcl) {
					t.Errorf("Empty: %d; %v", cl.TgID, clients)
				}
			}
		}
	}
}
