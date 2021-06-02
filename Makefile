mock:
	mockery --dir p2p/ --name IConnection --output p2p/mock --exported

test:
	go test -v -cover -count=1 ./p2p/...