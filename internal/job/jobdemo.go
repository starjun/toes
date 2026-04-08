// Package job 提供定时任务定义。
//
// 该包包含所有定时任务的实现，使用 jobrunner
// 框架进行任务调度。
//
// 主要任务:
//   - Job01: 示例任务
//   - 其他业务任务
//
// 使用示例:
//
//	jobrunner.Schedule("@every 10s", job.Job01{})
package job

import "log"

type Job01 struct {
	Test string
}

func (g Job01) Run() {
	log.Println("Hello, ", g.Test)
}
