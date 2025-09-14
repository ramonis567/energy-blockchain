const Joi = require("joi");
const { writeMeasurement } = require("./influxWriter")

const schema = Joi.object({
  timestamp: Joi.string().isoDate().required(),
  user_id: Joi.string().required(),
  user_type: Joi.string().valid("consumer", "producer", "prosumer", "battery").required(),
  energy_consumption: Joi.number().min(0).required(),
  energy_generation: Joi.number().min(0).required(),
  energy_storage: Joi.number().min(0).required(),
})

const processData = (topic, data) => {
  const { error, value } = schema.validate(data, { convert: true, aboutEarly: false });

  if (error) {
    console.warn(`[VALIDAÇÃO] ${topic}:`, error.details.map(d=>d.message).join("; "));
    return;
  }

  console.log(`[PROCESSANDO] ${topic}:\n`, data);
  writeMeasurement(data);
  
  // Adicionar aqui lógica para armazenar/processar os dados
  // Exemplo: salvar no banco de dados ou repassar para outra camada

};
  
  module.exports = { processData };
  