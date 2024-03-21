# Zero Trust ILT Course

## Prerequisites

In order to run these labs, you will need access to a device running a Linux OS. These labs have been tested on Ubuntu 
22.04.

You will require the following tools to be installed:

- [Golang](https://go.dev/doc/install)
- [curl](https://everything.curl.dev/get/linux)
- [OpenSSL](https://github.com/openssl/openssl#build-and-install)
- [Docker](https://docs.docker.com/desktop/install/linux-install/)
- [kind](https://kind.sigs.k8s.io/docs/user/quick-start#installation)
- [kubectl](https://kubernetes.io/docs/tasks/tools/)
- [jq](https://jqlang.github.io/jq/download/)
- [Contour](https://projectcontour.io/getting-started/)
- [Make](https://www.gnu.org/software/make/manual/make.html)
- [Helm](https://helm.sh/docs/intro/install/)

## Initial Steps

1. Clone the repository to get started with the labs.

    ```bash
    git clone https://github.com/lftraining/lfs482-labs
    ```

## Lab Process

After cloning the repository, you can begin running through the hands-on exercises provided for this course. Once you've 
selected a lab that you would like to practice, navigate to the directory of that exercise and set an environment 
variable to easily reference it when going through the demo by:

1. Change directory (`cd`) into the lab's folder.
2. Set the `LAB_DIR` environment variable:

    ```bash
    export LAB_DIR=$(pwd)
    ```

This will help us in the future when issuing commands within each of the labs, and minimizes the chance of errors when 
running demos.

## 🚢 Scenario Introduction

Welcome to the high-stakes, high-seas operations of Coastal Containers Ltd., an international freighting company 
specializing in transatlantic shipping. Coastal Containers is not just your average shipping company; it sits at the 
forefront of modernizing shipping logistics using cutting-edge technology. Despite the prosperous growth of the company, 
however, the high-seas can be an unruly place to conduct business...

### 🏴‍☠️ The Challenge

Ahoy! Captain Hashjack and his motley crew of cyber-pirates are sailing the high-seas and have set their sights on our 
precious cargo. These brash brigands are notorious for infiltrating supply ships, not just to loot and leave, but to 
compromise the integrity of our supply lanes and steer ships off their intended course. The menacing members of Captain 
Hashjack's crew could be anywhere, even aboard our own ships, making it impossible to trust anyone out on the open sea!

### 🛡️ Why Zero Trust?

The mantra of **"Never Trust, Always Verify"** has never been more relevant. To thwart the nefarious plans of Captain 
Hashjack and his crew, Coastal Containers is embarking on a journey to implement a Zero Trust Architecture between its 
HQ infrastructure in the UK and its supply lanes across the Atlantic ocean. This voyage will sail you through a series 
of hands-on labs that will equip you with the skills to:

- Implement robust identity and authentication mechanisms for both users and workloads.
- Create and enforce stringent policies at various layers of our architecture.
- Integrate dynamic security measures that adapt to ever-evolving threats.

By the end of this course, you will have a treasure trove of knowledge and skills to navigate through the perilous 
waters of cyber-threats and begin implementing a Zero Trust model that even the dastardly Captain Hashjack would tip his 
hat to.

So get ready to set sail with Coastal Containers on this Zero Trust Voyage!

## Labs

### Lab 0 - [Setup](lab-00-setup/README.md)

- **Purpose**: Download the required binaries and container images for the following labs.

### Lab 1 - [Getting Hands on with PKI](lab-01-pki-basics/README.md)

- **Purpose**: Understand the basics of setting up a Public Key Infrastructure using CFSSL and cert-manager.
- **Learning Outcome**: Learn how to create root and intermediate certificates and how to set up cert-manager.
- **Scenario**: To begin its voyage towards Zero Trust, Coastal Containers must first test the waters and set up a 
Public Key Infrastructure (PKI) *dry*-run to begin securing communication across its fleet of sea freighters.

### Lab 2 - [Installing and Configuring SPIRE from Binaries](lab-02-binary-install/README.md)

- **Purpose**: Configure a SPIRE server and agent on a Linux machine using binaries and config files.
- **Learning Outcome**: Gain hands-on experience in configuring SPIRE components on a Linux machine.
- **Scenario**: After the initial PKI *dry*-run sets sail, Coastal Containers now turns its focus to the companies 
coast-bound headquarters, which need to be securely configured to identify and communicate with shipping freighters 
on the move.

### Lab 3 - [Setup SPIRE on Kubernetes with Kind](lab-03-kind-install/README.md)

- **Purpose**: Deploy and configure a SPIRE server and agent on a Kind Kubernetes cluster.
- **Learning Outcome**: Get familiar with deploying SPIRE in a Kubernetes environment.
- **Scenario**: Coastal Containers is recruiting a new fleet admiral, Captain Kubernetes, to modernize and manage its 
fleet. The first order of business is to set up a secure communication channel for the admiral to relay orders to ship 
captains.

### Lab 4 - [Getting SVIDs with SPIFFE-Helper](lab-04-getting-svids/README.md)

- **Purpose**: Learn how to register workloads and obtain SVIDs (SPIFFE Verifiable Identity Documents).
- **Learning Outcome**: Use spiffe-helper to fetch SVIDs and understand the registration process.
- **Scenario**: Coastal Containers realizes that every ship and even individual shipping containers should have their 
own unique identity. To achieve this, Captain Kubernetes introduces a new identity protocol for all ships and 
containers, making them identifiable and accountable.

### Lab 5 - [Using the Workload API with go-spiffe](lab-05-go-spiffe/README.md)

- **Purpose**: Understand and use the SPIFFE Workload API for service identity.
- **Learning Outcome**: Implement mTLS authentication between services using Go SPIFFE library.
- **Scenario**: After clandestine pirate raids go unreported, Captain Kubernetes introduces a new signaling system using 
the Workload API to securely identify and communicate with each ship in the fleet.

### Lab 6 - [Navigating Basic Authorization with Open Policy Agent](lab-06-opa-basics/README.md)

- **Purpose**: Get introduced to the basics of using Open Policy Agent (OPA) for policy enforcement.
- **Learning Outcome**: Create and test simple policies using OPA's Rego language.
- **Scenario**: Coastal Containers starts drafting policies to determine which ships can enter certain waters and which 
containers can carry specific goods. To enforce these policies, the company begins training its personnel in the basics 
of Open Policy Agent.

### Lab 7 - [Integrating SPIRE with OPA and Envoy](lab-07-spire-opa/README.md)

- **Purpose**: Learn how to integrate SPIRE with OPA for dynamic policy enforcement based on workload identity.
- **Learning Outcome**: Set up a demo showing SPIRE providing identity tokens that OPA uses for decision-making.
- **Scenario**: Coastal Containers needs a way to dynamically enforce policies based on the identity of ships and their 
containers. To do so, they consider integrating SPIRE and OPA to make real-time policy decisions based on freighter 
workload identities.

### Lab 8 - [OpenID Connect Discovery](lab-08-oidc-discovery/README.md)

- **Purpose**: Explore how to integrate SPIRE with external systems through OIDC Discovery Providers.
- **Learning Outcome**: Utilize a simple command-line interface (CLI) tool to validate JWTs without requiring cloud 
environment setup.
- **Scenario**: Coastal Containers is forming alliances with other shipping companies to share certain routes. To verify 
the identity of these external partners, the company explores integrating SPIRE with external OIDC providers.

### Lab 9 - [Deploying SPIRE in High Availability Mode](lab-09-ha-mode/README.md)

- **Purpose**: Learn how to set up a highly available SPIRE deployment using Helm charts on a Kubernetes cluster.
- **Learning Outcome**: Understand the components and configuration needed for a high-availability SPIRE setup.
- **Scenario**: With increased pirate threats, Coastal Containers decides it's time to make their communication systems 
more resilient. The aim is to ensure that Captain Kubernetes can always reach the fleet, even if some ships are 
compromised.


### Lab 10 - [Advanced Configuration 1: Nested SPIRE](lab-10-nested-spire/README.md)

- **Purpose**: Dive into advanced SPIRE configurations by setting up a nested SPIRE topology.
- **Learning Outcome**: Understand the architecture and workflow of nested SPIRE deployments.
- **Scenario**: The company realizes that some of its larger ships operate like independent fleets, with multiple layers 
of hierarchy. To manage this complexity, Coastal Containers investigates the use of nested SPIRE deployments.

### Lab 11 - [Advanced Configuration 2: Federated SPIRE](lab-11-federated-spire/README.md)

- **Purpose**: Learn how to set up federated SPIRE deployments for cross-cluster and cross-organization identity.
- **Learning Outcome**: Understand the mechanisms and protocols used in SPIRE federation.
- **Scenario**: Coastal Containers is planning to branching out their operations to establish a transpacific shipping 
route. To ensure seamless and secure communication across the high-seas, they look into federated SPIRE deployments.

### Lab 12 - [Cilium with SPIRE](lab-12-cilium-spire/README.md)

- **Purpose**: Integrate SPIRE with the Cilium network security project to implement identity-based network policies.
- **Learning Outcome**: Configure Cilium to use SPIRE for assigning network policies based on workload identity.
- **Scenario**: As the fleet grows to staggering heights, Coastal Containers is interested in implementing more 
fine-grained policies for its ships and containers. To this end, they look into integrating SPIRE with the Cilium 
network security project.

# ⚓ Conclusion

🎉 Congratulations, sailor! 🎉 You've successfully navigated the tumultuous seas of cyber security and helped Coastal 
Containers anchor safely in the harbor of Zero Trust. By implementing a robust identity and authentication system, 
configuring dynamic policies, and ensuring high availability, you've thwarted Captain Hashjack's plans at every turn. 
The fleet is now sailing smoother than ever, and even the ever-elusive Hashjack would think twice before attempting to 
infiltrate Coastal Containers network of freighters.

Throughout this course, you've gained invaluable hands-on experience in setting up Public Key Infrastructure (PKI), 
configuring SPIRE servers and agents, integrating with Open Policy Agent, and much more. You've not only mastered the 
art of Zero Trust but also learned how to implement it in complex, real-world scenarios. From nested SPIRE topologies to 
federated deployments, you've begun diving deeper into how to realize Zero Trust with OSS!

But the journey doesn't end here. The skills and knowledge you've acquired are highly transferable to any organization 
seeking to bolster its cybersecurity measures. In today's world, where cyber threats can come from any quarter (or any 
quarterdeck), the principles of Zero Trust are more relevant than ever.

So go ahead, take the helm and set sail towards a more secure future. May your seas be calm and your ships ever secure!
