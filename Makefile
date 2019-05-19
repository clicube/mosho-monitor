fmt:
	go fmt ./...

build:
	mkdir -p ./bin
	go build -o ./bin/mosho-monitord ./cmd
	cp basic.txt ./bin

run: build
	cd bin && ./mosho-monitord

clean:
	go clean ./...
	rm -rf ./bin

test:
	go test ./...

remote-build:
	mkdir -p ./bin
	GOOS=linux GOARCH=arm GOARM=6 go build -o ./bin/mosho-monitord ./cmd

remote-copy: remote-build
	ssh raspi "mkdir -p services/mosho-monitord/bin"
	scp bin/mosho-monitord basic.txt raspi:services/mosho-monitord/bin/
	scp mosho-monitord.service raspi:services/mosho-monitord/

remote-install: remote-copy
	ssh raspi "cd services/mosho-monitord && sudo cp mosho-monitord.service /etc/systemd/system && sudo systemctl enable mosho-monitord && sudo service mosho-monitord start"
