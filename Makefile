.PHONY: pb_register

pb_register:
	protoc register/registerpb/register.proto --go_out=plugins=grpc:register/registerpb/
	#
	#protoc -I register/registerpb/ register/registerpb/register.proto --go_out=plugins=grpc:.