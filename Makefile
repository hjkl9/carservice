# generate or update.
generate:
	goctl api go -api api/carservice.api -dir .

# swagger related.
# todo

# migrate related.
LOCALADDR = "127.0.0.1:3306"
LOCALUSR = "root"
LOCALPWR = ""
migratelocal:
	migrate -database 'mysql://$(LOCALUSR):$(LOCALPWD)@$(LOCALADDR)/carservice?sslmode=disable' -path ./migrations up 1

buildapp:
	docker build . -t ${APP_NAME}

swagger:
	docker run --rm -d -p 8083:8080 --name swagger_ui -e "SWAGGER_JSON=/carservice.json" -v $(shell pwd)/carservice.json:/carservice.json  swaggerapi/swagger-ui

# deployment related.
APP_NAME=carservice
# existing mirrors and containers need to be deleted before calling.
dockerup:
	git pull
	docker build . -t ${APP_NAME}
	docker run -d --name ${APP_NAME} \
		--link carservice_mysql:mysql \
		--link carservice_redis:redis \
		--net deploy_default -p 8888:8888 \
		${APP_NAME}

dockerdown:
	docker stop ${APP_NAME}
	docker rm ${APP_NAME}
	docker rmi ${APP_NAME}