# Prometheus and Grafana Workshop

## Excercise 1

Task: Star the docker containers

1. Execute `make start` to start the docker containers.

Task: Connect prometheus to grafana

1. Log in to Grafana (default username: admin, password: admin).
2. Go to Configuration > Data Sources.
3. Click "Add data source" and select Prometheus.
4. Set the URL to http://prometheus:9090.
5. Click "Save & Test" to verify the connection.

Task: Write a PromQL query in Grafana to visualze how fast the CPU time is accumulating

1. Go to Dashboards
2. Create Dashboard > Add visualization
3. Select Prometheus
4. Write the PromQL query below

```sh
rate(process_cpu_seconds_total[1m])
```

## Exercise 2

Task: install `node_exporter` and connect it to Prometheus

1. Go to the `docker-compose.yaml` 
2. Uncomment the `node_exporter` service and execute `make down start`
3. Go to `http://localhost:9100/metrics` to verify that the metrics are being exposed (check TYPE and HELP)
4. Go to `http://localhost:9090/targets` to check which targets are being scraped
5. Update prometheus configuration `prometheus.yml` to scrape the `node_exporter` service and execute `make down start`

``` yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
  - job_name: 'node_exporter'
    static_configs:
      - targets: ['node-exporter:9100']
```
6. Go to `http://localhost:9090/targets` to check which targets are being scraped
7. Write a PromQL query to visualize the CPU usage of the node_exporter

```sh
100 - (avg by (instance) (irate(node_cpu_seconds_total{mode="idle"}[5m]))) * 100
```
8. Open Graph and inspect how the graph changes

## Exercise 3

Task: Install Grafana Alloy and connect it to Prometheus

1. Go to the `docker-compose.yaml`
2. Uncomment the `grafana-alloy` service and prometheus `command:` and execute `make down start`
3. Go to `http://localhost:12345/` to check components of Alloy (should be empty)
4. Copy the text below and replace the contents of `config.alloy` and execute `make down start`

``` sh
prometheus.exporter.unix "local_system" { }

prometheus.scrape "scrape_metrics" {
  targets         = prometheus.exporter.unix.local_system.targets
  forward_to      = [prometheus.relabel.filter_metrics.receiver]
  scrape_interval = "10s"
}

prometheus.relabel "filter_metrics" {
  rule {
    action        = "drop"
    source_labels = ["env"]
    regex         = "dev"
  }

  forward_to = [prometheus.remote_write.metrics_service.receiver]
}

prometheus.remote_write "metrics_service" {
    endpoint {
        url = "http://prometheus:9090/api/v1/write"

        basic_auth {
          username = "admin"
          password = "admin"
        }
    }
}
```

5. Go to `http://localhost:12345/` to check components of Alloy (should NOT be empty)
6. Go to `http://localhost:12345/` and open Graph and check the relationship between the components

## Exercise 4

Task: Add node_exporter to Grafana Alloy

1. Replace the contents of `prometheus.yml` with the contents from `prometheus.yml.old`
2. Copy the text below and replace the contents of `config.alloy` and execute `make down start`

``` sh
prometheus.exporter.unix "local_system" { }

prometheus.scrape "scrape_metrics" {
  targets         = prometheus.exporter.unix.local_system.targets
  forward_to      = [prometheus.relabel.filter_metrics.receiver]
  scrape_interval = "10s"
}

prometheus.relabel "filter_metrics" {
  rule {
    action        = "drop"
    source_labels = ["env"]
    regex         = "dev"
  }

  forward_to = [prometheus.remote_write.metrics_service.receiver]
}

prometheus.remote_write "metrics_service" {
    endpoint {
        url = "http://prometheus:9090/api/v1/write"

        basic_auth {
          username = "admin"
          password = "admin"
        }
    }
}

prometheus.scrape "node_exporter" {
  // Collect metrics from the default listen address.
  targets = [{
    __address__ = "node-exporter:9100",
  }]

  forward_to = [prometheus.relabel.filter_metrics.receiver]
}
```
3. Go to `http://localhost:9090/targets` to check which targets are being scraped
4. Write `up` to check if the node_exporter is being scraped
