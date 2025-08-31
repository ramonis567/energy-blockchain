const processData = (topic, data) => {
    console.log(`[PROCESSANDO] Dados recebidos no tópico ${topic}:`, data);
  
    // Adicionar aqui lógica para armazenar/processar os dados
    // Exemplo: salvar no banco de dados ou repassar para outra camada
  };
  
  module.exports = { processData };
  