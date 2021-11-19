gen-cal:
	protoc --proto_path=calculatorpb --go_out=calculatorpb --go_opt=paths=source_relative --go-grpc_out=calculatorpb --go-grpc_opt=paths=source_relative calculator.proto