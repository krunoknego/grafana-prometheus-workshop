# Prometheus and Grafana Workshop

## Create a Grafana Dashboard

### Connect Prometheus and Grafana
1. http://localhost:3000 - log in admin:admin
2. Connections -> Data Sources -> Choose Prometheus
3. Set `Prometheus server URL` to http://prometheus:9090
4. Press `Save & test`


### Create a Text Panel
1. Dashboard -> New -> New Dashboard 
2. Add visualization -> Select data source `prometheus`
3. Change visualization type to `Text` (upper right corner)
4. Add to the content:

```
# GO Application

hello ${__user.login}!

- [Documentation](http://example.com)
- [Incident Report](http://example.com)
```

The `__user.login` is a variable that will be replaced with the username of the user that is currently logged in.

5. Give the panel a title `Welcome`
6. Press `Save dashboard`
7. Press `Back to dashboard`

### Create a Bar Gauge Panel
1. Add -> Visualization (choose `Bar Gauge`)
2. For the query use `sum by (path)(rate(goapp_http_requests_total[5m]))`
3. Give the panel a title `Requests per path`
4. Set the unit to `requests/sec`, you can find this in the `Standard options` section, Unit -> Throughput -> requests/sec
5. Press `Save dashboard`
6. Press `Back to dashboard`

### Create a CPU Gauge Panel
1. Add -> Visualization (choose `Gauge`)
2. For the query use `100 - avg by(instance)(rate(node_cpu_seconds_total{mode="idle"}[1m])) * 100`
3. Give the panel a title `CPU Usage`
4. Set the unit to `%`, you can find this in the `Standard options` section, Unit -> Misc -> Percent (0-100)
5. Press `Save dashboard`
6. Press `Back to dashboard`

### Create a Memory Gauge Panel
1. Add -> Visualization (choose `Gauge`)
2. For the query use `100 * (1 - (node_memory_MemAvailable_bytes / node_memory_MemTotal_bytes))`
3. Give the panel a title `Memory Usage`
4. Set the unit to `%`, you can find this in the `Standard options` section, Unit -> Misc -> Percent (0-100)
5. Press `Save dashboard`
6. Press `Back to dashboard`

### HTTP Error Rate Panel (Stat Panel)
1. Add -> Visualization (choose `Stat`)
2. For the query use `sum(rate(goapp_http_requests_total{status=~"5.*"}[5m])) / sum(rate(goapp_http_requests_total[5m])) * 100`
3. Give the panel a title `Error Rate`
4. Give the panel a description `Monitor the percentage of HTTP errors (5xx) compared to total requests.`
5. Set the unit to `%`, you can find this in the `Standard options` section, Unit -> Misc -> Percent (0-100)
6. Press `Save dashboard`
7. Press `Back to dashboard`

### Requests Growth Over Time (Time Series Panel)
1. Add -> Visualization (choose `Time Series`)
2. For the query use `sum(rate(goapp_http_requests_total[1m]))`
3. Give the panel a title `Total HTTP Requests Over Time`
4. Give the panel a description `Monitor the growth of HTTP requests over time.`
5. Set the unit to `requests/sec`, you can find this in the `Standard options` section, Unit -> Throughput -> requests/sec
6. Press `Save dashboard`
7. Press `Back to dashboard`

### Generate the traffic
1. In terminal run `make generate-traffic`
2. Open the Grafana dashboard that you created
3. In the upper right corner set the time range to `Last 15 minutes` and Refresh every `5s`
4. You should see the traffic on the dashboard
5. If you want to stop traffic generation, press `Ctrl+C` in the terminal

## Import a Grafana Dashboard

### Import NodeExporter Full
1. Dashboard -> New -> Import
2. Enter the dashboard ID `1860` and press `Load`
3. Choose Prometheus as the data source
4. Press `Import`

### Import Blackbox Exporter
1. Dashboard -> New -> Import
2. Enter the dashboard ID `7587` and press `Load`
3. Choose Prometheus as the data source
4. Press `Import`

### Generate the traffic
1. In terminal run `make generate-traffic`
2. Open the Grafana dashboards that you created
3. You should see the traffic on the dashboard
4. If you want to stop traffic generation, press `Ctrl+C` in the terminal
