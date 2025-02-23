import paho.mqtt.client as mqtt
from config import MQTT_BROKER, MQTT_PORT

def publish_data(topic, data):
    client = mqtt.Client()
    client.connect(MQTT_BROKER, MQTT_PORT)
    client.publish(topic, data)
    client.disconnect()