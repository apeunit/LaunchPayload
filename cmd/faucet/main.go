package main

import (
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

var faucetConfig *FaucetConfig
var configOutputPath string

const cliBinaryPath = "/payload/launchpayloadcli"
const cliConfigPath = "/home/docker/nodeconfig/faucet_account"

type FaucetConfig struct {
	ListenAddr    string `yaml:"listen_addr"`
	ChainID       string `yaml:"chain_id"`
	CliBinaryPath string `yaml:"cli_binary_path"`
	CliConfigPath string `yaml:"cli_config_path"`
	FaucetAddr    string `yaml:"faucet_addr"`
	Unit          string `yaml:"unit"`
	NodeAddr      string `yaml:"node_addr"`
	Secret        string `yaml:"secret"`
}

func (fc *FaucetConfig) Parse(data []byte) error {
	return yaml.Unmarshal(data, fc)
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

// LoadFaucetConfig loads FaucetConfig from a .yaml file
func LoadFaucetConfig(filename string) (fc *FaucetConfig, err error) {
	f, err := ioutil.ReadFile(filename)
	if err != nil {
		return
	}
	fc = new(FaucetConfig)
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
func RunStatus(fc *FaucetConfig) (output string, err error) {
	c := fmt.Sprintf("%s status --node tcp://%s -o json", fc.CliBinaryPath, fc.NodeAddr)
	return RunCommand(c)
}

// HttpRunStatus is a HTTP wrapper function around RunStatus
func HttpRunStatus(w http.ResponseWriter, r *http.Request) {
	o, err := RunStatus(faucetConfig)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(o))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(o))
}

// RunSendTx runs launchpayloadcli tx send FROM_ADDR DEST_ADDR AMOUNT
func RunSendTx(fc *FaucetConfig, destAddr, amount string) (output string, err error) {
	cliOptions := fmt.Sprintf("--home /payload/config/faucet_account --keyring-backend test --chain-id %s --node tcp://%s -o json", fc.ChainID, fc.NodeAddr)
	cliSend := fmt.Sprintf("%s tx send %s %s %s %s --yes", fc.CliBinaryPath, fc.FaucetAddr, destAddr, amount, cliOptions)
	return RunCommand(cliSend)
}

// HttpRunSendTx is a http wrapper function around RunSendTx that uses a token for authentication
func HttpRunSendTx(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	destAddr := params["destAddr"]
	amount := params["amount"]
	err := r.ParseForm()
	if err != nil {
		log.Fatal(err)
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}
	token := r.Form.Get("token")
	if token != faucetConfig.Secret {
		log.Println("Someone sent the wrong token")
		w.WriteHeader(http.StatusNetworkAuthenticationRequired)
		return
	}

	o, err := RunSendTx(faucetConfig, destAddr, amount)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte(o))
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(o))
}

func startFaucet(c *cobra.Command, args []string) (err error) {
	fc, err := LoadFaucetConfig(args[0])
	faucetConfig = fc
	if err != nil {
		return err
	}

	router := mux.NewRouter()
	router.HandleFunc("/status", HttpRunStatus).Methods("GET")
	router.HandleFunc("/send/{destAddr}/{amount}", HttpRunSendTx).Methods("POST")
	log.Fatal(http.ListenAndServe(fc.ListenAddr, router))

	return nil
}

func generateConfig(c *cobra.Command, args []string) (err error) {
	evtID := args[0]
	faucetAddr := args[1]
	tokenSymbol := args[2]
	nodeIP := args[3]

	fc := FaucetConfig{
		ListenAddr:    "0.0.0.0:8000",
		ChainID:       evtID,
		CliBinaryPath: cliBinaryPath,
		CliConfigPath: cliConfigPath,
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
