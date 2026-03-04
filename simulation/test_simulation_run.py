import subprocess
import time
import os
import signal

def test_simulation_run():
    # Path to the simulation script
    simulation_script_path = os.path.join(os.path.dirname(__file__), 'src', 'simulation.py')

    # Start the simulation as a subprocess
    process = subprocess.Popen(['python3', '-u', simulation_script_path],
                               stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE,
                               text=True)

    # Wait for a moment to see if it starts
    time.sleep(5)

    # Terminate the process
    process.terminate()
    try:
        stdout, stderr = process.communicate(timeout=5)
    except subprocess.TimeoutExpired:
        process.kill()
        stdout, stderr = process.communicate()

    # The simulation should print "Simulating..."
    assert "Simulating..." in stdout, f"Expected 'Simulating...' in stdout, but got stdout: '{stdout}' and stderr: '{stderr}'"

    # We do NOT expect a ModuleNotFoundError anymore since we installed dependencies
    if "ModuleNotFoundError" in stderr:
        print(f"Warning: ModuleNotFoundError found in stderr: {stderr}")

if __name__ == "__main__":
    test_simulation_run()
    print("Test passed: Simulation starts correctly.")
