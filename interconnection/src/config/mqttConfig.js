require("dotenv").config();

module.exports = {
  MQTT_BROKER: process.env.MQTT_BROKER || "localhost",
  MQTT_PORT: process.env.MQTT_PORT || 1883,
  MQTT_TOPICS: [
    "energy/consumer/+",
    "energy/producer/+",
    "energy/prosumer/+",
    "energy/battery/+"
  ]
};