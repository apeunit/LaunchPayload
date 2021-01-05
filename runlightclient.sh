#!/bin/sh
# USAGE: runlightclient.sh NODE_IP CHAIN_ID

# --trust-node to avoid undebuggable nil pointer reference error that only occurs when lctrld issues the run command (works just fine from the shell)
/payload/launchpayloadcli rest-server --laddr tcp://0.0.0.0:1317 --node tcp://$1:26657 --unsafe-cors --trust-node --chain-id $2 --home /payload/config/cli
