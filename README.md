# Azure Kubernetes ServiceLabs

This is a set of scripts and tools use to generate a docker image that will have the l200labs binary used to evaluate your AKS troubleshooting skill.


## prerequisites

Run `az login` before the lab execrise

Run `az aks get-credentials` to download lab cluster config

Installed kubectl CLI

## How it works?

- Deploy a new AKS cluster with the lab environment (5-10 minutes)

- Follow the lab instructions to complete the task

- Troubleshooting (using ASC, Applen and etc)

- Validate the lab

- Manully delete the lab AKS cluster and resource group

### Example

- Login to the Auzre Cloud Shell with bash
- Run `wget https://github.com/yaoshipu/akslabs/releases/download/0.1.0/akslabs` to download the latest tool
- Run `chmod +x akslabs`
- Run `./akslabs describe` to get lab information

  ```bash 
  Note
  -------------
  Please run az login before the labs
  Please install kubectl
  recommanded to create a new resource group for each lab

  Lab 1 Networking
  ----------------
  Run 'akslabs -n lab1cluster -g lab1group deploy lab1' to deploy the lab1 cluster.
  Tasks:
    Run 'kubectl create deployment whoami --image=containous/whoami' to deploy the pod.
    Run 'kubectl get pods' to list pod.
    Run 'akslabs -n lab1cluster -g lab1group validate lab1' to check the results.
  ```

- Run `akslabs -n lab1cluster -g lab1group deploy lab1` to prepare lab environment, you may need to wait 5 to 10 minutes.
- Run `az aks get-credentials -n lab1cluster -g lab1group` to download kube config file
- Complate the above lab 1 tasks
- You will need to use ASC, Applens and Javis tools
- Save screen shots for your troubleshooting steps
- Run 'akslabs -n lab1cluster -g lab1group validate lab1' to check the results.

  ```
  spark@Azure:~$ kubectl create deployment whoami --image=containous/whoami
  deployment.apps/whoami created
  spark@Azure:~$ kubectl get pods
  NAME                     READY   STATUS    RESTARTS   AGE
  whoami-5c8d94f78-b82rg   1/1     Running   0          17s
  spark@Azure:~$ ./akslabs validate laba1
  No help topic for 'laba1'
  spark@Azure:~$ ./akslabs validate lab1
  Congratulations! You have passed the lab1 test. Good job!
  ```

