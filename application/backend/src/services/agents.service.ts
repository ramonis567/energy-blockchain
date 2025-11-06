import { getContract } from "../fabric/gateway";

export class AgentService {
  async getAllAgents() {
    const contract = await getContract();
    const resultBuffer = await contract.evaluateTransaction("GetAllAgents");
    const resultString = resultBuffer.toString();

    console.log("DEBUG: GetAllAgents raw response =>", resultString);

    if (!resultString || resultString.trim().length === 0) {
      console.warn("⚠️ Empty response from Fabric for GetAllAgents");
      return [];
    }

    try {
      return JSON.parse(resultString);
    } catch (e) {
      console.error("❌ JSON parse error:", e, "Raw:", resultString);
      throw e;
    }
  }

  async getAgent(id: string) {
    const contract = await getContract();
    const res = await contract.evaluateTransaction("GetAgentFullInfo", id);
    const str = res.toString();
    if (!str || str.trim().length === 0) return {};
    return JSON.parse(str);
  }

  async getCount() {
    const contract = await getContract();
    const res = await contract.evaluateTransaction("GetAgentCount");
    return parseInt(res.toString());
  }

  async registerAgent(id: string, type: string, name: string, address: string) {
    const contract = await getContract();
    await contract.submitTransaction("RegisterAgent", id, type, name, address);
    return { id, status: "registered" };
  }
}
