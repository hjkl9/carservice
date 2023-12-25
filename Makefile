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

SWAGGER_CONTAINER_NAME = "swagger_ui"
swaggerdown:
	docker stop $(SWAGGER_CONTAINER_NAME)
swaggerup:
	docker run --rm -d -p 8083:8080 --name $(SWAGGER_CONTAINER_NAME) -e "SWAGGER_JSON=/carservice.json" -v $(shell pwd)/carservice.json:/carservice.json  swaggerapi/swagger-ui

# deployment related.
APP_NAME=carservice
# existing mirrors and containers need to be deleted before calling.
up1:
	git pull
	docker build . -t ${APP_NAME}
	docker run -d --name ${APP_NAME} \
		--link carservice_mysql:mysql \
		--link carservice_redis:redis \
		--net deploy_default -p 8888:8888 \
		${APP_NAME}

down1:
	docker stop ${APP_NAME}
	docker rm ${APP_NAME}
	docker rmi ${APP_NAME}

# todo: should add link for redis.
# todo: specify the link of mysql and network.
PROJECT_LOGS_LINUX=/root/projects/master/carservice/logs
up2:
	git pull
	docker build . -t ${APP_NAME}
	docker run --rm -d \
		--name ${APP_NAME} \
		--link carservice_mysql:mysql \
		--net deploy_carservice_admin_network \
		-v ${PROJECT_LOGS_LINUX}:/app/logs \
		-p 8888:8888 \
		${APP_NAME}
down2:
	docker stop ${APP_NAME}
	docker rm ${APP_NAME}
	docker rmi ${APP_NAME}

