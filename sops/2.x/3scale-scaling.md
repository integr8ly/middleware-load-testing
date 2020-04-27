# 3Scale Scaling Guide
The purpose of this guide is to outline the existing scaling capabilities of 3Scale, including scaling options for high availability and performance driven implementations along with some known limitations within the product.

**Note :** Pod [anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity) rule is in place across a number of 3scale components. This may result in limited scaling options.

## Component Scaling Options
**Note :** Components are not in any specifc order.


### Backend-Listener
#### Scalability:
- Pod **anti-affinity** rule is important here
- You can scale this component **horizontally** by adding more PODs
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"backend": {"listenerSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```
- *Recommended HA --* Deploy a minimum of 2 PODs, but for best performance on OSD/POC clusters deploy a minimum of 5 PODS.


### Backend-Worker
#### Scalability:
- Pod **anti-affinity** rule is important here
- Critical functionality -- rate limits depend on this component
- Depending on the number of reports you should check the length of the redis queues for jobs
- Can be scaled **horizontally** by adding more PODs as needed.
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"backend": {"workerSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```
- *Recommended HA --* Deploy a minimum of 2 PODs

### Backend-Cron
#### Scalability:
- There should be only one copy of this component running. If it fails, OpenShift will restart a new pod to replace it and it will pick-up where the previous pod ended and process remaining work.
- No need to scale this service
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"backend": {"cronSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```


### System-App
#### Scalability:
- This component scales **horizontally** by adding more PODs
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"system": {"appSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```

### AppSpec-Zync
#### Scalability:
- This component scales **horizontally** by adding more PODs
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"zync": {"appSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```

### QueSpec-Zync
#### Scalability:
- This component scales **horizontally** by adding more PODs
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"zync": {"queSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```

### System-Sidekiq 
#### Scalability:
- Due to the nature of the background jobs they don't usually require to be scaled, one POD should be enough, but both components can be scaled **horizonatally**.
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"system": {"sidekiqSpec": {"replicas": <number-of-replicas> }}}}' --type=merge
```
- *Recommended HA --* Deploy 2 PODs

### Apicast-Staging & Apicast-Production
#### Scalability:
- Pod **anti-affinity** rule is important here
- You can scale this component **horizontally** by adding more PODs
```
oc patch ApiManager 3scale -n <ns-prefix>-3scale --patch '{"spec": {"apicast": {"<staging or production>Spec": {"replicas": <number-of-replicas> }}}}' --type=merge
```
- *Recommended HA --* Deploy a minimum of 2 PODs in different openshift nodes.


