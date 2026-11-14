import re

with open("sim/sim.go", "r") as f:
    content = f.read()

# Replace func Start() with func Start(algorithmChoice string)
old_sig = "func Start() {"
new_sig = "func Start(algorithmChoice string) {"

if old_sig in content:
    content = content.replace(old_sig, new_sig, 1)

# Replace fmt.Scan(&choice) logic with the passed parameter
old_scan = """	var choice string
	fmt.Scan(&choice)
	if choice == "QASRE" {"""

new_scan = """	if algorithmChoice == "QASRE" {"""

if old_scan in content:
    content = content.replace(old_scan, new_scan, 1)

with open("sim/sim.go", "w") as f:
    f.write(content)

print("sim.go updated to accept algorithmChoice parameter!")
