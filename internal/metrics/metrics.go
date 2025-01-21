package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

type Metrics struct {
	MessagesProcessed   prometheus.Counter
	MessagesFailedTotal prometheus.Counter
	ProcessingTime      prometheus.Histogram
	MessageSize         prometheus.Histogram
	QueueSize           prometheus.Gauge
}

func NewMetrics(namespace, subsystem string) *Metrics {
	return &Metrics{
		MessagesProcessed: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "messages_processed_total",
			Help:      "The total number of processed messages",
		}),
		MessagesFailedTotal: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "messages_failed_total",
			Help:      "The total number of failed messages",
		}),
		ProcessingTime: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "message_processing_duration_seconds",
			Help:      "The time spent processing messages",
			Buckets:   prometheus.DefBuckets,
		}),
		MessageSize: promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "message_size_bytes",
			Help:      "The size of processed messages in bytes",
			Buckets:   []float64{64, 128, 256, 512, 1024, 2048, 4096, 8192},
		}),
		QueueSize: promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "queue_size",
			Help:      "The current number of messages in the queue",
		}),
	}
}
