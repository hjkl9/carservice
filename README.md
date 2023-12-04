# Carservice API

## Database Migration CLI
### For windows
1. Open PowerShell(x64).
2. Set installation directory: ```$env:SCOOP='PATH'``` and ```[Environment]::SetEnvironmentVariable('SCOOP', $env:SCOOP, 'User')```
3. Allow PowerShell to execute the local scripts: ```Set-ExecutionPolicy RemoteSigned -scope CurrentUser```
4. Install scoop: ```iwr -useb get.scoop.sh | iex```
5. Install migrate ```scoop install migrate```

### For Linux.
1. ```curl -L https://packagecloud.io/golang-migrate/migrate/gpgkey | apt-key add -```
2. ```$ echo "deb https://packagecloud.io/golang-migrate/migrate/ubuntu/ $(lsb_release -sc) main" > /etc/apt/sources.list.d/migrate.list```
3. ```$ apt-get update```
4. ```$ apt-get install -y migrate```

## Create migration file
```migrate create -ext sql -dir ./migrations -seq create_xxx(s)_table```
## Execute migration
```migrate -database 'mysql://user:password@127.0.0.1:3306/carservice?sslmode=disable' -path ./migrations up 1```

## Swagger CLI
1. Install go cli: ```go install github.com/zeromicro/goctl-swagger@latest```
2. Run generation 1: ```goctl api plugin -plugin goctl-swagger="swagger -filename carservice.json" -api api/carservice.api -dir .```
3. Run generation 2: ```goctl api plugin -plugin goctl-swagger="swagger -filename carservice.json -host 127.0.0.1 -basepath /api" -api api/carservice.api -dir .```
4. **On server** just using docker(No test): ```docker run --rm -p 8083:8080 -e "SWAGGER_JSON=/carservice.json" -v $PWD/carservice.json:/carservice.json  swaggerapi/swagger-ui```