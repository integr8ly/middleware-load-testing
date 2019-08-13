# 3scale scaling on Openshift 3.11, Integreatly 1.4.x
This guide states the known limits of 3scale capacity on Openshift along with scaling options for 3scale Components.

## 3scale Limits
Known limits of 3scale
|| Up to (Requests/Day) | Sustained rate(rps) | Peak rate (x4)(rps)| Compute cores|Router cores|
|---|---|---|---|---|---|
||20M|231|924|6|2|

## 3scale Scaling Components
**Note :** Components are ordered in priority of known bottle necks.

1. Apicast-Staging & Apicast-Production
    - Pod anti-affinity rule is important here
    - You can scale this component **horizontally** by adding more pods
    - You can scale this component **vertically** by deploying one worker for each CPU core available to the apicast process. This is controlled by the variable `APICAST_WORKERS`
1. System-Provider & System-Developer
    - This component scales **horizontally** by adding more PODs
1. Backend-Worker
    - Pod anit-affinity rule is important here
    - Critical functionality -- rate limits depend on this component
    - Depending on the number of reports you should check the lenght of the redis queues for jobs
    - Can be scaled **horizontally** by adding more pods as needed.
1. System-Redis
    - **Vertical** scaling of each nore (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate Pods on a multi-core node or across nodes
1. Backend-Redis
    - **Vertical** scaling of each nore (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate Pods on a multi-core node or across nodes
1. System-memcache
    - **Vertical** scaling of each nore (core) redis run on: 
        - Redis is a single threaded application, so each redis pod at most can use one core. A single redis pod does not scale with number of cores on the host node, or the number of host nodes.
        - CPU Speed, RAM available splitting redis DMs to run on separate Pods on a multi-core node or across nodes
1. System-Sidekiq 
    - Due to the nature of the background jobs they don't usually require to be scaled, one POD should be enough, but both components can be scaled **horizonatally**.
    - These is also `RAILS_MAX_THREADS` which allows for **vertical** scaling, but is recommended to scale horizontally over this approach.
1. Backend-Cron
    - There should be only one copy of this component running. If it fails, OpenShift will restart a new pod to replace it and it will pick-up where the previous pod ended and process remaining work.
    - No need to scale this service
