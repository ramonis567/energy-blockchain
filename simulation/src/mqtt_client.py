import paho.mqtt.client as mqtt
from config import MQTT_BROKER, MQTT_PORT

client = mqtt.Client()

def connect():
    client.connect(MQTT_BROKER, MQTT_PORT)

def disconnect():
    client.disconnect()

def publish_data(topic, data):
    client.publish(topic, data)