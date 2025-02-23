import json
import time
import random     ## REMOVE THIS AFTER IMPLEMENTING THE SIMULATION

from mqtt_client import publish_data
from config import MQTT_UPDATE_INTERVAL

users = json.load(open("simulation/src/data/users.json"))["users"]

def simulate():
    while True:
        users = json.load(open("simulation/src/data/users.json"))["users"]
        
        for user in users:
            user_id= user.get("id")
            user_name = user.get("name")
            user_type = user.get("type")
            topic = f"energy/{user_type}/{user_id}"

            payload = {
                "user_id": user_id,
                "user_name": user_name,
                "user_type": user_type,
                "energy_consumption": random.randint(1, 5),
                "energy_generation": random.randint(1, 5)
            }

            payload_json = json.dumps(payload)
            publish_data(topic, payload_json)


        time.sleep(MQTT_UPDATE_INTERVAL)

if __name__ == "__main__":
    simulate()