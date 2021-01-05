FROM debian:stable

# # Copy our static executable + data
COPY dist/launchpayloadd dist/launchpayloadcli dist/faucet /payload/
COPY runlightclient.sh cmd/faucet/runfaucet.sh cmd/faucet/configurefaucet.sh /payload/

RUN mkdir /payload/config
VOLUME /payload/config
USER 1000:50
EXPOSE 1317 8000 26656 26657 26658
# Run the whole shebang.
CMD [ "/payload/launchpayloadd", "start", "--home", "/payload/config/daemon/"]
