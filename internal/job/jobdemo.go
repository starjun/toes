package job

import "log"

type Job01 struct {
	Test string
}

func (g Job01) Run() {
	log.Println("Hello, ", g.Test)
}
