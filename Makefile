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