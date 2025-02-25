Python 3.10.12
Mosquitto 2.0.11

Create .venv
pip install requirements.txt


Start mqtt broker:
    Install: sudo apt update && sudo apt install mosquitto mosquitto-clients -y
    Initialize: sudo systemctl start mosquitto
    Verify: sudo systemctl status mosquitto

    Test subscribe: mosquitto_sub -h localhost -t "energy/#" -v
