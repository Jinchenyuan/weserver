.PHONY: swagger
swagger:
	swag init -g api/s3/ginhandler/s3.go --parseDependency --parseInternal -o api/s3/docs
	swag init -g api/account/ginhandler/account.go --parseDependency --parseInternal -o api/account/docs