# AMQ Online Scaling Guide
The purpose of this guide is to outline the existing scaling capabilities of AMQ Online, including scaling options for high availability and performance driven implementations along with some known limitations within the product.


## AMQ Online Components 
The basic AMQ Online infrastructure consists of the following components: 
 - Operator
 - Address-Space-Controller 
 - Api-server,
 - Console
 - (optional) auth-service.

 There is no scaling options for the above components. It is rare that these components will be under load. The following components make up the messaging infrastructure which will be exposed to load under use:
- Controller
- Console
- Brokers
- Routers

While these components do not scale in the traditional sense. Those components that do scale are the `Routers` and `Brokers` as they experience the most load during use. Typically scaled horizontally and will scale automatically.

### AMQ Online Resource Provisioning
To understand how these components scale we need to first understand how resources are provisioned for AMQ 

When customers want to create addresses they request the necessary infrastructure by creating an address space. An address space provides them with an environment for their addresses. Messaging infrastructure is created specifically for that address space. Note that the customer requests result in new pods being deployed in the (central) infrastructure namespace, not in the project in which the custom resource has been created.

Thus, customers can create as many address spaces as they like and AMQ Online will create the necessary resources as needed (within the confines of the cluster). The customer will not directly see the resources they are creating, nor will they be able to see the amount of  resources they are using. Resources are also deleted as address spaces are deleted. 

## AMQ Online Scaling
Within the address space, the controller provides elastic scaling of the routers and brokers as addresses are created and deleted.

How address spaces scale the brokers and routers is based on the following custom resources: 
- Infra-configs
- Address-space-plans
- Address-plans.

When an address space is created an address-space-plan is specified. The plan sets resource limits such as max number of routers and brokers (or unlimited). The address-space-plan also defines which infra-config to use.

Infra configs define the resources (e.g. memory) allocated to brokers and routers and also defines the minimum number of routers and/or brokers.

Finally, address-plans define the size of an address. Therefore they also describe the expected resource usage of brokers and routers and as such will impact how brokers and routers scale.

Some upstream examples of the above CRs:
- [infra-config](https://github.com/EnMasseProject/enmasse/blob/master/templates/example-plans/020-StandardInfraConfig-default.yaml#L11-L20)
- [address-space-plan](https://github.com/EnMasseProject/enmasse/blob/master/templates/example-plans/020-AddressSpacePlan-standard-small.yaml#L14-L17)
- [address-plan](https://github.com/EnMasseProject/enmasse/blob/master/templates/example-plans/020-AddressPlan-standard-small-queue.yaml#L13-L15)

