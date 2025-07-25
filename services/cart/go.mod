module cart-service

go 1.24.2

require (
	github.com/gin-gonic/gin v1.10.1
	github.com/go-redis/redis/v8 v8.11.5
	github.com/google/uuid v1.6.0
	github.com/mansoormajeed/kubernetes-lab/proto/catalog v0.0.0
	github.com/prometheus/client_golang v1.22.0
	github.com/sirupsen/logrus v1.9.3
	go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin v0.62.0
	go.opentelemetry.io/contrib/instrumentation/google.golang.org/grpc/otelgrpc v0.62.0
	go.opentelemetry.io/otel v1.37.0
	go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp v1.37.0
	go.opentelemetry.io/otel/sdk v1.37.0
	go.opentelemetry.io/otel/trace v1.37.0
	google.golang.org/grpc v1.73.0
)