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

## Minikube (Kubernetes)
- spin up: `minikube start`
- spin down: `minikube stop`

#### Setting secrets
- set env vars on the cluster: `kubectl create -f mysql-secret.yaml`
- check secrets were created: `kubectl get secrets`
- check secret elements: `kubectl describe secrets mysql-secret`

#### Applying yaml commands
- Order: pv, pvc, deployment, service
```
kubectl apply -f mysql-db-pv.yaml
kubectl apply -f mysql-db-pvc.yaml
kubectl apply -f mysql-db-deployment.yaml
kubectl apply -f mysql-db-service.yaml
```
- view pod after service is spun up: `kubectl get pods`
- check logs:
```
kubectl describe pod <pod_name>
kubectl logs <pod_name>
```

## Pushing API to dockerhub
- Build image for kubernetes: `docker build -t go-blog-kubernetes .`
- Tag image to repo on dockerhub: `docker tag <image-name> <dockerhub-username>/<repository-name>:<tag-name>`
  - ie. `docker tag go-blog-kubernetes dmdinh/go-blog:X.X.X`

## Deploying API to Kubernetes
- Create and apply kubernetes from the yaml files:
```
kubectl create -f mysql-secret.yaml
kubectl apply -f app-mysql-deployment.yaml
kubectl apply -f app-mysql-service.yaml
```

- Check pod status: `kubectl get pods`
- Get services in cluster: `kubectl get services`
- Get url exposed: `minikube service <service-name> --url` (`minikube service go-blog-mysql --url`)
