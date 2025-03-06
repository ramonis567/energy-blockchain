import random

def get_average_energy_consumption(user_type, user_class):
    if user_type == "consumer" or user_type == "prosumer":
        if user_class == "residential":
            return round(random.uniform(0.3, 2), 3)
        elif user_class == "commercial":
            return round(random.uniform(1, 10), 3)
        elif user_class == "industrial":
            return round(random.uniform(5, 30), 3)
        else:
            return 0
    else:
        return 0
    
def get_average_energy_generation(user_type, user_class):
    if user_type == "prosumer" or user_type == "producer":
        if user_class == "residential":
            return round(random.uniform(0.1, 2.2), 3)
        elif user_class == "commercial":
            return round(random.uniform(0.5, 11), 3)
        elif user_class == "industrial":
            return round(random.uniform(5, 20), 3)
        elif user_class == "energy_supplier":
            return round(random.uniform(5, 25), 3)
        else:
            return 0
    else:
        return 0
    
def get_energy_storage_capacity(user_type):
    if user_type == "battery":
        return round(random.uniform(10, 50), 3)
    else:
        return 0