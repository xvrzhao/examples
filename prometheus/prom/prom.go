package prom

import (
	"fmt"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
)

var (
	// HandledTasksCounter 是一个 counter 指标收集器，代表已处理任务数。
	HandledTasksCounter *prometheus.CounterVec
	// RemainTasksInQueue 是一个 gauge 指标收集器，代表任务队列中剩余任务数。
	RemainTasksInQueueGauge prometheus.Gauge
	// TaskSecondHistogram 是一个 histogram 指标收集器，统计处理任务耗时时长的区间分布。
	TaskSecondHistogram prometheus.Histogram
	// TaskSecondSummary 是一个 summary 指标收集器，统计处理任务耗时时长的分位数大小。
	TaskSecondSummary prometheus.Summary
)

func init() {
	HandledTasksCounter = prometheus.NewCounterVec(prometheus.CounterOpts{
		Name: "handle_tasks_total",
		Help: "The count of handled tasks.",
	}, []string{"worker_id", "task_type"})

	RemainTasksInQueueGauge = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "remain_tasks_in_queue",
		Help: "Number of remain tasks in queue.",
	})

	TaskSecondHistogram = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name:    "task_second_histogram",
		Help:    "The time of processing task.",
		Buckets: []float64{1000, 2000, 3000, 4000, 5000},
	})

	TaskSecondSummary = prometheus.NewSummary(prometheus.SummaryOpts{
		Name:       "task_second_summary",
		Help:       "The time of processing task.",
		Objectives: map[float64]float64{0.25: 0.05, 0.5: 0.05, 0.75: 0.05, 0.99: 0.001},
	})

	// 注册指标收集器
	prometheus.MustRegister(HandledTasksCounter, RemainTasksInQueueGauge, TaskSecondHistogram, TaskSecondSummary)

	// 为 Prometheus Server 提供指标接口
	mux := http.NewServeMux()
	mux.Handle("/metrics", promhttp.Handler())
	go func() {
		if err := http.ListenAndServe(":9000", mux); err != nil {
			log.Fatal(fmt.Errorf("failed to start metric server: %w", err))
		}
	}()
}
