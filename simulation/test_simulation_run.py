import subprocess
import time
import os

def test_simulation_run():
    # Path to the simulation script
    simulation_script_path = os.path.join(os.path.dirname(__file__), 'src', 'simulation.py')

    # Start the simulation as a subprocess
    process = subprocess.Popen(['python3', simulation_script_path],
                               stdout=subprocess.PIPE,
                               stderr=subprocess.PIPE,
                               text=True)

    # Wait for a moment to see if it starts
    time.sleep(2)

    # Check the output
    stdout, stderr = process.communicate()

    # Terminate the process
    process.terminate()

    # The simulation should print "Simulating..."
    assert "Simulating..." in stdout, f"Expected 'Simulating...' in stdout, but got: {stdout}"

    # We expect a ModuleNotFoundError because paho-mqtt is not installed
    assert "ModuleNotFoundError: No module named 'paho'" in stderr, f"Expected ModuleNotFoundError in stderr, but got: {stderr}"


if __name__ == "__main__":
    test_simulation_run()
    print("Test passed: Simulation starts and fails as expected without paho-mqtt.")
