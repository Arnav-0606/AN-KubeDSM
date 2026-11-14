import json, sys
import numpy as np
import matplotlib.pyplot as plt

what = input()
with open(sys.argv[1], "r") as f:
    y1 = json.load(f)[what]
    print(f"QASRE mean usage: {sum(y1)/len(y1)}")

with open(sys.argv[2], "r") as f:
    y2 = json.load(f)[what]
    print(f"K8s mean usage: {sum(y2)/len(y2)}")

x = list(range(len(y1)))

fig, ax = plt.subplots()
plt.xlabel("Cycle")
plt.ylabel("Edge Resource Usage")

max_x = len(x) - 1
ax.set(xlim=(0, max_x), ylim=(0, 1))

ax.plot(x, y1, linewidth=2.0, label="QASRE")
ax.plot(x, y2, linewidth=2.0, label="k8s")
ax.legend(loc='center left', bbox_to_anchor=(1, 0.5))

plt.tight_layout()
plt.show()
