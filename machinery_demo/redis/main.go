package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/RichardKnop/machinery/v2"
	"github.com/opentracing/opentracing-go"
	"os"
	"time"

	"github.com/google/uuid"
	"github.com/urfave/cli"

	exampletasks "machinery/redis/task"

	"github.com/RichardKnop/machinery/v2/config"
	"github.com/RichardKnop/machinery/v2/log"
	"github.com/RichardKnop/machinery/v2/tasks"

	redisbackend "github.com/RichardKnop/machinery/v2/backends/redis"
	redisbroker "github.com/RichardKnop/machinery/v2/brokers/redis"
	"github.com/RichardKnop/machinery/v2/example/tracers"
	eagerlock "github.com/RichardKnop/machinery/v2/locks/eager"
	opentracinglog "github.com/opentracing/opentracing-go/log"
)

var (
	app *cli.App
)

func init() {
	// Initialise a CLI app
	app = cli.NewApp()
	app.Name = "machinery"
	app.Usage = "machinery worker and send example tasks with machinery send"
	app.Version = "0.0.0"
}

func main() {
	// Set the CLI app commands
	app.Commands = []cli.Command{
		{
			Name:  "worker",
			Usage: "launch machinery worker",
			Action: func(c *cli.Context) error {
				if err := worker(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
		{
			Name:  "send",
			Usage: "send example tasks ",
			Action: func(c *cli.Context) error {
				if err := send(); err != nil {
					return cli.NewExitError(err.Error(), 1)
				}
				return nil
			},
		},
	}

	// Run the CLI app
	_ = app.Run(os.Args)
}

func startServer() (*machinery.Server, error) {
	cnf := &config.Config{
		DefaultQueue:    "machinery_tasks",
		ResultsExpireIn: 3600,
		Redis: &config.RedisConfig{
			MaxIdle:                3,
			IdleTimeout:            240,
			ReadTimeout:            15,
			WriteTimeout:           15,
			ConnectTimeout:         15,
			NormalTasksPollPeriod:  1000,
			DelayedTasksPollPeriod: 500,
		},
	}

	// Create server instance
	broker := redisbroker.New(cnf, "localhost:6379", "", "", 0)
	backend := redisbackend.New(cnf, "localhost:6379", "", "", 0)
	lock := eagerlock.New()
	server := machinery.NewServer(cnf, broker, backend, lock)

	// Register tasks
	tasksMap := map[string]interface{}{
		"split": exampletasks.Split,
	}

	return server, server.RegisterTasks(tasksMap)
}

func worker() error {
	consumerTag := "machinery_worker"

	cleanup, err := tracers.SetupTracer(consumerTag)
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {
		return err
	}

	// The second argument is a consumer tag
	// Ideally, each worker should have a unique tag (worker1, worker2 etc)
	worker := server.NewWorker(consumerTag, 0)

	// Here we inject some custom code for error handling,
	// start and end of task hooks, useful for metrics for example.
	errorHandler := func(err error) {
		log.ERROR.Println("I am an error handler:", err)
	}

	preTaskHandler := func(signature *tasks.Signature) {
		fmt.Println("name = ", signature.Name)
		log.INFO.Println("I am a start of task handler for:", signature.Name)
	}

	postTaskHandler := func(signature *tasks.Signature) {
		log.INFO.Println("I am an end of task handler for:", signature.Name)
	}

	worker.SetPostTaskHandler(postTaskHandler)
	worker.SetErrorHandler(errorHandler)
	worker.SetPreTaskHandler(preTaskHandler)

	return worker.Launch()
}

func send() error {
	cleanup, err := tracers.SetupTracer("sender")
	if err != nil {
		log.FATAL.Fatalln("Unable to instantiate a tracer:", err)
	}
	defer cleanup()

	server, err := startServer()
	if err != nil {
		return err
	}
	server.SetPreTaskHandler(func(signature *tasks.Signature) {
		fmt.Println("server SetPreTaskHandler", signature)
	})

	span, ctx := opentracing.StartSpanFromContext(context.Background(), "send")
	defer span.Finish()

	batchID := uuid.New().String()
	span.SetBaggageItem("batch.id", batchID)
	span.LogFields(opentracinglog.String("batch.id", batchID))

	log.INFO.Println("Starting batch:", batchID)

	t, _ := time.Parse("2006-01-02 15:04:05", "2023-12-26 16:20:00")
	fmt.Println(t)

	eta := time.Now().UTC().Add(time.Second * 5)

	fmt.Println(server.GetRegisteredTaskNames())

	asyncResult, err := server.SendTaskWithContext(ctx, &tasks.Signature{
		Name: "split",
		ETA:  &eta,
		Args: []tasks.Arg{
			{
				Type:  "[]string",
				Value: []string{"cloud", "jinan"},
			},
			{
				Type:  "string",
				Value: "test",
			},
		},
	})
	if err != nil {
		return fmt.Errorf("Could not send task: %s", err.Error())
	}

	results, err := asyncResult.Get(time.Millisecond * 5)
	if err != nil {
		return fmt.Errorf("Getting task result failed with error: %s", err.Error())
	}
	log.INFO.Printf("data = %v", tasks.HumanReadableResults(results))
	data := &exampletasks.Demo{}
	json.Unmarshal([]byte(tasks.HumanReadableResults(results)), data)
	fmt.Println(data)
	return nil
}
