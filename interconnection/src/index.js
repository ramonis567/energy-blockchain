require("dotenv").config();

const mqttClient = require("./mqtt/mqttClient");
const app = require("./api/server");

const PORT = process.env.PORT || 3000;

app.listen(PORT, () => {
  console.log(`[API] Servidor rodando na porta ${PORT}`);
});
