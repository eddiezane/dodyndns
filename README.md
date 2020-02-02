# dodyndns

Updates a DigitalOcean domain record with your public IP

# Usage

```
dodyndns --domain example.com --record home --token $DIGITALOCEAN_TOKEN
```

| Flag | Description |
| --- | --- |
| domain | The name of the domain on your DigitalOcean account |
| record | The record name to update |
| token | Your DigitalOcean API token |

I run this from a K3s cluster of Raspberry Pi 4's.

To build it into an image I use Docker Buildx.

```
docker buildx create --name dodyndns
docker buildx inspect --bootstrap
docker buildx build --platform linux/arm/v7 -t eddiezane/dodyndns:latest --push .
```

See [job.yml](job.yml) for an example manifest.

Create a Kubernetes secret with:

```
kubectl create secret generic dodyndns --from-literal=token=YOUR_TOKEN
```

# License

MIT
