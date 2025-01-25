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

dev-auth-start:
	cd services/auth && air

dev-category-start:
	cd services/category && air

dev-product-start:
	cd services/product && air init && air

dev-file-start:
	cd services/files  && air init && air
	
dev-gateway-start:
	cd services/gateway   && air
	
