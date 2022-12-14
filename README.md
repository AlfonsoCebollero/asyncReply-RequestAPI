# asyncReply-RequestAPI
This project runs a [cadence worker](https://github.com/uber-go/cadence-client) in parallel with a [gin-gonic](https://github.com/gin-gonic/gin) server to implement an
async Reply-Request API.

It is thought to have a local instance of cadence running locally which can be easily accomplished following these [instructions](https://cadenceworkflow.io/docs/get-started/installation/#run-cadence-server-using-docker-compose).
This set up does not provide persintance when deleted, but it can be easily added by docker volume means.

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
If the app is wanted to be run with docker, no modifications are needed, in case it is run directly with go, the hostPort must be changed to "127.0.0.1:7933". Both options are thought to be deployed in a local environment where a cadence instance is already running.

There is also the **domain** entry, which value can be changed to any other **existing domain** within cadence. The workflows will be created inside this domain.

A workflows map, which contains all the available workflows inside the application, is included, so, when a new workflow is added, this value must be updated as well.

## Run with docker
```
>> docker build -t asyncapi . 
>> docker run -p 8080:8080 -d --name AsyncAPI asyncapi
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
