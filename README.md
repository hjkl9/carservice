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

## Swagger CLI
Install go cli: ```go install github.com/zeromicro/goctl-swagger@latest```
Run generation 1: ```goctl api plugin -plugin goctl-swagger="swagger -filename carservice.json" -api api/carservice.api -dir .```
Run generation 2: ```goctl api plugin -plugin goctl-swagger="swagger -filename carservice.json -host 127.0.0.1 -basepath /api" -api api/carservice.api -dir .```
**On server** just using docker(No test): ```docker run --rm -p 8083:8080 -e SWAGGER_JSON=carservice.json -v $PWD:/ swaggerapi/swagger-ui```