path = "sim/sim.go"
with open(path, "r") as f:
    content = f.read()

lines = content.splitlines()
cleaned_lines = []
for line in lines:
    if "report = struct" in line or "report.QoS = nil" in line or "report.Edge = nil" in line:
        continue
    cleaned_lines.append(line)
content = "\n".join(cleaned_lines)

target = "func Start() {"
replacement = "func Start() {\n\treport.QoS = nil\n\treport.Edge = nil"
if target in content and "report.QoS = nil" not in content:
    content = content.replace(target, replacement, 1)

with open(path, "w") as f:
    f.write(content)
print("sim/sim.go successfully fixed!")
