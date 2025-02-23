import json

with open('simulation/src/data/config.json', 'r') as f:
    config = json.load(f)

MQTT_BROKER = config['mqtt_broker']
MQTT_PORT = config['mqtt_port']
MQTT_UPDATE_INTERVAL = config['update_interval']