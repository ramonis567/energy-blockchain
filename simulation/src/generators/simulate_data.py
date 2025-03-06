import random
import numpy as np

def simulate_user_consumption(user_type, user_class, user_avg_consumption):
    if user_type == "producer" or user_type == "battery":
        return 0
    else:
        return np.random.normal(loc=user_avg_consumption, scale=0.4167)

def simulate_user_generation(user_type, user_class, user_avg_generation):
    if user_type == "consumer" or user_type == "battery":
        return 0
    else:
        return np.random.normal(loc=user_avg_generation, scale=0.4167)

def simulate_storage(user_type, user_class, energy_storage_cap):
    if user_type == "consumer" or user_type == "producer" or user_type == "prosumer":
        return 0
    else:
        return np.random.normal(loc=energy_storage_cap, scale=0.4167)
