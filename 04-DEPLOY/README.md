# Deploy

```bash
cd 04-DEPLOY/

docker build -t operator:latest -f Dockerfile ..

kind load docker-image operator:latest

kubectl apply -f manifest/
```