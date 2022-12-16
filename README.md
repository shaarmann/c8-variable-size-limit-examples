# Error Handling when Exceeding Variable Size

Camunda 8 has a limit for the variable size. One command, i.e., completing a job, is limited to a 4 MiB payload. If this payload is exceeded, external task workers fail. If this error is not handled properly, Camunda offers the job again to the external task workers.

## Golang

In Go, a failed complete command returns an error. If the error is not `nil`, a failure command can be raised (decrementing or) setting the number of retries.
```Go
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
```

## Java

In Java, a failed complete command raises an exception. The exception can be caught and a fail command can be send (decrementing or) setting the number of retries.
```Java
String twoMiBString = generate4MiBString();
client.newCompleteCommand(job.getKey())
	.variables("{\"twoMiB\": \"" + twoMiBString + "\"}")
	.send()
	.exceptionally(throwable -> {
	  client.newFailCommand(job).retries(0).errorMessage(throwable.getMessage()).send();
	  throw new RuntimeException("Could not complete job " + job, throwable);});
```