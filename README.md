# asyncReply-RequestAPI
This project runs a [cadence worker](https://github.com/uber-go/cadence-client) in parallel with a [gin-gonic](https://github.com/gin-gonic/gin) server to implement an
async Reply-Request API.

It is thought to have a local instance of cadence running locally which can be easily accomplished following these [instructions](https://cadenceworkflow.io/docs/get-started/installation/#run-cadence-server-using-docker-compose).
This set up does not provide persintance when stopped, but it can be easily added by docker volume means.

I have inspired myself from this [medium post](https://medium.com/stashaway-engineering/building-your-first-cadence-workflow-e61a0b29785) to implement the cadence structure.

## Configuration file
Inside resources folder, a configuration file can be found with the following content:
```
env: "development"
cadence:
  domain: "test-domain"
  service: "cadence-frontend"
  hostPort: "host.docker.internal:7933"
  serverBaseURL: "http://localhost:8080"
  workflows:
    waitingwf: "WaitingWorkflow"

```


## Workflows
The implemented worker only counts with one workflow, which is designed to simulate long async tasks that are in charge of the worker.
When accessing to the local cadence web instance the created workflows can be seen:

![image](https://user-images.githubusercontent.com/34543261/185809387-9188a116-2f45-4a21-934d-fea52fb24170.png)

Its results are also available, this workflow has an activity that is performed after the waiting time that just returns "Completed!".

![image](https://user-images.githubusercontent.com/34543261/185809432-6f78f01e-f806-45f6-b09c-9bc1574e6933.png)

## API and swagger
The endpoints documentation can be found in the app swagger, which is available in: http://localhost:8080/swagger/index.html

![image](https://user-images.githubusercontent.com/34543261/185809499-3568cf53-00fc-45e7-baab-18cb6ff9bbe8.png)
