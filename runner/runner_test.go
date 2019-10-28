package runner

import (
	"testing"
	"time"
	"fmt"
)

func TestRun(t *testing.T) {
	cost := time.Second * 2
	runnerObj := New(cost)
	runnerObj.AddTasks(TaskFunc(task), TaskFunc(task), TaskFunc(task))
	if err := runnerObj.Start(); err != nil {
		fmt.Printf("runner error:%v\n", err)
	}
	fmt.Println("process runned.")
}

func task(id int) {
	fmt.Printf("exec task, id:%d\n", id)
	time.Sleep(time.Second)
}