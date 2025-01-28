gen-auth:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/auth.proto

gen-category:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/category.proto

gen-product:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/product.proto

gen-file:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/files.proto

gen-stock:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/stock.proto

gen-order:
	@protoc --experimental_allow_proto3_optional \
	--go_out=. --go_opt=paths=source_relative \
	--go-grpc_out=. --go-grpc_opt=paths=source_relative \
	grpc/pb/order.proto

down :
	docker compose -f docker-compose-dev.yml down

build :
	docker compose -f docker-compose-dev.yml build

up :
	docker compose -f docker-compose-dev.yml up


logs-auth:
	docker logs golang-backend-auth-1

logs-files:
	docker logs golang-backend-files-1

logs-gateway:
	docker logs golang-backend-gateway-1

logs-category:
	docker logs golang-backend-category-1

logs-product:
	docker logs golang-backend-product-1

logs-stock:
	docker logs golang-backend-stock-1

logs-order:
	docker logs golang-backend-order-1

start-gateway:
	docker compose -f docker-compose-dev.yml up -d gateway
	
