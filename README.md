# go-blog

### Tech
- Go
- GORM
- JWT
- Postgres
- Mysql
- Gorilla Mux (For HTTP routing and URL matcher)

### Running Test Suite Locally
- Make sure API and DB is running locally
- `cd test`
- `go test -v ./...`

### Running Tests for a Module
- Make sure API and DB is running locally
- `cd test/$whatever_test_dir_your_test_is_in`
- `go test -v`

## Docker
#### Docker Commands
- From root dir of app
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

## Swagger Docs
- When the environment is running, go to `localhost:8080/swagger/index.html`
  - `localhost:8080` or `localhost:8080/swagger` in your browser will also redirect to the swagger documentation

## Minikube (Kubernetes) Deployment
- spin up: `minikube start`
- spin down: `minikube stop`

#### Setting secrets
- set env vars on the cluster: `kubectl create -f mysql-secret.yaml`
- check secrets were created: `kubectl get secrets`
- check secret elements: `kubectl describe secrets mysql-secret`

#### Deploying DB to Kubernetes
- Order: pv, pvc, deployment, service
```
kubectl apply -f mysql-db-pv.yaml
kubectl apply -f mysql-db-pvc.yaml
kubectl apply -f mysql-db-deployment.yaml
kubectl apply -f mysql-db-service.yaml
```
- view pod after service is spun up: `kubectl get pods`
- check logs if there's an error:
```
kubectl describe pod <pod_name>
kubectl logs <pod_name>
kubectl logs <pod-name> -c <container-name>
```
- `container-name` comes from `describe` command

## Pushing API to dockerhub
- Build image for kubernetes: `docker build -t go-blog-kubernetes .`
- Tag image to repo on dockerhub: `docker tag <image-name> <dockerhub-username>/<repository-name>:<tag-name>`
  - ie. `docker tag go-blog-kubernetes dmdinh/go-blog:X.X.X`

#### Deploying API to Kubernetes
- Create and apply kubernetes from the yaml files:
```
kubectl apply -f api-mysql-deployment.yaml
kubectl apply -f api-mysql-service.yaml
```

- Check pod status: `kubectl get pods`
- Get services in cluster: `kubectl get services`
- Get url exposed: `minikube service <service-name> --url` (`minikube service go-blog-mysql --url`)
