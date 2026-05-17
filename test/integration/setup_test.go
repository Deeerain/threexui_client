package integration

import (
	"os"
	"testing"

	threexuiclient "github.com/Deeerain/threexui_client"
)

var testClient *threexuiclient.Client

func TestMain(m *testing.M) {
	token := "CnAktTL5TvlAuBdaeBRwIb0zzQkN0nKRLupmztGmgup0wmy8"
	testClient = threexuiclient.CreateClient(threexuiclient.NewDefaultClientOptions(token))

	code := m.Run()

	os.Exit(code)
}

// func TestLogin_Integration(t *testing.T) {
// 	if err := testClient.Login("admin", "admin"); err != nil {
// 		t.Fatal(err)
// 	}
// }

func TestGetInbounds_Integration(t *testing.T) {
	if inbounds, err := testClient.Inbounds(); err != nil {
		t.Fatal(err)
	} else {
		t.Log("Inbounds", inbounds)
	}
}
