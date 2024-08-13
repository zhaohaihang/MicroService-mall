package initialize

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/flow"
	"go.uber.org/zap"
)

// 限流器
func InitSentinel() {
	err := sentinel.InitDefault()
	if err != nil {
		zap.S().Fatalw("init sentinel failed ", "err", err)
	}
	_, err = flow.LoadRules([]*flow.Rule{
		{
			Resource:               "userop_api",
			TokenCalculateStrategy: flow.Direct,   // 支持两种启动策略：WarmUp 冷启动策略; Direct 直接启动 
			ControlBehavior:        flow.Reject,  // 超过限制时执行的策略，支持两种策略: Reject 直接拒绝 ; Throttling 匀速通过
			Threshold:              10000,        // 允许通过请求数量的阈值
			StatIntervalInMs:       1000,          // 时间间隔
		},
	})
	if err != nil {
		zap.S().Fatalw("sentinel load rules failed", "err", err)
	}
}
