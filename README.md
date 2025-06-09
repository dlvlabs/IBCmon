# IBCmon

`IBCmon` is the monitoring program whether working well all of the connected IBC from a specific chain(base chain).

## Features

- Monitoring

    - **IBC TAO**: Automatically search and store all of the well functioning IBC TAO information related with the base chain automatically and vice versa(this mean not only base chain but also couterparties)

    - **Client Health**: Monitoring whether ibc clients are update well and there's a risk for expired

    - **IBC Packet**: Monitoring IBC tx is sent, received well through specific IBC TAO

- JSON API

    - `/ibc-info`: List of well functioning IBC TAO information

    - `/client-health`: List of client health

    - `/ibc-packet`: List of ibc channels and packets information

- Prometheus 

    - `/metrics`: Metrics for IBC TAO, client health, and ibc packets

## Quick Guide

1. **Build**
```bash
make docker-build
```

2. **Fill `config.toml`**
```bash
cp config.toml.example config.toml
vim config.toml
```

3. **Change `docker-compose.yml`**
```bash
vim docker-compose.yml
```

4. **Run**
```bash
docker-compose up -d
```
