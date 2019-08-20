# 3Scale Scaling Guide
The purpose of this guide is to outline the existing scaling capabilities of 3Scale, including scaling options for high availability and performance driven implementations along with some known limitations within the product.

**Note :** Pod [anti-affinity](https://kubernetes.io/docs/concepts/configuration/assign-pod-node/#affinity-and-anti-affinity) rule is in place across a number of 3scale components. This may result in limited scaling options.

## Component Scaling Options
**Note :** Components are ordered in priority of known bottle necks.


1. Backend-Listener
    - Pod **anti-affinity** rule is important here
    - You can scale this component **horizontally** by adding more PODs
        - `oc scale dc backend-listener -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
1. Backend-Worker
    - Pod **anti-affinity** rule is important here
    - Critical functionality -- rate limits depend on this component
    - Depending on the number of reports you should check the lenght of the redis queues for jobs
    - Can be scaled **horizontally** by adding more PODs as needed.
        - `oc scale dc backend-worker -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
    - *Recommended HA --* Deploy a minimum of 2 PODs
    - Configuration Options:
        - You can direct worker PODs to the Redis pod they should use by changing the value for the environment variable `CONFIG_QUEUES_MASTER_NAME` (defaults to “backend-redis:6379/1”). This allows you to have separate Redis PODs for background jobs and for listener data in order to scale them independently.
        - **Note** that if you change the CONFIG_QUEUES_MASTER_NAME variable, the PODs for backend-listener and backend-cron should also have matching variables.
1. System-App
    - This component scales **horizontally** by adding more PODs
        - `oc scale dc system-app -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
1. Apicast-Staging & Apicast-Production
    - Pod **anti-affinity** rule is important here
    - You can scale this component **horizontally** by adding more PODs
        - `oc scale dc apicast-<staging or production> -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
    - You can scale this component **vertically** by deploying one worker for each CPU core available to the apicast process. This is controlled by the variable `APICAST_WORKERS`
    - *Recommended HA --* Deploy a minimum of 2 PODs in different openshift nodes.
1. System-Provider & System-Developer
    - This component scales **horizontally** by adding more PODs
        - `oc scale dc system-<provider or developer> -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
    - **Note :** As it depends on system-mysql, keep an eye open for common mysql problems.
    - *Recommended HA --* Deploy 2 PODs
1. System-Redis
    - **Vertical** scaling of each node (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate PODs on a multi-core node or across nodes
1. Backend-Redis
    - **Vertical** scaling of each node (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate PODs on a multi-core node or across nodes
1. System-memcache
    - **Vertical** scaling of each node (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate PODs on a multi-core node or across nodes
1. System-Sidekiq 
    - Due to the nature of the background jobs they don't usually require to be scaled, one POD should be enough, but both components can be scaled **horizonatally**.
        - `oc scale dc system-sidekiq -n <ns-prefix>-3scale --replicas=<number-of-replicas>`
    - These is also `RAILS_MAX_THREADS` which allows for **vertical** scaling, but is recommended to scale horizontally over this approach.
    - *Recommended HA --* Deploy 2 PODs
1. Backend-Cron
    - There should be only one copy of this component running. If it fails, OpenShift will restart a new pod to replace it and it will pick-up where the previous pod ended and process remaining work.
    - No need to scale this service

## Additional Information

### Redis
- Redis is used for storage of traffic data, and also for the job queue, hence redis performance can affect the performance of Listeners and Workers.

- Redis is a single threaded application, so each redis pod at most can use one core.
A single redis pod does not scale with number of cores on the host node, or the number of host nodes.

- It can be scaled by:
    - Vertical scaling of each node (core) redis runs on: CPU Speed, RAM available
    - Splitting redis DBs to run on separate Pods on a multi-core node or across nodes.

    - Recommend ensuring backend-redis is running on a high-memory, fast core machine that is not heavily loaded by other components

- Redis should be as close as possible in the network to the component(s) using it
