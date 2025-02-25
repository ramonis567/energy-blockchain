import random
import numpy as np

def simulate_user_consumption(user_type, user_class, user_avg_consumption):
    return user_avg_consumption + random.uniform(-0.5, 0.5)

def simulate_user_generation(user_type, user_class, user_avg_generation):
    return user_avg_generation + random.uniform(-0.5, 0.5)

def simulate_storage(user_type, user_class, energy_storage_cap):
    return energy_storage_cap + random.uniform(-0.5, 0.5)
