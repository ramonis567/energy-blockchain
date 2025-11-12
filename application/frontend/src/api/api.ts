import axios from "axios";

export const api = axios.create({
  baseURL: "http://localhost:4000",
  headers: {
    "Content-Type": "application/json",
  },
});

/**
 * === AGENTES ===
 */
export async function getAgents() {
  const response = await api.get("/agents");
  return response.data;
}

export async function registerAgent(agent: {
  id: string;
  type: string;
  name: string;
  address: string;
}) {
  const response = await api.post("/agents/register", agent);
  return response.data;
}

/**
 * === TOKENS ===
 */
export async function mintTokens(agentId: string, tokenType: string, amount: string) {
  const response = await api.post("/tokens/mint", {
    agentId,
    tokenType,
    amount,
  });
  return response.data;
}

export async function transferTokens(
  from: string,
  to: string,
  tokenType: string,
  amount: string
) {
  const response = await api.post("/tokens/transfer", {
    from,
    to,
    tokenType,
    amount,
  });
  return response.data;
}

export async function getOffers() {
  const res = await api.get("/offers");
  return res.data;
}

export async function createOffer(payload: {
  id: string;
  sellerId: string;
  energyAmount: number;
  pricePerKWh: number;
}) {
  const res = await api.post("/offers/create", payload);
  return res.data;
}

export async function acceptOffer(payload: { id: string; buyerId: string }) {
  const res = await api.post("/offers/accept", payload);
  return res.data;
}