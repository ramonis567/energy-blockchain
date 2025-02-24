import json
import time
import random     ## REMOVE THIS AFTER IMPLEMENTING THE SIMULATION

from generators.init_user_data import get_average_energy_consumption, get_average_energy_generation, get_energy_storage_capacity
from mqtt_client import publish_data
from config import MQTT_UPDATE_INTERVAL

users_static = []

def simulate():
    global users_static

    while True:
        users_dynamic = json.load(open("simulation/src/data/users.json"))["users"]

        ids_users_dynamic = {user["id"] for user in users_dynamic}
        id_users_static = {user["id"] for user in users_static}

        removed_users = id_users_static - ids_users_dynamic
        added_users = ids_users_dynamic - id_users_static

        users_static = [user for user in users_static if user["id"] not in removed_users]

        for user in users_dynamic:
            if user["id"] in added_users:
                user["average_energy_consumption"] = get_average_energy_consumption(user["type"], user["class"])
                user["average_energy_generation"] = get_average_energy_generation(user["type"], user["class"])
                user["energy_storage_capacity"] = get_energy_storage_capacity(user["type"])

                users_static.append(user)

        # for user in users_dynamic:
        #     user_id= user.get("id")
        #     user_type = user.get("type")
        #     user_class = user.get("class")
        #     topic = f"energy/{user_type}/{user_id}"

        #     payload = {
        #         "user_id": user_id,
        #         "user_type": user_type,
        #         "energy_consumption": random.randint(1, 5),
        #         "energy_generation": random.randint(1, 5)
        #     }

        #     payload_json = json.dumps(payload)
        #     publish_data(topic, payload_json)

        time.sleep(MQTT_UPDATE_INTERVAL)

if __name__ == "__main__":
    simulate()