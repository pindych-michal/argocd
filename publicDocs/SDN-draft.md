
# Architectire 

Red Hat OpenShifts networking architecture delivers a robust and scalable foundation for containerized applications. Services provide load balancing based on pod selectors, while routes expose 
these services to external networks. This architecture is well-suited for cloud-native microservices.
However, applications running as virtual machines on OpenShift Virtualization may require different networking considerations. 
Organizations with existing server management infrastructure often expect direct, unrestricted access to virtual machines. The standard Services and Routes model, 
while effective for containerized workloads, may not fully address the networking requirements of VM-based workloads that rely on traditional network access patterns.


IPAM Options:

- OVN-Kubernetes IPAM: Automatic IP allocation from defined subnets
- External DHCP: Use existing DHCP servers
- Cloud-init: Static IP configuration in VM images
- Disabled: Manual IP configuration after deployment


Physical connection consideration 
- 1 bound interfaces
- 2 bound interface (1 dedicated for VM worklaod ) 

Explanation if VM can have only one network interface: 

https://claude.ai/public/artifacts/6e3bba52-ef44-48af-a284-37b907d66d15



Assumption 
- VM should have 1 interface (bridge ovn or linux) if possible - no need for additional interfaces
- If for some reason VM will have POD interfacves - it should be placed in user defined network (static IP) 

Migration approach
- Lift and shift - investigate if in this approach we need POD networking ? If yes can we replace POD networking with User Defined Networking ?
waht really is the role of the POD interface ? 

