const mqtt = require("mqtt");
const { MQTT_BROKER, MQTT_PORT, MQTT_TOPICS } = require("../config/mqttConfig");
const { processData } = require("../services/dataProcessor");

const client = mqtt.connect(`mqtt://${MQTT_BROKER}:${MQTT_PORT}`);

client.on("connect", () => {
  console.log("[MQTT] Conectado ao broker!");
  
  MQTT_TOPICS.forEach(topic => {
    client.subscribe(topic, (err) => {
      if (err) {
        console.error(`[MQTT] Erro ao se inscrever no tópico ${topic}:`, err);
      } else {
        console.log(`[MQTT] Inscrito no tópico: ${topic}`);
      }
    });
  });
});

client.on("message", (topic, message) => {
  try {
    const data = JSON.parse(message.toString());
    console.log(`[RECEBIDO] Tópico: ${topic} | Dados:`, data);
    
    processData(topic, data);

  } catch (error) {
    console.error(`[ERRO] Falha ao processar mensagem do tópico ${topic}:`, error);
  }
});

module.exports = client;