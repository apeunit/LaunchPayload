package main

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/stretchr/testify/assert"
)

var tFc = &FaucetConfig{
	ListenAddr:    "127.0.0.1:8000",
	ChainID:       "drop-28b10d4eff415a7b0b2c",
	CliBinaryPath: "/tmp/workspace/bin/launchpayloadcli",
	CliConfigPath: "/tmp/workspace/evts/drop-28b10d4eff415a7b0b2c/nodeconfig/extra_accounts/dropgiver",
	FaucetAddr:    "cosmos1hlquhkh96p67jja733whk2ydea3vhwvjw4tn75",
	Unit:          "drop",
	NodeAddr:      "127.0.0.1:26657",
	Secret:        "abadjoke",
}

func TestLoadFaucetConfig(t *testing.T) {
	fc, err := LoadFaucetConfig("testdata/dropevent.yaml")
	if err != nil {
		t.Error(err)
	}
	expectedFc := &FaucetConfig{
		ListenAddr:    "127.0.0.1:8000",
		ChainID:       "drop-28b10d4eff415a7b0b2c",
		CliBinaryPath: "/tmp/workspace/bin/launchpayloadcli",
		CliConfigPath: "/tmp/workspace/evts/drop-28b10d4eff415a7b0b2c/nodeconfig/extra_accounts/dropgiver",
		FaucetAddr:    "cosmos1hlquhkh96p67jja733whk2ydea3vhwvjw4tn75",
		Unit:          "drop",
		NodeAddr:      "192.168.99.104:26657",
		Secret:        "abadjoke",
	}
	assert.Equal(t, expectedFc, fc)
}

func TestRunCommand(t *testing.T) {
	o, err := RunCommand("uname -a")
	t.Log(o)
	if err != nil {
		t.Error(err)
	}
}

func TestRunStatus(t *testing.T) {
	o, err := RunStatus(tFc)
	if err != nil {
		t.Error(err)
	}
	t.Log(o)
}

func TestHttpRunStatus(t *testing.T) {
	req, err := http.NewRequest("GET", "/status", nil)
	if err != nil {
		t.Fatal(err)
	}

	faucetConfig = tFc
	// We create a ResponseRecorder (which satisfies http.ResponseWriter) to record the response.
	rr := httptest.NewRecorder()
	handler := http.HandlerFunc(HttpRunStatus)

	// Our handlers satisfy http.Handler, so we can call their ServeHTTP method
	// directly and pass in our Request and ResponseRecorder.
	handler.ServeHTTP(rr, req)

	// Check the status code is what we expect.
	if status := rr.Code; status != http.StatusOK {
		t.Errorf("handler returned wrong status code: got %v want %v",
			status, http.StatusOK)
	}

	// Check the response body is what we expect.
	expected := `{"alive": true}`
	if rr.Body.String() != expected {
		t.Errorf("handler returned unexpected body: got %v want %v",
			rr.Body.String(), expected)
	}

}

func testRunSendTx(t *testing.T) {
	o, err := RunSendTx(tFc, "cosmos109cjtu8vaperd7hxfyx90a4p62le4el76jlut0", "1stake")
	if err != nil {
		t.Error(err)
	}
	t.Log(o)
}
