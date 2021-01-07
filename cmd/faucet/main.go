package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/gorilla/mux"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
)

var fConfig *faucetConfig
var configOutputPath string

const defaultCLIBinaryPath = "/payload/launchpayloadcli"
const defaultCLIConfigPath = "/home/docker/nodeconfig/faucet_account"

type faucetConfig struct {
	ListenAddr    string `yaml:"listen_addr"`
	ChainID       string `yaml:"chain_id"`
	CliBinaryPath string `yaml:"cli_binary_path"`
	CliConfigPath string `yaml:"cli_config_path"`
	FaucetAddr    string `yaml:"faucet_addr"`
	Unit          string `yaml:"unit"`
	NodeAddr      string `yaml:"node_addr"`
	Secret        string `yaml:"secret"`
}

func (f *faucetConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, f)
}

// SendRequest represents the send request that will come in as JSON
type SendRequest struct {
	ToAddr string `json:"to_address"`
	Amount string `json:"amount"`
	Memo   string `json:"memo"`
	Token  string `json:"token"`
}

func main() {
	var rootCmd = &cobra.Command{
		Use:   "faucet <configfile.yml>",
		Short: "A faucet to dispense some drops using launchpayloadcli",
		Args:  cobra.ExactArgs(1),
		RunE:  startFaucet,
	}

	var genCfgCmd = &cobra.Command{
		Use:   "generate-config EVENT_ID FAUCET_ADDR TOKEN_SYMBOL NODE_IP",
		Short: "Generate a configuration file from the arguments",
		Args:  cobra.ExactArgs(4),
		RunE:  generateConfig,
	}

	pwd, err := filepath.Abs(".")
	if err != nil {
		return
	}
	genCfgCmd.Flags().StringVarP(&configOutputPath, "output", "o", filepath.Join(pwd, "faucetconfig.yml"), "filename to write the config to")
	rootCmd.AddCommand(genCfgCmd)
	err = rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

// loadFaucetConfig loads FaucetConfig from a .yaml file
func loadFaucetConfig(filename string) (fc *faucetConfig, err error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	fc = new(faucetConfig)
	err = yaml.Unmarshal(f, fc)
	if err != nil {
		return
	}
	return
}

// RunCommand is a general purpose CLI program running function
func RunCommand(c string) (output string, err error) {
	csplit := strings.Split(c, " ")
	log.Println("RunCommand", c)
	out, err := exec.Command(csplit[0], csplit[1:]...).CombinedOutput()
	output = string(out)
	return
}

// RunStatus runs launchpayloadcli status --node tcp://.... -o json
func RunStatus(fc *faucetConfig) (output string, err error) {
	c := fmt.Sprintf("%s status --node tcp://%s -o json", fc.CliBinaryPath, fc.NodeAddr)
	return RunCommand(c)
}

// HTTPRunStatus is a HTTP wrapper function around RunStatus
func HTTPRunStatus(w http.ResponseWriter, r *http.Request) {
	o, err := RunStatus(fConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(o))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(o))
}

// RunSendTx runs launchpayloadcli tx send FROM_ADDR DEST_ADDR AMOUNT
func RunSendTx(fc *faucetConfig, destAddr, amount string) (output string, err error) {
	cliOptions := fmt.Sprintf("--home %s --keyring-backend test --chain-id %s --node tcp://%s -o json", fc.CliConfigPath, fc.ChainID, fc.NodeAddr)
	cliSend := fmt.Sprintf("%s tx send %s %s %s %s --yes", fc.CliBinaryPath, fc.FaucetAddr, destAddr, amount, cliOptions)
	return RunCommand(cliSend)
}

// HTTPRunSendTx is a http wrapper function around RunSendTx that uses a token for authentication
func HTTPRunSendTx(w http.ResponseWriter, r *http.Request) {
	headerContentType := r.Header.Get("Content-Type")
	if headerContentType != "application/json" {
		errorResponse(w, "Content-Type was not set to application/json", http.StatusUnsupportedMediaType)
	}

	var req SendRequest
	var unmarshalErr *json.UnmarshalTypeError

	decoder := json.NewDecoder(r.Body)
	decoder.DisallowUnknownFields()

	err := decoder.Decode(&req)
	if err != nil {
		if errors.As(err, &unmarshalErr) {
			errorResponse(w, "Bad Request. Wrong Type for field "+unmarshalErr.Field, http.StatusBadRequest)
		} else {
			errorResponse(w, "Bad Request "+err.Error(), http.StatusBadRequest)
		}
		return
	}

	if req.Token != fConfig.Secret {
		log.Println("Someone sent the wrong token")
		errorResponse(w, "Wrong authentication token", http.StatusNetworkAuthenticationRequired)
		return
	}

	o, err := RunSendTx(fConfig, req.ToAddr, req.Amount)
	if err != nil {
		errorResponse(w, fmt.Sprintf("An error occurred while running the CLI tx send command: %s", err), http.StatusInternalServerError)
		return
	}
	errorResponse(w, o, http.StatusOK)
	return

}

func errorResponse(w http.ResponseWriter, message string, httpStatusCode int) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(httpStatusCode)
	resp := make(map[string]string)
	resp["message"] = message
	jsonResp, _ := json.Marshal(resp)
	w.Write(jsonResp)
}

func startFaucet(c *cobra.Command, args []string) (err error) {
	fc, err := loadFaucetConfig(args[0])
	fConfig = fc
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/status", HTTPRunStatus).Methods("GET")
	router.HandleFunc("/send", HTTPRunSendTx).Methods("POST")
	log.Fatal(http.ListenAndServe(fc.ListenAddr, router))

	return nil
}

func generateConfig(c *cobra.Command, args []string) (err error) {
	evtID := args[0]
	faucetAddr := args[1]
	tokenSymbol := args[2]
	nodeIP := args[3]

	fc := faucetConfig{
		ListenAddr:    "0.0.0.0:8000",
		ChainID:       evtID,
		CliBinaryPath: defaultCLIBinaryPath,
		CliConfigPath: defaultCLIConfigPath,
		FaucetAddr:    faucetAddr,
		Unit:          tokenSymbol,
		NodeAddr:      fmt.Sprintf("%s:26657", nodeIP),
		Secret:        "abadjoke",
	}

	fBytes, err := yaml.Marshal(fc)
	if err != nil {
		return
	}
	fmt.Println("Writing to", configOutputPath)
	err = ioutil.WriteFile(configOutputPath, fBytes, 0644)
	return
}
