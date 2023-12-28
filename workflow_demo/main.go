package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"workflow/workflow"
)

const TaskGroupApplication = "application"

const (
	TaskFunction_Application_Sync = "application_sync"
)

func (p *ApplicationProcessor) ProvideFuntions() map[string]interface{} {
	return map[string]interface{}{
		TaskFunction_Application_Sync: p.Sync,
	}
}

type ApplicationProcessor struct {
}

func (p *ApplicationProcessor) Sync(ctx context.Context, srt string) (interface{}, error) {
	fmt.Println(srt, "sync", "-------------")
	return nil, nil
}

func main() {
	app := &ApplicationProcessor{}
	ctx := context.Background()
	Backend := workflow.NewRedisBackend("127.0.0.1:6379", "", "")

	client := workflow.NewClientFromBackend(Backend)

	server := workflow.NewServerFromBackend(Backend)

	for k, v := range app.ProvideFuntions() {
		if err := server.Register(k, v); err != nil {
			log.Fatal(err)
		}
	}

	go func() {
		err := server.Run(context.TODO())
		if err != nil {
			log.Println("server error ", err)
		}
		log.Println("server run ", err)
	}()

	steps := []workflow.Step{
		{
			Name:     "sync",
			Function: TaskFunction_Application_Sync,
			Args:     workflow.ArgsOf("abc"),
		},
	}
	id, err := client.SubmitCronTask(ctx, workflow.Task{
		Name:  "task_name",
		Group: TaskGroupApplication,
		Steps: steps,
	}, "49 13 * * ?")
	if err != nil {
		log.Fatal(err)
	}

	client.DisableCronTask(id)
	fmt.Println(id)

	//fmt.Println(client.ListTasks(ctx, TaskGroupApplication, "task_name"))

	log.Println("success")
	time.Sleep(10 * time.Hour)
}
