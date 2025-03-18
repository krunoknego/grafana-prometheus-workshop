# Prometheus and Grafana Workshop

## Configure a Prometheus Agent

Agent mode is when Prometheus runs in a reduced mode, where it only scrapes metrics and sends them to the Prometheus server using remote_write protocol.

1. In your `docker-compose.yaml`, below `service` add prometheus server

``` yaml
  prometheus:
    image: prom/prometheus:latest
    container_name: prometheus
    volumes:
      - ./prometheus_server.yml:/etc/prometheus/prometheus.yml:ro
    ports:
      - "9090:9090"
    command:
      - "--config.file=/etc/prometheus/prometheus.yml"
      - "--web.enable-remote-write-receiver"
```

The remote_write needs to be enabled on the Prometheus server. This allows Prometheus to receive metrics from the Prometheus agent.

2. In `docker-compose.yaml` add the Prometheus agent

``` yaml
  prometheus-agent:
    image: prom/prometheus:latest
    container_name: prometheus-agent
    volumes:
      - ./prometheus_agent.yml:/etc/prometheus/prometheus.yml:ro
    depends_on:
      - node-exporter
    command:
      - "--enable-feature=agent"
      - "--config.file=/etc/prometheus/prometheus.yml"
```

Notice the image of the Prometheus agent is the same as the Prometheus server. The only difference is the command that is passed to the agent.

For the agent, we need to enable the agent feature `--enable-feature=agent`.

3. In `docker-compose.yaml` add the node-exporter

``` yaml
  node-exporter:
    image: prom/node-exporter:latest
    container_name: node-exporter
    ports:
      - "9100:9100"
```

4. Create a `prometheus_agent.yml` file

``` yaml
global:
  scrape_interval: 15s

scrape_configs:
  - job_name: 'node'
    static_configs:
      - targets: ['node-exporter:9100']

remote_write:
  - url: "http://prometheus:9090/api/v1/write"
```

5. Create a `prometheus_server.yml` file

``` yaml
global:
  scrape_interval: 2s

rule_files: []

scrape_configs:
  - job_name: 'prometheus'
    static_configs:
      - targets: ['localhost:9090']
```

Here we have the Prometheus server configuration. The server scrapes itself every 2 seconds.

6. Start the containers `make start`

7. Open the Prometheus server UI at `http://localhost:9090` and check if the Prometheus agent is sending metrics. 
   You can check this by going to the `Status` -> `Target health` tab.
   
   How many targets are up?
   
8. Write a PromQL query to get the number of targets that are up. `up`
   
   How many targets are up?
   
   Why is it different?

## Configure Telegraf with Prometheus

Now we will configure Telegraf to send metrics to Prometheus. 
This architecture will consist out of two Telegraf clients that 
gather the data and one Telegraf Aggregator that aggregates the data and sends it to Prometheus.

First let's create the Telegraf clients.

1. In your `docker-compose.yaml` add the Telegraf client 1

``` yaml
  telegraf-client1:
    image: telegraf:latest
    container_name: telegraf-client1
    volumes:
      - ./telegraf_client.conf:/etc/telegraf/telegraf.conf:ro
    environment:
      - HOSTNAME=telegraf-client1
```

2. In your `docker-compose.yaml` add the Telegraf client 2

``` yaml
  telegraf-client2:
    image: telegraf:latest
    container_name: telegraf-client2
    volumes:
      - ./telegraf_client.conf:/etc/telegraf/telegraf.conf:ro
    environment:
      - HOSTNAME=telegraf-client2
```

3. Create a `telegraf_client.conf` file

``` toml
[agent]
  interval = "10s"
  round_interval = true

[[inputs.cpu]]
  percpu = true
  totalcpu = true
  collect_cpu_time = false
  report_active = false

[[inputs.mem]]

[[inputs.disk]]
  ignore_fs = ["tmpfs", "devtmpfs", "overlay"]

[[outputs.socket_writer]]
  address = "tcp://telegraf-aggregator:8094"
  data_format = "influx"
```

This configuration will gather CPU, memory, and disk metrics and send them to the Telegraf Aggregator.
Both Telegraf clients will use this configuration. The only difference is the HOSTNAME.

4. Now let's create the Telegraf Aggregator. In your `docker-compose.yaml` add the Telegraf Aggregator

``` yaml
  telegraf-aggregator:
    image: telegraf:latest
    container_name: telegraf-aggregator
    volumes:
      - ./telegraf_aggregator.conf:/etc/telegraf/telegraf.conf:ro
    ports:
      - "9273:9273"
    depends_on:
      - telegraf-client1
      - telegraf-client2
```

5. Create a `telegraf_aggregator.conf` file

``` toml
[agent]
  interval = "10s"
  round_interval = true

[[inputs.socket_listener]]
  service_address = "tcp://:8094"
  data_format = "influx"

[[outputs.prometheus_client]]
  listen = ":9273"
  metric_version = 2
```

The Telegraf Aggregator listens on port 8094 for incoming metrics from the Telegraf clients.
The two clients and the aggregator usethe same data format, Influx.

The aggregator gives the metrics to Prometheus on port 9273.