By default, pods/VMs get one network interface (the cluster's internal network)
OpenShift provides the following CNI plug-ins for the Multus CNI plug-in to chain:
The Five Main CNI Plug-in Types:
- Brdige
- Host-Device
- IPVLAN
- MACVLAN
- SR-IOV 


Two types of accessing openshift resources from external world
- Kubernetes native (ingress/routes, NodePort, MetaLLB)
- Virtual machine native (bridging, routing, dedicate physical interface )


Three types of Virtual Machine networks:
- bridge types  (Vmware portgroups)
- routing type  (Vmware segmetns and T1/T0 distributed routers)
- dedicateed physical connection 


| OpenShift Virtualization menu item | VMware comparable                   | Explanation                                                                                                                                                                      |
|------------------------------------|--------------------------------------|----------------------------------------------------------------------------------------------------------------------------------------------------------------------------------|
| NodeNetworkConfigurationPolicy<br>(NMState Operator needed) | vSwitch/DvSwitch                     | Desired network configuration on cluster nodes                                                                                                                                   |
| NodeNetworkState<br>(NMState Operator needed) | Similar to vSwitch/DvSwitch view at ESX/vCenter | Network status on nodes                                                                                                                                                         |
| Service                            | N/A                                  | Layer4 load balancing configurations with self-discovery and automatic DNS internal to the SDN. Combined with ingress LB solutions (MetalLB, cloud LB) to expose services outside the cluster |
| Routes and Ingresses               | NSX Load Balancer                    | **Routes:** Application load balancing configurations to expose web services outside the cluster<br>**Ingresses:** Accessing application with unique hostname                     |
| NetworkPolicy                      | NSX-T Firewall (Microsegmentation Rules) | Manage application-centric network policies                                                                                                                                      |
| NetworkAttachmentDefinitions       | Port Groups                          | Virtual machine connectivity to networks (VLANs, private networks, etc.)                                                                                                        |
| UserDefinedNetwork                 | NSX-T Overlay Segments               | Create and manage overlay network segments                                                                                                                                       |


## Direct access to the VM using OVN bridegs 

## Direct access to the VM using Linux bridegs 

## Direct access to the VM using distributed OVN routing

## Direct access using dedicated host physical interface 

## Security 








# Configuration 

## Direct access to the VM using OVN bridegs 


Node Network Configuration Policy   (shared interface) 

```
apiVersion: nmstate.io/v1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: br-ex-network
spec:
  nodeSelector:
    node-role.kubernetes.io/worker: '' 
  desiredState:
    ovn:
      bridge-mappings:
      - localnet: br-ex-network
        bridge: br-ex 
        state: present

```


Node Network Configuration Policy (dedicated interfaces) 

```
apiVersion: nmstate.io/v1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: localnet-network
spec:
  nodeSelector:
    node-role.kubernetes.io/worker: ""
	
  desiredState:
    interfaces:
    - name: br0
      bridge:
        options: {}
        port:
        - name: ens192
      ipv4:
        dhcp: "false"
        enabled: "false"
      state: up
      type: ovs-bridge	  
    ovn:
      bridge-mappings:
      - localnet: localnet-network       
		bridge: br0
        state: present
```

ovs-vsctl set bridge br0 stp_enable=false

https://grok.com/share/bGVnYWN5LWNvcHk_5c577409-dc58-4dc7-ab7a-064192dd3f88

Network attachment definition 
```

spec:
  config: |-
    {
        "cniVersion": "0.3.1",
        "name": "br-ex-network"",
        "type": "ovn-k8s-cni-overlay",
        "netAttachDefName": "michalp/a559772",
        "topology": "layer2"
    }
	
spec:
  config: |-
    {
        "cniVersion": "0.4.0",
        "name": "br-ex-network",
        "type": "ovn-k8s-cni-overlay",
        "netAttachDefName": "michalp/a559772-2",
        "topology": "localnet"
    }
# check diffrtent cniVersion, why is that ? 
```


## Direct access to the VM using Linux bridegs 

Node Network Configuration Policy 


```
apiVersion: nmstate.io/v1
kind: NodeNetworkConfigurationPolicy
metadata:
  name: br-flat
spec:
  nodeSelector:
    node-role.kubernetes.io/worker: ""
  desiredState:
    interfaces:
      - name: br-flat
        description: Linux bridge with enp11s0 as a port
        type: linux-bridge
        state: up
        ipv4:
          dhcp: false
          enabled: false
        bridge:
          options:
            stp:
              enabled: false
          port:
            - name: enp11s0
```

https://grok.com/share/bGVnYWN5LWNvcHk_b0914f36-0fbb-439e-b856-199dcd18ce63


## Direct access to the VM using distributed OVN routing 

## Direct access using dedicated host physical interface 

## Security


# Glossary 

------------------------------------------
## Generic Terms 

Container Network Interface (CNI) - Is a standard specification and set of libraries that enables different container network providers to plug into Kubernetes, 
handling essential tasks like assigning IP addresses to pods, setting up network connectivity, and managing network resource.

OVN Kubernetes

Open vSwitch (OVS) 

Multus - Is a container network interface (CNI) plugin for Kubernetes that enables attaching multiple network interfaces to pods.

Operator 

Custom resource definition (CRD)

-------------------------------------------
## Kubernetise networking terms 

Services - (Service types: ClusterIP, Load Balancer, NodePort, ExternalkName, Headless )

Routes 

Ingress 

Network Policies 

-------------------------------------------
## Virtualization networking terms terms 

NMSTAE operator - Manage the network configuration of the cluster nodes themselves in a declarative, Kubernetes-native way.

Node network configuration policy (NNCP)

Node network configuration enactment (NNCE) 

Network Attachment Definition (NAD) - Describes how to attach an additional network interface to a pod or VM.

User Defined Network (UDN)

Cluster User Defined Network (CUDN)

--------------------------------------------

# Trobuelshooting  


0. Global networking trubleshooting

```
oc get network/cluster -o yaml

```

2. Access direcly to the worker nodes

```
oc get nodes
oc debug node/worker1.dev-ocp.openshift.local
```

2. Basic network troubleshooting on worker nodes



3. Applaying nncp

```
oc.exe apply -f ovn-bridge.yml
oc.exe apply -f linux-bridge.yml
oc get nnce
oc get nncp

```

5. Troubleshooting PODs 

```
oc.exe get vmi -o wide 
oc.exe exec -it virt-launcher-rhel-10-turquoise-wombat-53-52jwb -- bin/bash
virtctl.exe console windows-2k25-black-alpaca-86
```



# Links specyfic 

Official redhat documeantion: 

https://docs.redhat.com/en/documentation/openshift_container_platform/4.18/html-single/virtualization/index#virt-connecting-vm-to-ovn-secondary-network

Official redhat blog (3 diffrent scenarios for 1 and 2 dedicatged physical interfaces):

https://www.redhat.com/en/blog/access-external-networks-with-openshift-virtualization

Linux brdige setup:

https://medium.com/@ahmeddraz/connect-openshift-to-the-external-network-70a0362d5d03

OVN brdige setup:

https://blog.epheo.eu/articles/openshift-localnet/index.html

Openshift Virtualziation cookbook  (good section for troubleshooting): 

https://redhatquickcourses.github.io/ocp-virt-cookbook/ocp-virt-cookbook/1/networking/index.html

Youtube videos: 

- https://www.youtube.com/watch?v=vqhJokrTzbs
- https://www.youtube.com/watch?v=RWjvzNH1d0A





# Links general 

Nmstate operator configuration examples:

https://nmstate.io/examples.html#nmstate-state-examples

Openshift examples: 

https://examples.openshift.pub/kubevirt/networking/

https://examples.openshift.pub/

OVN Kuberneties:

https://ovn-kubernetes.io/

https://docs.ovn.org/_/downloads/en/latest/pdf/

MetaLB Configuration:

https://myopenshiftblog.com/openshift-virtualization-networking-101-routes-and-metallb-to-load-balance-vms/

OKD networking documetnation:

https://docs.okd.io/4.19/networking/networking_overview/understanding-networking.html

NSX-T to OVN comparasion: 

https://cloud.redhat.com/learning/learn:high-level-guide-red-hat-openshift-virtualization-vmware-admin/resource/resources:understanding-red-hat-openshift-networking-options-vmware-admin

Exploring the Capabilities of OpenShift Networking:

https://medium.com/@kiran.soft/exploring-the-capabilities-of-openshift-networking-4de7dc51c1c9

https://meatybytes.io/posts/openshift/ocp-features/overview/core-foundations/networking/cni/

https://medium.com/@yakovbeder/mastering-openshift-multi-networking-a-deep-dive-into-udn-and-cudn-912339b3e813   (UDN) - DONE !!! 

Openshift Virtualziation cookbook  (good section for troubleshooting): 

https://redhatquickcourses.github.io/ocp-virt-cookbook/ocp-virt-cookbook/1/networking/index.html

https://github.com/RedHatQuickCourses/

- openshift-virt-roadshow-partner-converg
- architect-the-ocpvirt
- ocpvirt-migration
- rhoso-arch
- rhoso-intro
- map-vSphere-ocpvirt
- devspaces-plugins

https://redhatquickcourses.github.io/

Youtube videos:

https://www.youtube.com/watch?v=_1mULoOtTwA&feature=youtu.be
https://www.youtube.com/watch?v=7L4nNNN6lqc
https://www.youtube.com/watch?v=BMLmHgYYfDI












# IPAM

https://chatgpt.com/share/694005bb-c8ec-800a-8098-65eaab0ad234



