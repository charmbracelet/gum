import subprocess

print("What's your favorite language?")

result = subprocess.run(["gum", "choose", "Go", "Python"], stdout=subprocess.PIPE, text=True)

print(f"I like {result.stdout.strip()}, too!")
