# KubeDSM (named ecmus before)
*KubeDSM* is a Kubernetes scheduler module designed to operate on cloud-assisted edge clusters. The following figure displays an overview of the cluster:
<p align="center">
  <img src="https://github.com/AUT-Cloud-Lab/ecmus/blob/main/archive/ClusterOverview.png" width="80%" height="80%">
</p>

## Paper
You can find more details about this scheduler and the experiments in our paper. You can access the paper here: 
[KubeDSM](http://arxiv.org/abs/2501.07130).

## Goal
KubeDSM's goal is to minimize user response time by using edge resources as much as possible. 
To do so, KubeDSM is dynamically migrating pods from cloud to edge. Also, it will migrate pods between edge nodes, to reduce resource fragmentation. This will help the upcoming pods to fit better in the edge.
The scheduler is also able to do batch scheduling (i.e. scheduling multiple pods simultaneously). You can see the scheduler's component diagram in the figure below:
<p align="center">
  <img src="https://github.com/AUT-Cloud-Lab/ecmus/blob/main/archive/ComponentDiagram.png" width="80%" height="80%">
</p>

## Algorithm
The underlying scheduling algorithm takes advantage of the edge cluster's small resource size to make almost optimal decisions to maximize its goal function (Quality Of Service, QoS).
The algorithm's data flow diagram:
<p align="center">
  <img src="https://github.com/AUT-Cloud-Lab/ecmus/blob/main/archive/AlgorithmDFD.png" width="80%" height="80%">
</p>
The nature of the edge environment is error-prone and the scheduler is developed to be resilient to this. The scheduler is able to *heal* its state whenever it gets into a malfunctioned state, and it expects errors and interruptions to happen in each step.

## Experiment
The scheduler's performance has been experimented with various generated multi-pattern workloads for hours in a heterogeneous edge cluster connected to a cloud node. The results show that KubeDSM achieved a **50%** reduction in user response time, a **20%** increase in resource utilization, and fulfilled **100%** of theoretically satisfiable QoS constraints.

## Helper projects
The following is the list of the projects developed alongside KubeDSM for this research:
- [Genny](https://github.com/AUT-Cloud-Lab/genny): Genny generates scenarios based on the given config. It is used to generate normal distributed and wavy scenarios.
- [DrStress](https://github.com/AUT-Cloud-Lab/DrStress): DrStress implements the given scenario to generate workload for the scheduler.
- [MineDraft](https://github.com/AUT-Cloud-Lab/MineDraft): MineDraft extracts draft diagrams and tables from the metrics aggregated during the scenario execution, these drafts are further used to analyze and create diagrams for the paper.
- [Sencillo](https://github.com/AUT-Cloud-Lab/sencillo): As Kubernetes default scheduler does not perform well in heterogeneous edge clusters, several scheduling algorithms are implemented in Sencillo to compete with KubeDSM in the experiments.

- # AN-KubeDSM

*AN-KubeDSM* is an enhanced, priority-aware multi-objective scheduling framework built upon the KubeDSM architecture, optimized for heterogeneous cloud-edge clusters. 

<p align="center">
  <img src="https://github.com/AUT-Cloud-Lab/ecmus/blob/main/archive/ClusterOverview.png" width="80%" height="80%">
</p>

## Overview & Key Enhancements
While traditional edge-cloud schedulers struggle with resource fragmentation and rigid tier prioritization, **AN-KubeDSM** introduces a priority-aware scoring function favoring low-latency edge nodes for critical workloads ($\nu = 1.0$) while dynamically offloading batch workloads ($\nu = 0.0$) to high-capacity cloud nodes.

### Key Features:
- **Priority-Aware Scheduling:** Guarantees critical tier protection by routing 100% of high-priority pods ($\nu = 1.0$) to low-latency edge infrastructure.
- **Dynamic Resource Reclamation:** Intelligent pod migration and multi-objective scoring to minimize user response latency.
- **Resilient Execution:** Built to handle error-prone edge environments with automated state self-healing capabilities.

---

## Paper & Foundation
AN-KubeDSM builds upon the foundational research of KubeDSM. You can review the original design, metrics, and architecture in the published paper:
* [KubeDSM Paper (arXiv:2501.07130)](http://arxiv.org/abs/2501.07130)

---

## Architecture & Component Diagram
The scheduler operates across distributed edge nodes backed by a centralized cloud datacenter:

<p align="center">
  <img src="https://github.com/AUT-Cloud-Lab/ecmus/blob/main/archive/ComponentDiagram.png" width="80%" height="80%">
</p>

---

## Large-Scale Evaluation Results ($N = 100$)
Comparative simulation benchmarks run across 100 heterogeneous workload instances demonstrate AN-KubeDSM's performance advantage over baseline KubeDSM:

```text
=======================================================================================
             AN-KubeDSM vs. Baseline: Scale Summary Report (N=100)                     
=======================================================================================
Evaluation Metric                | Baseline KubeDSM         | AN-KubeDSM (Ours)       
-----------------------------------------------------------------------------------------------
Average Latency                  | 83.8 ms                  | 69.4 ms                 
Critical Tier (nu=1.0) Edge      | 7/29 (24.1%)             | 29/29 (100.0%)          
Normal Tier (nu=0.5) Edge        | 4/29 (13.8%)             | 5/29 (17.2%)            
Batch Tier (nu=0.0) Cloud        | 35/42 (83.3%)            | 42/42 (100.0%)          
=======================================================================================
[Success] Clean scale evaluation successfully completed.

