# Deploy

```bash
docker build -t operator:latest -f Dockerfile ..

kind load docker-image operator:latest

kubectl apply -f manifests/
```