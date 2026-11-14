import json, os

qasre_path = "sim/QASRE.json"
k8s_path = "sim/KUBERNETES_DEFAULT.json"

if os.path.exists(qasre_path) and os.path.exists(k8s_path):
    with open(qasre_path) as f: qasre = json.load(f)
    with open(k8s_path) as f: k8s = json.load(f)

    print("==========================================")
    print("         QASRE vs K8s DEFAULT             ")
    print("==========================================")
    
    for metric in ["edge_usage", "qos", "latency"]:
        if metric in qasre and metric in k8s:
            q_vals = qasre[metric]
            k_vals = k8s[metric]
            
            q_mean = sum(q_vals) / len(q_vals) if q_vals else 0
            k_mean = sum(k_vals) / len(k_vals) if k_vals else 0
            
            improvement = ((k_mean - q_mean) / k_mean * 100) if k_mean != 0 else 0
            
            print(f"\nMetric: {metric.upper()}")
            print(f"  - QASRE Mean:       {q_mean:.4f}")
            print(f"  - K8s Default Mean: {k_mean:.4f}")
            print(f"  - Difference/Delta: {improvement:+.2f}%")
else:
    print("Simulation JSON reports not found. Run 'go run run_eval.go' first.")
