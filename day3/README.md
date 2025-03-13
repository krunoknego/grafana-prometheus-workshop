# Prometheus and Grafana Workshop

## Create a Grafana User and assign it a role

1. Execute the following command `make start`
2. Log into grafana http://localhost:3000
3. Go to Administration -> Users and Access -> Users
4. Press `New User` button and enter the credentials
5. Edit the user's Role to `Editor``

## Grafana OAuth with Github

### 1. Register a New Github OAuth Application

1. Open your Github account
2. Profile Picture -> Settings -> Developer settings -> OAuth Apps
3. Click `New OAuth App`
   - Application name: `Grafana OAuth`
   - Homepage URL: `http://localhost:3000`
   - Authorization callback URL: `http://localhost:3000/login/github`
4. Click `Register application`
5. Note the generated `Client ID` and `Client Secret`

### 2. Configure Grafana to Use Github OAuth

Modify your `docker-compose.yml` file under `grafana` service to include these environment variables:

``` yaml
grafana:
  image: grafana/grafana
  ports:
    - 3000:3000
  environment:
    - GF_SECURITY_ADMIN_PASSWORD=admin
    - GF_AUTH_GITHUB_ENABLED=true
    - GF_AUTH_GITHUB_CLIENT_ID=<Your GitHub Client ID>
    - GF_AUTH_GITHUB_CLIENT_SECRET=<Your GitHub Client Secret>
    - GF_AUTH_GITHUB_SCOPES=user:email,read:org
    - GF_AUTH_GITHUB_AUTH_URL=https://github.com/login/oauth/authorize
    - GF_AUTH_GITHUB_TOKEN_URL=https://github.com/login/oauth/access_token
    - GF_AUTH_GITHUB_API_URL=https://api.github.com/user
  networks:
    - monitoring
```

Replace <Your GitHub Client ID> and <Your GitHub Client Secret> with the values obtained in step 1 (Register a New Github OAuth Application).

### 3. Restart Grafana to Apply Changes

1. Execute `make down start`
2. Open `http://localhost:3000`
3. Click `Sign in with Github`
4. What Role do you have? What can you do with this Role?


## Configure prometheus rule and make sure that it is firing

### 1. Create Alerting Rule

Add a rule to your `prometheus.yml` configuration file

``` yaml
rule_files:
  - 'alerts.yml'
```

Modify file named `alerts.yml`:

``` yaml
groups:
  - name: example_alert
    rules:
      - alert: InstanceDown
        expr: up == 0
        for: 10s
        labels:
          severity: critical
        annotations:
          summary: "Instance {{ $labels.instance }} down"
          description: "The instance {{ $labels.instance }} has been down for more than 10 seconds."
```

Update you `docker-compose.yml`

``` yaml
services:
  prometheus:
    image: prom/prometheus
    ports:
      - 9090:9090
    volumes:
      - ./prometheus.yml:/etc/prometheus/prometheus.yml
      - ./alerts.yml:/etc/prometheus/alerts.yml # this lined needs to be added
    networks:
      - monitoring
```

Execute `make down start`

### 2. Verify the Alert state

1. Open `http://localhost:9090`
2. Go to the `Alerts` tab. What do you see? In what state is the alarm?
3. In your terminal execute the following command `docker compose kill node-exporter`
4. Open `http://localhost:9090`
5. Go to the `Alerts` tab. What do you see? In what state is the alarm? You have to wait 10 seconds, why is that?
6. Execute `make start` and once again check the state of the Alarm.


## Configure an alert using alertmanager

### 1. Get a Webhook.site Unique URL

1. Go to https://webhook.site 
2. Copy the unique URL

### 2. Update `alertmanager.yml`

Modify your `alertmanager.yml` file to send alerts to Webhook.site

``` yaml
route:
  receiver: "webhook-site"

receivers:
  - name: "webhook-site"
    webhook_configs:
      - url: "<Your Webhook.site Unique URL>" # Replace this line
        send_resolved: true
```

### 3. Connect Prometheus to Alertmanager

Modify `prometheus.yml` to include the following lines.

``` yaml
alerting:
  alertmanagers:
    - scheme: http
      static_configs:
        - targets: [ 'alertmanager:9093' ]
```

Restart the services `make down start`.

### 4. Verify Alert is Routed to Webhook.site

1. First verify that prometheus alertmanager component is running successfully. Go to http://localhost:9093
2. Execute the following command to simulate an alert `docker compose kill node-exporter`
3. Open Alertmanager (http://localhost:90903) and check if the alert appears
   Remember that the alert in Prometheus (http://localhost:9090) needs to be in firing state for it to appear in the Alertmanager.
4. Go to https://webhook.site and verify that you have received the notification
5. Execute `make start` to stop the error

## Create an alert in Grafana

### 1. Configure the Data Source

1. Go to `http://localhost:3000` and log in (you need to have Admin role)
2. Click `Data Sources` on the left side
3. Click `Add data source` and choose Prometheus
4. For connection set `http://prometheus:9090` 
5. Click `Save & test`

### 2. Configure the Contact Points

1. Click `Contact points` on the left side
2. Click `Create contact point`
3. Configure the contact point:
   Name: `TestContactPoint`
   Integration: `Webhook`
   URL: `<Put here your webhook.site unique>` (go to https://webhook.site)
4. Click `Save contact point`

### 3. Configure the Alert Rule

1. Click `Alert rules` on the left side
2. Click `New Alert Rule`
3. For the Alert name, set it to `Instance Down`

4. Under `2. Define query and alert condition` make sure that Prometheus data source is selected and write query `up == 0`
5. Press `Run queries` to be sure that it works
6. Also under `2. Define query and alert condition` for Expression make sure that `Threshold` is selected and select `IS BELOW` (there is dropdown menu), inside the value field write `1`

7. Under `3. Add folder and labels` press `New folder`, give it name `TestFolder`

8. Under `4. Set evaluation behavior` for the Pending Period select `None`
9. Also under `4. Set evaluation behavior` press `New evaluation group` and give it a name `TestEvaluationGroup` for the interval set 10s

9. Under `5. Configure notifications` choose the contact point `TestContactPoint`

10. Press `Save rule and exit` on top right side

## 4. Trigger the error

1. In your terminal execute `docker compose kill node-exporter`
2. Go to `http://localhost:3000/alerting/list`
3. Wait until the state is set to firing
4. Go to `https://webhook.site` and make sure that you received the notification
5. In your terminal execute `make start`
6. Go to `http://localhost:3000/alerting/list`, now the error should be resolved

