package job

import "log"

type Job01 struct {
	Test string
	Cnt  int
}

func (g *Job01) Run() {
	log.Println("Hello, ", g.Test, g.Cnt)
	g.Cnt++
}
