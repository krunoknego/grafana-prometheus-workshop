# Prometheus and Grafana Workshop

## Getting familiar with Kubernetes

In your Terminal execute the following command:

``` sh
minikube start
```

This will start-up kubernetes environment for you.

Usually when working with kubernetes, alias k is used.

Try the following commands

``` sh
kubectl get pods -A
```

``` sh
k get pods -A
```

Compare the output of these two commands. It should be the same.

These commands list all available pods in the cluster.

### Create a pod that runs nginx container

``` sh
k run nginx --image=nginx --port=80
```

This command will run an nginx pod with nginx container in the cluster. 

Use the command that you used previously to list all available pods.

Do you see the new pod nginx? What is the STATUS of the pod?

### Make the pod available outside of the cluster

``` sh
k expose pod/nginx --type=NodePort
```

This command will create a Service that exposes the pod outside of the cluster.

``` sh
k get svc
```

This will show you the k8s services in default namespace.

You should see nginx service.

``` sh
minikube service nginx --url
```

this command should return an URL that you can use to access your nginx container.

In your browser open the link.

If everything worked as expected you should see the "Welcome to nginx!" page.

## Manually installing prometheus in the Kubernetes cluster

First execute the following command (for this command you have to be in /day4 directory)

```sh
k apply -f prometheus.yml
```

Inspect the pods and identify if the prometheus pod is running.

```sh
k get pods
```

List all services and make sure that prometheus-service has been created.
```sh
k get svc
```

List the external URL that you can use to connect to the Prometheus that is running inside the cluster.
``` sh
minikube service prometheus-service --url
```

Open the URL in your web browser. Now you should see the Prometheus GUI that you are familiar with.

In the Prometheus list all targets that Prometheus is currently tracking.

```
up
```

There should be only one. Meaning prometheus is only reporting its own metrics.

## Configure Node Exporter on Kubernets cluster and connect it to Prometheus

Create Node Exporter DaemonSet

``` sh
k apply -f node-exporter.yaml
```

Create a service so that Prometheus Pod can scrape it.

Before executing the command below make sure that you replace <NODE-EXPORTER-POD-NAME>.
You get the pods name by writing `k get pods`, it should be something like this for example

NAME                          READY   STATUS    RESTARTS      AGE
nginx                         1/1     Running   0             116m
node-exporter-lx8zf           1/1     Running   0             74m
prometheus-58cd678dfc-pmqpt   1/1     Running   0             3m37s

In this case you would replace <NODE-EXPORTER-POD-NAME> with node-exporter-lx8zf.


``` sh
k expose pod <NODE-EXPORTER-POD-NAME> --type=ClusterIP --port=9100 --name=node-exporter
```


Update `prometheus.yml` file.
``` yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: prometheus-config
data:
  prometheus.yml: |
    global:
      scrape_interval: 15s
    scrape_configs:
      - job_name: 'prometheus'
        static_configs:
          - targets: ['localhost:9090']
      - job_name: 'node-exporter'
        static_configs:
          - targets: ['node-exporter:9100']
```

Now execute the following two commands:

``` sh
k delete -f prometheus.yml
k create -f prometheus.yml
```

This will recreate the prometheus infrastructure.

Now you can execute

``` sh
minikube service prometheus-service --url
```

Open that URL and in Prometheus GUI write

```
up
```

How many active exporters do you see?

## Connect Prometheus running inside of the cluster to Grafana

In your terminal execute

``` sh
make start
```

This will start your Grafana instance.

Grafana will be available on `http://localhost:3000`

### Add Prometheus Data Source

1. Open the Grafana menu
2. Connections -> Data Sources
3. Press `Add new data source`
4. For Prometheus Server URL enter output of the following command
``` sh
minikube service prometheus-service --url
```
5. Press `Save & test`

### Create a simple dashboard

1. Open the Grafana menu
2. Dashboards
3. Press `New` -> `New dashboard`
4. Add visualization
5. Select prometheus Data Source
6. There are two small tabs `Builder` and `Code`, make sure that `Code` is selected
7. Replace `Enter a PromQL query` with `rate(node_cpu_seconds_total[1m])`
8. Press `Run queries`

## Configure Prometheus with Thanos

1. In `docker-compose.yaml` uncomment the lines from line number 11 till line number 117
2. Execute `make down start`

Before going on with the exercise take a look at what was created. `docker ps`.

Also inspect the `docker-compose.yaml` file. How many separate prometheus instances do you identify?

3. Open a prometheus instance `http://localhost:9091` and execute up
4. Open a prometheus instance `http://localhost:9092` and execute up
5. Open a prometheus instance `http://localhost:9093` and execute up

6. Open thanos `http://localhost:9090` and execute up

What do you see? How many instances are monitored by Thanos?
