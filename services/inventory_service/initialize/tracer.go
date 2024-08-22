package initialize

import (
	"fmt"
	"io"

	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/zhaohaihang/inventory_service/global"
	"go.uber.org/zap"
)

func InitTracer() (opentracing.Tracer, io.Closer) {
	jaegerInfo := global.ServiceConfig.JaegerInfo
	jaegerURL := fmt.Sprintf("http://%s:%d/api/traces", jaegerInfo.Host, jaegerInfo.Port)
	cfg := &config.Configuration{
		ServiceName: global.ServiceConfig.Name,
		Sampler:     &config.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
		Reporter:    &config.ReporterConfig{LogSpans: true, CollectorEndpoint: jaegerURL},
	}
	tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
	if err != nil {
		zap.S().Fatalw("New Tracer failed: %s", "err", err.Error())
	}
	return tracer, closer
}
