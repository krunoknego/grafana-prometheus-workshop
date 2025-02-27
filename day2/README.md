
# Prometheus and Grafana Workshop

## Exercise 1
Write a query that selects only the CPU time spent (node_cpu_seconds_total) in “system” mode.

## Exercise 2
Write a query that selects all CPU metrics (node_cpu_seconds_total) except those in “idle” mode on CPU “0”.

## Exercise 3
Write a query that selects metrics (node_cpu_seconds_total) where the mode label doesn't start with “i” (using regex) on CPU "1".

## Exercise 4
For each network device on a node, compute the total network traffic by adding the bytes received (node_network_receive_bytes_total) and the bytes transmitted (node_network_transmit_bytes_total). Use vector matching on the instance and device labels.

## Exercise 5
Calculate the difference between idle and system CPU times across all CPUs on a node (node_cpu_seconds_total). Remember, vectors need to match!

## Exercise 6
Your goal is to combine filesystem availability data with build information of node_exporter to correlate available storage with the version of the exporter running on each node. In this exercise, you will use the metric node_filesystem_avail_bytes (which shows the available bytes on filesystems) and the metric node_exporter_build_info (which contains build details, including the version label). Multiply these two metrics while matching on the common labels instance and job, and use group_left(version) to attach the extra version label from node_exporter_build_info to your result.

## Exercise 7
Write a query to monitor the current value of active memory using a gauge metric (node_memory_Active_bytes). Execute it multiple times, waiting inbetween approx. 5 seconds. Notice how the value goes up and down.

## Exercise 8
Write a query to calculate the per-second increase (rate) of a counter metric (node_cpu_seconds_total) over the last 5 minutes.
Check the graph.

## Exercise 9
Write a query to calculate the instantaneous per-second increase (irate) of node_cpu_seconds_total over the last 5 minutes.
Check the graph. Compare the graph to the graph in exercise number 8.

## Exercise 10
Your task is to analyze HTTP request duration data using the histogram metric prometheus_http_request_duration_seconds_bucket. Over the last 5 minutes, calculate the 90th percentile of request durations. To achieve this, first determine the maximum observed value for each bucket during that time window (max_over_time), and then use the histogram_quantile function to derive the 90th percentile value.

Try other quantiles and see how the value changes. Notice that you can put any quantile that you want.

## Exercise 11
Write a query to retrieve the 75th percentile of garbage collection duration from a summary metric go_gc_duration_seconds.
Try to select another percentile, like for example 99th and see what is the result going to be.
