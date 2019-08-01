# AMQ Online Perf Setup

   * [Description of test plans](https://docs.google.com/document/d/1d8HXkpuxtHu1vFoUPmx0BW8AyVIAcvmyFXlKzVK9woI/edit?usp=drive_web&ouid=114847416645920015745) taken from original work carried out by the EnMasse team and adapted to Integreatly` needs

# Test setup

This guide covers setting up/tearing down tests, changing storage settings etc.

## Prerequisite

- Integreatly POC Cluster 

## Running tests
Ansible can be used to automate the configuration and execution of the PERF tests. The steps needed are detailed below

## Using Ansible
First, we need to update our inventory to point to your master cluster. Navigate to the `ansible` directory and run the following command
```
cp inventories/inventory.template inventories/inventory
```
Update the `[master]` variable to the public URL of your master node

To setup the testing infrastructure, this will include testing `addressspaces` and `addresses` along with the `maestro` framework run the following:
```
ansible-playbook -i inventories/inventory playbooks/install.yml
```
Once the playbook has completed run the following to execute the load tests
```
ansible-playbook -i inventories/inventory playbooks/execute_tests.yml
```
Test results can be found at `reports-your-subdomain.com`

## Manual Setup and Execution
### Creating addresses
To meet the description use cases outlined above a number of addresses need to be created.

First, we need remove the standard authservice and to create a none authentication service.
```
oc project openshift-enmasse
oc delete authenticationservice standard-authservice
oc create -f enmasse/authservice/none-authservice.yaml
```
Now we will create our address spaces and addresses in a different project
```
oc new-project maestro-test-addresses
oc create -f enmasse/addressspace/standard.yaml
oc create -f enmasse/addressspace/brokered.yaml

# Wait until the address spaces `standard` and `brokered` status equals true
oc get addressspace standard -o jsonpath={.status.isReady}
oc get addressspace brokered -o jsonpath={.status.isReady}

oc create -f enmasse/address/integreatly_standard.yaml 
oc create -f enmasse/address/integreatly_brokered.yaml 
```
**Note** to retrieve the address space messaging endpoints for the `standard` and `brokered` host names
```
# standard
oc get addressspace standard -o 'jsonpath={.status.endpointStatuses[?(@.name=="messaging")].externalHost}'
# brokered
oc get addressspace brokered -o 'jsonpath={.status.endpointStatuses[?(@.name=="messaging")].externalHost}'
```

## Deploying test infrastructure

## Running tests

In order to isolate the maestro infra and workers from the SUT hosts, all deployments should have a node affinity so 
that you can control where Maestro components are running. 

The affinity is defined by applying the `maestro-node-role` label to a node and setting it to one of the following 
predefined values;

* worker
* infra

Nodes with the `worker` label have an affinity towards worker type components. These are: the Maestro Worker and the 
Maestro Agent. Nodes with the `infra` label have an affinity towards infra type components. These are low resource 
intensive components that provide support for the test execution or the infrastructure required for running Maestro. 
These components are: the Maestro Inspector, Maestro Reports and the Maestro Broker.    

Nodes can be labeled like this:

```
oc label node my-node-01.corp.com maestro-node-role=infra
```

### Create new Maestro Project
```
oc new-project maestro
```

### Deploy Broker

To deploy the Mosquitto broker:

```
oc apply -f maestro/oc/broker/mosquitto-deployment.yaml -f maestro/oc/broker/cluster-service.yaml -f maestro/oc/broker/external-service.yaml
```

### Deploy The Reports Tool

To deploy the reports, first we create volume claim to store the data permanently:

```
oc apply -f maestro/oc/reports/reports-data-pvc.yaml
```

Then we create the services, the deployment and expose the services:

```
oc apply -f maestro/oc/reports/reports-service.yaml -f maestro/oc/reports/reports-deployment.yaml
oc expose -f maestro/oc/reports/reports-service.yaml --hostname=reports.apps.<<my.hostname.com>>
```

**Note**: be sure to replace `<<my.hostname.com>>` with the correct hostname for the reports server.

### Deploy worker
The worker connects to the maestro broker through the broker service.

To deploy the worker:

```
oc apply -f maestro/oc/worker/worker-deployment.yaml
```

Scale the workers up to the amount of replicas that you need for your test: 

```
oc scale deployment --replicas=2 maestro-worker
```

### Deploy inspector (optional)
**Note** Both the inspector and the agent are optional as they are currently experimental components and may cause some unwanted affects.

To deploy the inspector:

```
oc apply -f maestro/oc/inspector/inspector-deployment.yaml
``` 
### Deploy agent (optional)

To deploy the agent:

```
oc apply -f maestro/oc/agent/agent-deployment.yaml
``` 

### Executing Tests on the Client
The maestro client is configured using 2 config maps. The test-scripts config map is used to load the test scripts you wish to select from. This configmap can be created by referencing the directory containing the scripts on creation:
```
oc create configmap test-scripts --from-file=./maestro/test-scripts/src/main/groovy/singlepoint
```
Each test case is run by creating a configmap with the parameters for the test, and creating the pod that runs the client and collects the reports. The reports are stored in a pod volume that you can access through a sidecar container after running the client.

To run a test first you must update `testcase-1.yaml` to match your Receive Url
```
oc apply -f maestro/oc/client/testcase-1.yaml
oc apply -f maestro/oc/client/client.yaml
```
The test is finished once the job is marked as complete. If the test fails it will be rerun. You can view the progress in the client pod's logs. 
You can now start a new test case by replacing the client-config configmap with a new test case, and recreating the pod.
