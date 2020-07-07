NAME := bank
PKG := github.com/gidyon/$(NAME)
CMD := ${PKG}/cmd
API_IN_PATH := ./api/proto
API_OUT_PATH := ./pkg/api
API_DOC_PATH := ./api/swagger
BINARY := service

setup_dev: ## Sets up a development environment for the banking project
	@cd deployments/compose &&\
	docker-compose up -d

teardown_dev: ## Tear down development environment for the banking project
	@cd deployments/compose &&\
	docker-compose down

protoc_bank:
	@protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/bank bank.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/bank bank.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(API_DOC_PATH) bank.proto

protoc_account:
	@protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/bankaccount bank_account.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/bankaccount bank_account.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(API_DOC_PATH) bank_account.proto

protoc_customer:
	@protoc -I=$(API_IN_PATH) -I=third_party --go_out=plugins=grpc:$(API_OUT_PATH)/customer customer.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --grpc-gateway_out=logtostderr=true:$(API_OUT_PATH)/customer customer.proto
	@protoc -I=$(API_IN_PATH) -I=third_party --swagger_out=logtostderr=true:$(API_DOC_PATH) customer.proto
	
run_bank:
	@cd cmd/bank && go run ./main.go --config-file=./config.yml