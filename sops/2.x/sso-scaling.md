# SSO Scaling Guide
The purpose of this guide is to outline the existing scaling capabilities of SSO, including scaling options for high availability and performance driven implementations along with some known limitations within the product.

## SSO Components
### RH-SSO Server
- There is no auto-scaling support
- Scaling strategies are left to customers
- There are no metrics provided for scaling
- Response time is expensive due to password hashing
- User sessions are kept in memory
- Heavy reliance on caching 

### RHSSO
#### Scalability:
- The default deployment of 1 pod allows a maximum rate of 7 users to login/logout per second. Scaling to 2 pods increases this number to 14 users per second. Adding further pods gives a minimal performance increase
- This component scales **horizontally** by adding more pods:
```
oc patch Keycloak rhsso -n <ns-prefix>-rhsso --patch '{"spec": {"instances": <number-of-replicas> }}' --type=merge
```

### RHSSOUser
#### Scalability:
- You can scale this component **horizontally** by adding more PODs
```
oc patch Keycloak rhssouser -n <ns-prefix>-user-sso --patch '{"spec": {"instances": <number-of-replicas> }}' --type=merge
```

