import json
import time

from datetime import datetime, timedelta
from generators.init_user_data import get_average_energy_consumption, get_average_energy_generation, get_energy_storage_capacity
from generators.simulate_data import simulate_user_consumption, simulate_user_generation, simulate_storage
from mqtt_client import publish_data
from config import MQTT_UPDATE_INTERVAL

users_static = []
simulation_time = datetime(2025, 1, 1, 0, 0)  

def simulate():
    global users_static, simulation_time

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

        for user in users_static:
            user_id= user.get("id")
            user_type = user.get("type")
            user_class = user.get("class")

            user_avg_consumption = user.get("average_energy_consumption")
            user_avg_generation = user.get("average_energy_generation")
            energy_storage_cap = user.get("energy_storage_capacity")

            topic = f"energy/{user_type}/{user_id}"

            energy_consumption = simulate_user_consumption(user_type, user_class, user_avg_consumption)
            energy_generation = simulate_user_generation(user_type, user_class, user_avg_generation)
            energy_storage = simulate_storage(user_type, user_class, energy_storage_cap)

            payload = {
                "timestamp": simulation_time.isoformat(),
                "user_id": user_id,
                "user_type": user_type,
                "energy_consumption": round(energy_consumption, 3),
                "energy_generation": round(energy_generation, 3),
                "energy_storage": round(energy_storage, 3)
            }

            payload_json = json.dumps(payload)
            publish_data(topic, payload_json)

        simulation_time = simulation_time + timedelta(minutes=15) 
        time.sleep(MQTT_UPDATE_INTERVAL)

if __name__ == "__main__":
    print("Simulating...")
    simulate()