### Apicurito
#### Scalability:
- You can scale this component **horizontally** by adding more PODs
```
oc patch Apicurito apicurito -n <ns-prefix>-apicurito --patch '{"spec": {"size": <number-of-replicas> }}' --type=merge
```
- *Recommended HA --* Deploy a minimum of 2 PODs, but for best performance on OSD/POC clusters deploy a minimum of 5 PODS.


