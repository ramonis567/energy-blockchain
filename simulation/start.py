#!/usr/bin/env python3
import subprocess
import sys
import os
import time
from pathlib import Path

def run_command(command, description, check=True):
    print(f"{description}...")
    print(f"Command: {command}")
    
    try:
        if isinstance(command, str):
            result = subprocess.run(command, shell=True, capture_output=True, text=True)
        else:
            result = subprocess.run(command, capture_output=True, text=True)
        
        if result.returncode == 0:
            print(f"{description} finished successfully")
            # if result.stdout:
            #     print(f"Output: \n{result.stdout.strip()}")
            # return True
        else:
            print(f"Error: {description}")
            print(f"Stderr: {result.stderr}")
            if result.stdout:
                print(f"Stdout: {result.stdout}")
            if not check:
                return False
            raise Exception(f"Fail to execute: {description}")
    
    except Exception as e:
        print(f"Exception in {description}: {e}")
        if check:
            raise
        return False

def create_venv():
    if not Path(".venv").exists():
        run_command("python3 -m venv .venv", "Creating .venv")
    else:
        print("Python .venv already exists")
    
    pip_path = Path(".venv/bin/pip")
    if not pip_path.exists():
        pip_path = Path(".venv/Scripts/pip.exe")  # for windows
    
    if pip_path.exists():
        run_command(f"{pip_path} install -r ./simulation/requirements.txt", "Installing requirements.txt")
    else:
        print("There is no pip in the virtual environment")

def install_mosquitto():
    result = subprocess.run(["which", "mosquitto"], capture_output=True, text=True)
    if result.returncode == 0:
        print("Mosquitto already installed")
        return True
    
    print("Installing Mosquitto MQTT broker...")
    run_command("sudo apt install mosquitto mosquitto-clients -y", "Installing Mosquitto")
    
    return True

def start_mosquitto():
    print("🚀 Starting Mosquitto...")
    run_command("sudo systemctl start mosquitto", "Starting Mosquitto service")
    run_command("sudo systemctl enable mosquitto", "Enabling Mosquitto to start on boot")
    run_command("sudo systemctl status mosquitto --no-pager", "Checking Mosquitto status", check=False)
    
    time.sleep(2)
    
    try:
        test_result = subprocess.run(
            ["timeout", "5s", "mosquitto_sub", "-h", "localhost", "-t", "test/connection", "-C", "1"],
            capture_output=True, text=True, timeout=10
        )
        print("✅ MQTT Broker is running and reachable")
    except:
        print("⚠️ MQTT Broker test failed, please check the Mosquitto service")

def run_simulation():
    print("🎯 Starting simulation.py...")
    
    if not Path("./simulation/src/simulation.py").exists():
        print("Arquivo simulation.py not reached!")
        return False
    
    python_path = Path(".venv/bin/python")
    if not python_path.exists():
        python_path = Path(".venv/Scripts/python.exe")  # Windows
    

    if python_path.exists():
        run_command(f"{python_path} ./simulation/src/simulation.py", "Finally starting simulation")
    else:
        run_command("python3 ./simulation/src/simulation.py", "Finally starting simulation")
    
    return True

def main():
    print("=" * 30)
    print("Starting MQTT Simulation Setup")
    print("=" * 30)
    
    try:
        create_venv()
        install_mosquitto()
        start_mosquitto()
        run_simulation()
        
    except Exception as e:
        print(f"❌ ERROR DURING CONFIGURATION: {e}")
        sys.exit(1)

if __name__ == "__main__":
    main()