const express = require("express");
const app = express();

app.use(express.json());

app.get("/", (req, res) => {
  res.send("API da Camada de Interligação ativa!");
});

// Adicione endpoints conforme necessário
// Exemplo: /dados, /status

module.exports = app;
