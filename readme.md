# redisKeySample-exporter
  - This is a prometheus exporter to monitor the statistics of Redis keys,including the statistics of bigkey, memkey, hotkey.
  - The principle is to use the redis-cli --bigkeys --hotkeys --memkeys command line for key statistics.
# Parameter Variable Description
| -h | -p | -s                                                                  | -c  | -P|
|:---- | :----|:--------------------------------------------------------------------|:----|:----|
| Redis address example: 127.0.0.1:6379  | Redis password   | Sample type example: big or big\|hot or big\|hot\|mem  default: big |whether to enable concurrent execute command|ports exposed by the service|
## Note:
1. REDIS_ADDR is the required when you need to monitor single-node redis. if you need to monitor redis cluster, leave the value empty.
2. SAMPLE_TYPE has three values that you can choose: big,hot,mem.they correspond respectively to the --bigkeys, --hotkeys, and --memkeys flags of redis-cli.if you leave the value empty, default --bigkeys flag.
3. enable redis hotkeys: redis-cli config set maxmemory-policy allkeys-lfu
# Prometheus integration example
- if you use the static configs of prometheus and monitor the single-node redis ,you can reference the examples below:

      - job_name: 'rediskeysample_exporter'
        scrape_interval: 200s
        metrics_path: "/metrics"
        static_configs:
        - targets: ["rediskeysample_exporter_svc:port"] # example
      
      - job_name: 'rediskeysample_targets'
        scrape_interval: 200s
        static_configs:
        - targets:
          - redis-node-0.redis:6379       # example
          - redis-node-1.redis:6379       # example
        metrics_path: /metrics
        relabel_configs:
        - source_labels: [__address__]
          target_label: __param_target
        - source_labels: [__param_target]
          target_label: instance
        - target_label: __address__
          replacement: rediskeysample_exporter_svc:port # example



- or you can use the kubernetes_sd_configs of prometheus and monitor the redis cluster on k8s , you can reference the examples below:

        - job_name: 'rediskeysample_exporter'
          kubernetes_sd_configs:
          - role: pod
          metrics_path: /metrics
          scrape_interval: 200s
          relabel_configs:
          - source_labels: [__meta_kubernetes_pod_name]
            action: keep
            regex: 'redis.*'
          - source_labels: [__meta_kubernetes_pod_ip, __meta_kubernetes_pod_container_port_number]
            target_label: __param_target
            regex: ^(.+);(.+)$
            replacement: $1:$2
          - source_labels: [__param_target]
            target_label: instance
          - target_label: __address__
            replacement: rediskeysample_exporter_svc:port # example
          - source_labels: [__meta_kubernetes_pod_name]
            target_label: pod_name
          - source_labels: [__meta_kubernetes_namespace]
            action: replace
            target_label: namespace


# Build a binary executable file
-  Run the build.sh file

        sh build.sh

# Deploy to k8s
- [k8s Deployment](redisKeySample-exporter.yaml)

# docker:
       docker pull zhenghan008/rediskeysample-exporter:v1.1.0

# grafana dashboard
- [JSON Model](Redis MEM HOT BIG Key Statistics-grafana-dashboard.json)
- [grafana dashboard example](redis-key.png)