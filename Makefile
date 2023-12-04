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

# deployment related.
APP_NAME=carservice
# existing mirrors and containers need to be deleted before calling.
prod:
	git pull
	docker build . -t ${APP_NAME}
	docker run -d --name ${APP_NAME} -p 8888:8888 ${APP_NAME}