6. Now we need to update the Prometheus server configuration to scrape the metrics from the Telegraf Aggregator.
   Update the `prometheus_server.yml` file
   
   Under `scrape_configs` add another job (don't delete the previous ones):
   ```yaml
  - job_name: 'telegraf'
    static_configs:
      - targets: ['telegraf-aggregator:9273']
   ```
   
7. Now that we have everything set up, start the containers `make down start`

8. Open the Prometheus server UI at `http://localhost:9090` and check if the Telegraf Aggregator is sending metrics. 
   You can check this by going to the `Status` -> `Target health` tab.
   
   How many targets are up?
   
9. Now open Telegraf Aggregator metrics at `http://localhost:9273/metrics`, you should see the metrics that are being sent to Prometheus.
   Choose a metric and write a PromQL query to get the value of that metric.

## Configure Blackbox Exporter with Prometheus

Now we will configure the Blackbox Exporter to monitor the availability of the Go applications.
In you day5 directory there should be `go-app1` and `go-app2` directories.

Inspect the code of the applications.

1. In your `docker-compose.yaml` add the Go application 1

``` yaml
  go-app1:
    build: ./go-app1
    container_name: go-app1
    ports:
      - "8081:8080"
```

2. Now do the same and add the Go application 2

``` yaml
  go-app2:
    build: ./go-app2
    container_name: go-app2
    ports:
      - "8082:8080"
```

The Go applications will be available at `http://localhost:8081` and `http://localhost:8082`.

3. Now let's add the Blackbox Exporter to the `docker-compose.yaml`

``` yaml
  blackbox_exporter:
    image: prom/blackbox-exporter:latest
    container_name: blackbox_exporter
    ports:
      - "9115:9115"
    volumes:
      - ./blackbox.yml:/etc/blackbox_exporter/config.yml:ro
```

4. Now we need to create configuration for the Blackbox Exporter. Create a `blackbox.yml` file and add the following configuration

``` yaml
modules:
  http_2xx:
    prober: http
    timeout: 5s
    http:
      valid_http_versions: ["HTTP/1.1", "HTTP/2"]
      method: GET
```

We are done with the configuration of the Blackbox Exporter.
Notice how there is no mention of the Go applications in the configuration.

Why is that?

5. Now we need to update the Prometheus server configuration to scrape the metrics from the Blackbox Exporter.
   Update the `prometheus_server.yml` file
   
   Under `scrape_configs` add another job (don't delete the previous ones):
   ```yaml
  - job_name: 'blackbox'
    metrics_path: /probe
    params:
      module: [http_2xx]
    static_configs:
      - targets: ["go-app1:8080", "go-app2:8080"]
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - source_labels: [__param_target]
        target_label: instance
      - target_label: __address__
        replacement: blackbox_exporter:9115
   ```
   
Take a look at the `relabel_configs` section. What is happening here?

Can you explain the `source_labels`, `target_label`, and `replacement`?

6. Now that we have everything set up, start the containers `make down start`

7. Open the Prometheus server UI at `http://localhost:9090` and check if the Blackbox Exporter is sending metrics. 
   You can check this by going to the `Status` -> `Target health` tab.
   
   How many targets are up?
   
8. Write a PromQL query to get the duration of the requests in seconds. 
   How would you go about finding out the metrics that are being sent by the Blackbox Exporter?
   
9. Notice how there is only one exporter for both Go applications. 
   How would you go about getting the metrics for another Go application (go-app3)?

## Configure JSON Exporter with Prometheus

Now we will configure the JSON Exporter to extract metrics from a Go application.
The Go application will return requests in JSON format.
The JSON Exporter will extract the metrics from the JSON response and convert them to Prometheus format.

1. In your `docker-compose.yaml` add the Go JSON application

``` yaml
  go-json:
    build: ./go-json
    container_name: go-json
    ports:
      - "8083:8080"
```

2. There should be a `go-json` directory in your day5 directory. Inspect the code of the application.

3. Now let's add the JSON Exporter to the `docker-compose.yaml`

``` yaml
  json_exporter:
    image: prometheuscommunity/json-exporter:latest
    container_name: json_exporter
    command:
      - "--config.file=/etc/json_exporter/config.yml"
    ports:
      - "7979:7979"
    volumes:
      - ./json_exporter.yml:/etc/json_exporter/config.yml:ro
```

4. Now we need to create configuration for the JSON Exporter. Create a `json_exporter.yml` file and add the following configuration

``` yaml
---
modules:
  default:
    metrics:
    - name: random_number_value
      path: '{ .random_number }'
      help: Example of a top-level global value scrape in the json
      labels:
        env: test
        timestamp: '{ .timestamp }'
        id: '{ .unique_id }'
```

This configuration will extract the `random_number`, `timestamp`, and `unique_id` from the JSON response.

5. Now we need to update the Prometheus server configuration to scrape the metrics from the JSON Exporter.
   Update the `prometheus_server.yml` file
   
   Under `scrape_configs` add another job (don't delete the previous ones):
   ```yaml
  - job_name: 'json_exporter'
    metrics_path: /probe
    params:
      module: [default]
    static_configs:
      - targets: [ "http://go-json:8080" ]
    relabel_configs:
      - source_labels: [__address__]
        target_label: __param_target
      - target_label: __address__
        replacement: json_exporter:7979
   ```

Take a look at the `relabel_configs` section. What is happening here?

Can you explain the `source_labels`, `target_label`, and `replacement`?

6. Now that we have everything set up, start the containers `make down start`

7. Open the Prometheus server UI at `http://localhost:9090` and check if the JSON Exporter is sending metrics. 
   You can check this by going to the `Status` -> `Target health` tab.
   
   How many targets are up?
   
8. What metrics are being sent by the JSON Exporter? Write a PromQL query to get the value of the metric.
