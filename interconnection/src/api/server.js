const express = require("express");
const app = express();

app.use(express.json());

app.get("/", (req, res) => {
  res.send("Camada de interconexão ativa!");
});

// Adicionar endpoints

module.exports = app;
