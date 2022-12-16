package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"

	"github.com/camunda/zeebe/clients/go/v8/pkg/entities"
	"github.com/camunda/zeebe/clients/go/v8/pkg/worker"
	"github.com/camunda/zeebe/clients/go/v8/pkg/zbc"
)

func main() {

	shutdownBarrier := make(chan bool, 1)
	SetupShutdownBarrier(shutdownBarrier)

	config := zbc.ClientConfig{UsePlaintextConnection: true, GatewayAddress: "localhost:26500"}
	client, err := zbc.NewClient(&config)
	if err != nil {
		panic(err)
	}

	w := client.NewJobWorker().
		JobType("generates2MiB").
		Handler(func(c worker.JobClient, job entities.Job) {

			ctx, cancelFn := context.WithTimeout(context.Background(), 5*time.Second)
			defer cancelFn()

			request, _ := client.NewCompleteJobCommand().JobKey(job.Key).VariablesFromString("{\"twoMiBVar\": \"" + strings.Repeat("b", 1<<22) + "\"}")
			_, err = request.Send(ctx)
			if err != nil {
				log.Printf("failed to complete job with key %d: [%s]", job.Key, err)
				ctx, cancelFn = context.WithTimeout(context.Background(), 5*time.Second)
				defer cancelFn()
				_, err = client.NewFailJobCommand().JobKey(job.Key).Retries(0).ErrorMessage("Failed to complete job because of error " + err.Error()).Send(ctx)
			} else {
				log.Printf("completed job %d successfully", job.Key)
			}
		}).
		MaxJobsActive(10).
		RequestTimeout(10 * time.Second).
		PollInterval(1 * time.Second).
		Name("goGenerates2MiB").
		Open()
	defer w.Close()
	<-shutdownBarrier
}

func SetupShutdownBarrier(done chan bool) {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT, syscall.SIGTERM)
	go func() {
		<-sigs
		done <- true
	}()
}
