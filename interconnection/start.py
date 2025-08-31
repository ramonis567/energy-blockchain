#!/usr/bin/env python3
import subprocess
import sys
import os
import time

print("=" * 30)
print("Starting Interconnection Layer")
print("=" * 30)

os.chdir("./interconnection/")
result = subprocess.run("npm install", shell=True, capture_output=True, text=True)
print(result.stdout.strip())

print("=" * 30)
print("Running Interconnection Layer in localhost:3000")
print("=" * 30)
result = subprocess.Popen("node ./src/index.js", shell=True, text=True, stdout=subprocess.PIPE, stderr=subprocess.PIPE, universal_newlines=True)

while True:
    stdout_line = result.stdout.readline()
    if stdout_line:
                print(stdout_line, end='')

    if result.poll() is not None:
        break