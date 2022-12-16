package com.camunda.consulting;

import io.camunda.zeebe.client.api.response.ActivatedJob;
import io.camunda.zeebe.client.api.worker.JobClient;
import io.camunda.zeebe.spring.client.EnableZeebeClient;
import io.camunda.zeebe.spring.client.annotation.ZeebeWorker;
import org.springframework.boot.SpringApplication;
import org.springframework.boot.autoconfigure.SpringBootApplication;

@SpringBootApplication
@EnableZeebeClient
public class Application {

  public static void main(String[] args) {
    SpringApplication.run(Application.class, args);
  }

  @ZeebeWorker(type = "generates4MiB")
  public void generate4MiBVariable(final JobClient client, final ActivatedJob job) {
    String twoMiBString = generate4MiBString();
    client.newCompleteCommand(job.getKey())
        .variables("{\"twoMiB\": \"" + twoMiBString + "\"}")
        .send()
        .exceptionally(throwable -> {
          client.newFailCommand(job).retries(0).errorMessage(throwable.getMessage()).send();
          throw new RuntimeException("Could not complete job " + job, throwable);});
  }

  private String generate4MiBString() {
    StringBuilder builder = new StringBuilder("");
    for (long i = 0L; i < (1 << 21); i++) {
      builder.append('a');
    }
    return builder.toString();
  }
}
