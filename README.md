# 🐳 Docker
Add this line in daemon.json <br />
- **LINUX:** /etc/docker/daemon.json

```json
"metrics-addr": "127.0.0.1:9323"
```

# 🆙 Start
```bash
docker-compose up -d --build
docker-compose exec app go run main.go
```
# 📦 Services
- [API GoLang](http://localhost:8181)
- [Prometheus](http://localhost:9090)
- [Grafana](http://localhost:3000)

# Grafana
- User: admin
- Password: admin
#### Add DataSource
- url: http://prometheus:9090