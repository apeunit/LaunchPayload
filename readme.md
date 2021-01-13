# launchpayload

**launchpayload** is the default blockchain application generated with [Starport](https://github.com/tendermint/starport). It serves as a template of how to package the Cosmos blockchain as a Docker image so that multiple instances thereof can be spun up with [lctrld](https://github.com/apeunit/launchcontrold).

launchpayload includes a simple faucet service under `cmd/faucet`.

Script files save `lctrld` from having to know too many details of how to configure and run any particular payload.
- `cmd/faucet/configurefaucet.sh` and `cmd/faucet/runfaucet.sh` decouples `lctrld` from the details of the faucet
- `runlightclient.sh` decouples `lctrld` from knowing how the light client daemon should be started

Currently, the process of generating the genesis.json and validator node configuration is tightly coupled with `lctrld`. This will be improved in the future.

## Building
Build the binaries to `dist/`
```sh
> make build
```

Build the docker image (only possible after building binaries)
```sh
> make docker
```

## Starport Local Development
### Get started

```
starport serve
```

`serve` command installs dependencies, initializes and runs the application.

### Configure

Initialization parameters of your app are stored in `config.yml`.

### `accounts`

A list of user accounts created during genesis of your application.

| Key   | Required | Type            | Description                                       |
| ----- | -------- | --------------- | ------------------------------------------------- |
| name  | Y        | String          | Local name of the key pair                        |
| coins | Y        | List of Strings | Initial coins with denominations (e.g. "100coin") |
