# launchpayload

**launchpayload** is the default blockchain application generated with [Starport](https://github.com/tendermint/starport). It serves as a template of how to package the Cosmos blockchain as a Docker image so that multiple instances thereof can be spun up with [LaunchControlD](https://github.com/apeunit/launchcontrold). It is also used as the blockchain backend behind [Eventivize](https://blog.apeunit.com/were-drop-ing-eventivize-ctm-2021/) - context and details as to how everything works together are explained in [this blog post](https://dev.to/apeunit/the-tech-behind-eventivize-drops-3o7k).

launchpayload includes a simple faucet service under `cmd/faucet`.

Script files save `lctrld` from having to know too many details of how to configure and run any particular payload.
- `cmd/faucet/configurefaucet.sh` and `cmd/faucet/runfaucet.sh` decouples `lctrld` from the details of the faucet
- `runlightclient.sh` decouples `lctrld` from knowing how the light client daemon should be started

Currently, the process of generating the genesis.json and validator node configuration is tightly coupled with `lctrld`. This will be improved in the future.

## Usage
Build the binaries to `dist/`
```sh
> make build
```

Build the docker image (only possible after building binaries)
```sh
> make docker
```

Push to your docker repository and follow the instructions in [LaunchControlD](https://github.com/apeunit/launchcontrold) to tell `lctrld` to deploy your image to the virtual machines.
