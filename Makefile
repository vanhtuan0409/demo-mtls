gen_ca:
	step-cli certificate create ca.company.org certs/ca.crt certs/ca.key \
		--profile root-ca \
		--not-after="2200-01-01T00:00:00Z" \
		--insecure --no-password
	cp -f ./certs/ca.crt ./cmds/client/ca.crt

gen_intermediate:
	step-cli certificate create intermediate.company.org certs/intermediate.crt certs/intermediate.key \
		--profile intermediate-ca \
		--ca ./certs/ca.crt --ca-key ./certs/ca.key \
		--not-after="2200-01-01T00:00:00Z" \
		--insecure --no-password

gen_server_cert:
	step-cli certificate create server.company.org certs/server.crt certs/server.key \
		--profile leaf \
		--ca ./certs/intermediate.crt --ca-key ./certs/intermediate.key \
		--san=127.0.0.1 \
		--insecure --no-password
	cat certs/server.crt certs/intermediate.crt > certs/bundle_server.crt

gen_client_cert:
	step-cli certificate create client.yolo certs/client.crt certs/client.key \
		--profile leaf \
		--ca ./certs/intermediate.crt --ca-key ./certs/intermediate.key \
		--not-after="2026-01-01T00:00:00Z" \
		--insecure --no-password
	cat certs/client.crt certs/intermediate.crt > certs/bundle_client.crt

gen_unauthenticated_client_cert:
	step-cli certificate create some_othe.client certs/unauthenticated.crt certs/unauthenticated.key \
		--profile leaf \
		--ca ./certs/intermediate.crt --ca-key ./certs/intermediate.key \
		--not-after="2026-01-01T00:00:00Z" \
		--insecure --no-password
	cat certs/unauthenticated.crt certs/intermediate.crt > certs/bundle_unauthenticated.crt

build_server:
	go build -o bin/server ./cmds/server/*.go

build_client:
	go build -o bin/client ./cmds/client/*.go
