package middlewares

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/opentracing/opentracing-go"
	"github.com/uber/jaeger-client-go"
	"github.com/uber/jaeger-client-go/config"
	"github.com/zhaohaihang/userop_api/global"
	"go.uber.org/zap"

	"io"
)

func Trace() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jaegerConfig := global.ApiConfig.JaegerInfo
		jaegerURL := fmt.Sprintf("%s:%d", jaegerConfig.Host, jaegerConfig.Port)
		cfg := &config.Configuration{
			ServiceName: global.ApiConfig.Name,
			Sampler:     &config.SamplerConfig{Type: jaeger.SamplerTypeConst, Param: 1},
			Reporter:    &config.ReporterConfig{LogSpans: true, LocalAgentHostPort: jaegerURL},
		}
		tracer, closer, err := cfg.NewTracer(config.Logger(jaeger.StdLogger))
		if err != nil {
			zap.S().Fatalw("init trace failed", "err", err)
		}
		opentracing.SetGlobalTracer(tracer)

		defer func(closer io.Closer) {
			err := closer.Close()
			if err != nil {
				zap.S().Fatalw("close trace failed", "err", err)
			}
		}(closer)

		startSpan := tracer.StartSpan(ctx.Request.URL.Path)
		defer startSpan.Finish()

		ctx.Set("tracer", tracer)
		ctx.Set("parentSpan", startSpan)
		ctx.Next()
	}

}
