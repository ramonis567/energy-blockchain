import axios from "axios";

export const api = axios.create({
  baseURL: "http://localhost:4000",
  headers: {
    "Content-Type": "application/json",
  },
});

// 🔹 Obter todos os agentes
export async function getAgents() {
  const response = await api.get("/agents");
  return response.data;
}

// 🔹 Registrar um ou mais agentes
export async function registerAgent(agent: {
  id: string;
  type: string;
  name: string;
  address: string;
}) {
  // A rota correta conforme o backend é /agents/register
  const response = await api.post("/agents/register", agent);
  return response.data;
}

// 🔹 Mint tokens (adição de saldo)
export async function mintTokens(agentId: string, tokenType: string, amount: string) {
  const response = await api.post("/tokens/mint", { agentId, tokenType, amount });
  return response.data;
}