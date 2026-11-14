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

