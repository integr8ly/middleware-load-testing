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
- Manual Scaling is possible : 
```
oc scale dc sso -n <ns-prefix>-sso --replicas=<number-of-replicas>
```

### Relational Database
- There is no auto-scaling support
- There are no metrics provided for scaling
- User sessions are kept in memory
- Heavy reliance on caching 