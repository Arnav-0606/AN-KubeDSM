with open("sim/sim.go", "r") as f:
    content = f.read()

# Replace any lingering usage of choice with algorithmChoice
content = content.replace('choice + ".json"', 'algorithmChoice + ".json"')

with open("sim/sim.go", "w") as f:
    f.write(content)

print("sim/sim.go fixed for choice reference!")
