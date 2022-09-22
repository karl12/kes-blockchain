genproto: clean
	protoc --go_out=. \
	--go-grpc_out=. \
	-I=$(PWD) pb/node/v1/*.proto

clean:
	rm -rf gen