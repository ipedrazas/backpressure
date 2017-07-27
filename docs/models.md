# Models

## API

```
version: v1
kind: LoadJob
meta
    uid
    creationTimestamp
service
    uid: svc.uid 
    name: svc.name
    namespace: svc.namespace
    labels: svc.meta.labels
    selector: svc.spec.selector
spec
    agents: num agents
    threads: num threads
    sleep: sleep between requests
    app: app 
status: [running, finished, aborted]
results:[
    podmetrics:
        timestamp:
        container:

]

```