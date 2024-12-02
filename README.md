# theme-park-data
This repo will contain:
- Database schema for app `theme-park-data`
- Test dataset
- Dockerfile for deploying its database (using PSQL RDBS)
- App written in go for wrapping up the database (Rest API)
- Dockerfile for deploying the Rest API app
- Docker-compose file for deploying all infra


## How to run it using docker-compose
- Make sure you have docker installed
- Checkout the repo locally
- Build the rest-api container and run the db. 
```bash
docker-compose up --build -d
```
The db will be exposed on `locahost:5432` and restapi will be available on `localhost:8080`.

Once done, kill it all using:
```bash
docker-compose down -v
```

Current implemented public (`/pub` means without token and `/priv` needs a token)
- http://localhost:8080/pub/login
- http://localhost:8080/pub/getCategories
- http://localhost:8080/pub/register
- http://localhost:8080/priv/users/{id}

To check what JSON info to pass, refer to `cheatsheet.txt`.

