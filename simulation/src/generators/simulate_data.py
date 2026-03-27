import random

def simulate_user_consumption(user_type, user_class, user_avg_consumption):
    if user_type == "producer" or user_type == "battery":
        return 0
    else:
        return abs(random.gauss(user_avg_consumption, 0.4167))

def simulate_user_generation(user_type, user_class, user_avg_generation):
    if user_type == "consumer" or user_type == "battery":
        return 0
    else:
        return abs(random.gauss(user_avg_generation, 0.4167))

def simulate_storage(user_type, user_class, energy_storage_cap):
    if user_type == "consumer" or user_type == "producer" or user_type == "prosumer":
        return 0
    else:
        return abs(random.gauss(energy_storage_cap, 0.4167))
