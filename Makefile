.PHONY: swagger
swagger:
	swag init -g api/s3/main.go --parseDependency --parseInternal -o api/s3/docs