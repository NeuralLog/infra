# Metrics and Monitoring

This document provides guidance on implementing metrics and monitoring for the NeuralLog Tenant Operator.

## Table of Contents

- [Introduction to Metrics and Monitoring](#introduction-to-metrics-and-monitoring)
- [Prometheus Metrics](#prometheus-metrics)
  - [Default Metrics](#default-metrics)
  - [Custom Metrics](#custom-metrics)
  - [Metric Types](#metric-types)
- [Alerting](#alerting)
  - [Alert Rules](#alert-rules)
  - [Alert Receivers](#alert-receivers)
- [Logging](#logging)
  - [Log Levels](#log-levels)
  - [Structured Logging](#structured-logging)
  - [Log Filtering](#log-filtering)
- [Tracing](#tracing)
  - [OpenTelemetry Integration](#opentelemetry-integration)
  - [Trace Context Propagation](#trace-context-propagation)
- [Dashboards](#dashboards)
  - [Grafana Dashboards](#grafana-dashboards)
  - [Dashboard Examples](#dashboard-examples)
- [Health Checks](#health-checks)
  - [Readiness Probe](#readiness-probe)
  - [Liveness Probe](#liveness-probe)
- [Best Practices](#best-practices)

## Introduction to Metrics and Monitoring

Metrics and monitoring are essential for understanding the health and performance of the NeuralLog Tenant Operator. They provide visibility into the operator's behavior and help identify issues before they become critical.

Key components of metrics and monitoring:

1. **Metrics**: Numerical data about the operator's performance
2. **Alerting**: Notifications when metrics exceed thresholds
3. **Logging**: Textual records of events and errors
4. **Tracing**: Detailed records of request flows
5. **Dashboards**: Visual representations of metrics
6. **Health Checks**: Probes to verify the operator's health

## Prometheus Metrics

The NeuralLog Tenant Operator uses Prometheus for metrics collection. Prometheus is a popular open-source monitoring system that collects metrics from configured targets at specified intervals.

### Default Metrics

The controller-runtime framework provides default metrics for the operator:

1. **Controller Metrics**:
   - `controller_runtime_reconcile_total`: Total number of reconciliations per controller
   - `controller_runtime_reconcile_errors_total`: Total number of reconciliation errors per controller
   - `controller_runtime_reconcile_time_seconds`: Time taken to reconcile per controller
   - `controller_runtime_max_concurrent_reconciles`: Maximum number of concurrent reconciles per controller

2. **Workqueue Metrics**:
   - `workqueue_depth`: Current depth of the workqueue
   - `workqueue_adds_total`: Total number of adds to the workqueue
   - `workqueue_queue_duration_seconds`: Time spent in the workqueue before processing
   - `workqueue_work_duration_seconds`: Time spent processing items from the workqueue
   - `workqueue_retries_total`: Total number of retries in the workqueue

3. **Cache Metrics**:
   - `rest_client_requests_total`: Total number of HTTP requests to the Kubernetes API
   - `rest_client_request_latency_seconds`: Latency of HTTP requests to the Kubernetes API

### Custom Metrics

Custom metrics provide insights specific to the NeuralLog Tenant Operator:

```go
var (
    // TenantCreationTotal is the total number of tenants created
    TenantCreationTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "neurallog_tenant_creation_total",
            Help: "Total number of tenants created",
        },
        []string{"result"},
    )

    // TenantDeletionTotal is the total number of tenants deleted
    TenantDeletionTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "neurallog_tenant_deletion_total",
            Help: "Total number of tenants deleted",
        },
        []string{"result"},
    )

    // TenantReconciliationDuration is the duration of tenant reconciliations
    TenantReconciliationDuration = prometheus.NewHistogramVec(
        prometheus.HistogramOpts{
            Name:    "neurallog_tenant_reconciliation_duration_seconds",
            Help:    "Duration of tenant reconciliations in seconds",
            Buckets: prometheus.DefBuckets,
        },
        []string{"phase"},
    )

    // TenantCount is the number of tenants by phase
    TenantCount = prometheus.NewGaugeVec(
        prometheus.GaugeOpts{
            Name: "neurallog_tenant_count",
            Help: "Number of tenants by phase",
        },
        []string{"phase"},
    )

    // ResourceCreationTotal is the total number of resources created
    ResourceCreationTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "neurallog_resource_creation_total",
            Help: "Total number of resources created",
        },
        []string{"resource_type", "result"},
    )

    // ResourceDeletionTotal is the total number of resources deleted
    ResourceDeletionTotal = prometheus.NewCounterVec(
        prometheus.CounterOpts{
            Name: "neurallog_resource_deletion_total",
            Help: "Total number of resources deleted",
        },
        []string{"resource_type", "result"},
    )
)

func init() {
    // Register custom metrics with the global prometheus registry
    metrics.Registry.MustRegister(
        TenantCreationTotal,
        TenantDeletionTotal,
        TenantReconciliationDuration,
        TenantCount,
        ResourceCreationTotal,
        ResourceDeletionTotal,
    )
}
```

### Metric Types

Prometheus supports several metric types:

1. **Counter**: A cumulative metric that only increases
   ```go
   counter := prometheus.NewCounter(prometheus.CounterOpts{
       Name: "neurallog_tenant_creation_total",
       Help: "Total number of tenants created",
   })
   counter.Inc() // Increment by 1
   counter.Add(5) // Increment by 5
   ```

2. **Gauge**: A metric that can increase and decrease
   ```go
   gauge := prometheus.NewGauge(prometheus.GaugeOpts{
       Name: "neurallog_tenant_count",
       Help: "Number of tenants",
   })
   gauge.Inc() // Increment by 1
   gauge.Dec() // Decrement by 1
   gauge.Set(10) // Set to 10
   ```

3. **Histogram**: A metric that samples observations and counts them in configurable buckets
   ```go
   histogram := prometheus.NewHistogram(prometheus.HistogramOpts{
       Name:    "neurallog_tenant_reconciliation_duration_seconds",
       Help:    "Duration of tenant reconciliations in seconds",
       Buckets: prometheus.DefBuckets,
   })
   histogram.Observe(0.5) // Observe a value
   ```

4. **Summary**: Similar to a histogram, but calculates configurable quantiles over a sliding time window
   ```go
   summary := prometheus.NewSummary(prometheus.SummaryOpts{
       Name:       "neurallog_tenant_reconciliation_duration_seconds",
       Help:       "Duration of tenant reconciliations in seconds",
       Objectives: map[float64]float64{0.5: 0.05, 0.9: 0.01, 0.99: 0.001},
   })
   summary.Observe(0.5) // Observe a value
   ```

## Alerting

Alerting notifies operators when metrics exceed thresholds. Prometheus AlertManager handles alerting.

### Alert Rules

Alert rules define when alerts should be triggered:

```yaml
# prometheus-alerts.yaml
groups:
- name: neurallog-tenant-operator
  rules:
  - alert: TenantReconciliationErrors
    expr: sum(rate(controller_runtime_reconcile_errors_total{controller="tenant"}[5m])) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tenant reconciliation errors"
      description: "Tenant reconciliation errors have been detected in the last 5 minutes."

  - alert: TenantReconciliationLatency
    expr: histogram_quantile(0.9, sum(rate(controller_runtime_reconcile_time_seconds_bucket{controller="tenant"}[5m])) by (le)) > 30
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tenant reconciliation latency"
      description: "90th percentile tenant reconciliation latency is above 30 seconds."

  - alert: TenantCreationFailures
    expr: sum(rate(neurallog_tenant_creation_total{result="failure"}[5m])) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tenant creation failures"
      description: "Tenant creation failures have been detected in the last 5 minutes."

  - alert: TenantDeletionFailures
    expr: sum(rate(neurallog_tenant_deletion_total{result="failure"}[5m])) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Tenant deletion failures"
      description: "Tenant deletion failures have been detected in the last 5 minutes."

  - alert: ResourceCreationFailures
    expr: sum(rate(neurallog_resource_creation_total{result="failure"}[5m])) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Resource creation failures"
      description: "Resource creation failures have been detected in the last 5 minutes."

  - alert: ResourceDeletionFailures
    expr: sum(rate(neurallog_resource_deletion_total{result="failure"}[5m])) > 0
    for: 5m
    labels:
      severity: warning
    annotations:
      summary: "Resource deletion failures"
      description: "Resource deletion failures have been detected in the last 5 minutes."

  - alert: OperatorDown
    expr: up{job="neurallog-tenant-operator"} == 0
    for: 5m
    labels:
      severity: critical
    annotations:
      summary: "Operator down"
      description: "The NeuralLog Tenant Operator is down."
```

### Alert Receivers

Alert receivers define where alerts should be sent:

```yaml
# alertmanager-config.yaml
global:
  resolve_timeout: 5m

route:
  group_by: ['alertname', 'job']
  group_wait: 30s
  group_interval: 5m
  repeat_interval: 12h
  receiver: 'slack'
  routes:
  - match:
      severity: critical
    receiver: 'pagerduty'

receivers:
- name: 'slack'
  slack_configs:
  - api_url: 'https://hooks.slack.com/services/T00000000/B00000000/XXXXXXXXXXXXXXXXXXXXXXXX'
    channel: '#alerts'
    send_resolved: true
    title: '{{ .GroupLabels.alertname }}'
    text: '{{ .CommonAnnotations.description }}'

- name: 'pagerduty'
  pagerduty_configs:
  - service_key: 'XXXXXXXXXXXXXXXXXXXXXXXX'
    send_resolved: true
    description: '{{ .CommonAnnotations.description }}'
```

## Logging

Logging provides detailed information about the operator's behavior. The NeuralLog Tenant Operator uses structured logging with different log levels.

### Log Levels

The operator supports the following log levels:

1. **Debug**: Detailed information for debugging
2. **Info**: General information about the operator's behavior
3. **Warning**: Potential issues that don't affect normal operation
4. **Error**: Issues that affect normal operation
5. **Fatal**: Severe issues that prevent the operator from running

```go
// Set up logging
ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

// Log at different levels
logger := log.FromContext(ctx)
logger.V(1).Info("Debug message") // Debug
logger.Info("Info message") // Info
logger.Info("Warning message", "level", "warning") // Warning
logger.Error(err, "Error message") // Error
logger.Error(err, "Fatal message", "level", "fatal") // Fatal
```

### Structured Logging

Structured logging provides context for log messages:

```go
logger.Info("Reconciling tenant",
    "tenant", tenant.Name,
    "namespace", tenant.Status.Namespace,
    "phase", tenant.Status.Phase,
)

logger.Error(err, "Failed to reconcile tenant",
    "tenant", tenant.Name,
    "namespace", tenant.Status.Namespace,
    "phase", tenant.Status.Phase,
)
```

### Log Filtering

Log filtering allows you to focus on specific log messages:

```go
// Set log level
ctrl.SetLogger(zap.New(zap.Level(zapcore.InfoLevel)))

// Filter logs by component
logger.WithValues("component", "reconciler").Info("Reconciling tenant")

// Filter logs by tenant
logger.WithValues("tenant", tenant.Name).Info("Reconciling tenant")
```

## Tracing

Tracing provides detailed information about request flows. The NeuralLog Tenant Operator uses OpenTelemetry for tracing.

### OpenTelemetry Integration

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() (*sdktrace.TracerProvider, error) {
    // Create Jaeger exporter
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://jaeger-collector:14268/api/traces")))
    if err != nil {
        return nil, err
    }

    // Create resource
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String("neurallog-tenant-operator"),
        ),
    )
    if err != nil {
        return nil, err
    }

    // Create tracer provider
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithSampler(sdktrace.AlwaysSample()),
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )

    // Set global tracer provider
    otel.SetTracerProvider(tp)

    return tp, nil
}
```

### Trace Context Propagation

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/trace"
)

func (r *TenantReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
    // Create a span for the reconciliation
    tracer := otel.Tracer("neurallog-tenant-operator")
    ctx, span := tracer.Start(ctx, "Reconcile")
    defer span.End()

    // Add attributes to the span
    span.SetAttributes(
        attribute.String("tenant", req.Name),
    )

    // Create a child span for a specific operation
    ctx, childSpan := tracer.Start(ctx, "GetTenant")
    tenant := &neurallogv1.Tenant{}
    err := r.Get(ctx, req.NamespacedName, tenant)
    childSpan.End()

    if err != nil {
        // Record error in span
        span.RecordError(err)
        span.SetStatus(codes.Error, err.Error())
        return ctrl.Result{}, client.IgnoreNotFound(err)
    }

    // Continue reconciliation...
    return ctrl.Result{}, nil
}
```

## Dashboards

Dashboards provide visual representations of metrics. The NeuralLog Tenant Operator uses Grafana for dashboards.

### Grafana Dashboards

```json
{
  "annotations": {
    "list": [
      {
        "builtIn": 1,
        "datasource": "-- Grafana --",
        "enable": true,
        "hide": true,
        "iconColor": "rgba(0, 211, 255, 1)",
        "name": "Annotations & Alerts",
        "type": "dashboard"
      }
    ]
  },
  "editable": true,
  "gnetId": null,
  "graphTooltip": 0,
  "id": 1,
  "links": [],
  "panels": [
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 0
      },
      "hiddenSeries": false,
      "id": 2,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(neurallog_tenant_count) by (phase)",
          "interval": "",
          "legendFormat": "{{phase}}",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Tenant Count by Phase",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 0
      },
      "hiddenSeries": false,
      "id": 3,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(controller_runtime_reconcile_total{controller=\"tenant\"}[5m]))",
          "interval": "",
          "legendFormat": "Reconciliations",
          "refId": "A"
        },
        {
          "expr": "sum(rate(controller_runtime_reconcile_errors_total{controller=\"tenant\"}[5m]))",
          "interval": "",
          "legendFormat": "Errors",
          "refId": "B"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Reconciliation Rate",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 0,
        "y": 8
      },
      "hiddenSeries": false,
      "id": 4,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "histogram_quantile(0.5, sum(rate(neurallog_tenant_reconciliation_duration_seconds_bucket[5m])) by (le, phase))",
          "interval": "",
          "legendFormat": "{{phase}} (p50)",
          "refId": "A"
        },
        {
          "expr": "histogram_quantile(0.9, sum(rate(neurallog_tenant_reconciliation_duration_seconds_bucket[5m])) by (le, phase))",
          "interval": "",
          "legendFormat": "{{phase}} (p90)",
          "refId": "B"
        },
        {
          "expr": "histogram_quantile(0.99, sum(rate(neurallog_tenant_reconciliation_duration_seconds_bucket[5m])) by (le, phase))",
          "interval": "",
          "legendFormat": "{{phase}} (p99)",
          "refId": "C"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Reconciliation Duration",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "s",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    },
    {
      "aliasColors": {},
      "bars": false,
      "dashLength": 10,
      "dashes": false,
      "datasource": "Prometheus",
      "fieldConfig": {
        "defaults": {
          "custom": {}
        },
        "overrides": []
      },
      "fill": 1,
      "fillGradient": 0,
      "gridPos": {
        "h": 8,
        "w": 12,
        "x": 12,
        "y": 8
      },
      "hiddenSeries": false,
      "id": 5,
      "legend": {
        "avg": false,
        "current": false,
        "max": false,
        "min": false,
        "show": true,
        "total": false,
        "values": false
      },
      "lines": true,
      "linewidth": 1,
      "nullPointMode": "null",
      "options": {
        "alertThreshold": true
      },
      "percentage": false,
      "pluginVersion": "7.3.7",
      "pointradius": 2,
      "points": false,
      "renderer": "flot",
      "seriesOverrides": [],
      "spaceLength": 10,
      "stack": false,
      "steppedLine": false,
      "targets": [
        {
          "expr": "sum(rate(neurallog_resource_creation_total[5m])) by (resource_type, result)",
          "interval": "",
          "legendFormat": "{{resource_type}} ({{result}})",
          "refId": "A"
        }
      ],
      "thresholds": [],
      "timeFrom": null,
      "timeRegions": [],
      "timeShift": null,
      "title": "Resource Creation Rate",
      "tooltip": {
        "shared": true,
        "sort": 0,
        "value_type": "individual"
      },
      "type": "graph",
      "xaxis": {
        "buckets": null,
        "mode": "time",
        "name": null,
        "show": true,
        "values": []
      },
      "yaxes": [
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        },
        {
          "format": "short",
          "label": null,
          "logBase": 1,
          "max": null,
          "min": null,
          "show": true
        }
      ],
      "yaxis": {
        "align": false,
        "alignLevel": null
      }
    }
  ],
  "refresh": "5s",
  "schemaVersion": 26,
  "style": "dark",
  "tags": [],
  "templating": {
    "list": []
  },
  "time": {
    "from": "now-1h",
    "to": "now"
  },
  "timepicker": {},
  "timezone": "",
  "title": "NeuralLog Tenant Operator",
  "uid": "neurallog-tenant-operator",
  "version": 1
}
```

### Dashboard Examples

1. **Tenant Overview**: Shows the number of tenants by phase
2. **Reconciliation Performance**: Shows reconciliation rate and duration
3. **Resource Management**: Shows resource creation and deletion rates
4. **Error Rates**: Shows error rates by type
5. **API Server Interaction**: Shows API server request rates and latencies

## Health Checks

Health checks verify the operator's health. The NeuralLog Tenant Operator provides readiness and liveness probes.

### Readiness Probe

The readiness probe verifies that the operator is ready to handle requests:

```go
// main.go
func main() {
    // ... existing code
    
    // Add readiness probe
    mgr.AddReadyzCheck("ping", healthz.Ping)
    
    // ... existing code
}
```

### Liveness Probe

The liveness probe verifies that the operator is running:

```go
// main.go
func main() {
    // ... existing code
    
    // Add liveness probe
    mgr.AddHealthzCheck("ping", healthz.Ping)
    
    // ... existing code
}
```

## Best Practices

Follow these best practices for metrics and monitoring:

1. **Use Meaningful Metrics**: Choose metrics that provide actionable insights
2. **Set Appropriate Thresholds**: Set alert thresholds that balance sensitivity and noise
3. **Use Structured Logging**: Use structured logging for better filtering and analysis
4. **Include Context in Logs**: Include relevant context in log messages
5. **Use Different Log Levels**: Use different log levels for different types of information
6. **Monitor Resource Usage**: Monitor CPU, memory, and disk usage
7. **Monitor API Server Interaction**: Monitor API server request rates and latencies
8. **Use Distributed Tracing**: Use distributed tracing for complex request flows
9. **Create Comprehensive Dashboards**: Create dashboards that provide a complete view of the operator
10. **Regularly Review Metrics and Logs**: Regularly review metrics and logs to identify issues
