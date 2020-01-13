# go-blog

### Tech
- Go
- GORM
- JWT
- Postgres
- Mysql
- Gorilla Mux (For HTTP routing and URL matcher)


### Running Test Suite
- `cd test`
- `go test -v ./...`

### Running Tests for a Module
- `cd test/$whatever_test_dir_your_test_is_in`
- `go test -v`

## Docker
#### Docker Commands
- Run: `docker-compose up`
- Run detached: `docker-compose up -d`
- Rebuild container: `docker-compose up --build`
- Bring down container: `docker-compose down`
- Stop running process and clear volumes: `docker-compose down --remove-orphans --volumes`
- Get IP address of docker container: `docker inspect <container_id> | grep IPAddress` -> hostname
- Remove dangling images: `docker system prune`
- Remove all unused images: `docker system prune -a`
- Remove all unused images and volumes: `docker system prune -a --volumes`

#### Running tests
- from root dir, run `docker-compose -f docker-compose.test.yml up --build --abort-on-container-exit`
