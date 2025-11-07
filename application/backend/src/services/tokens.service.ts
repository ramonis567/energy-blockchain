import { submit } from "../fabric/invoke";

export class TokenService {
  /**
   * Mint new tokens (energy or credit) to a given agent.
   */
    async mint(agentId: string, tokenType: string, amount: string) {
        try {
            const result = await submit("Mint", [agentId, tokenType, amount]);
            console.log(`[MINT] ${agentId} +${amount} ${tokenType}`);
            return {
                status: "ok",
                result: {
                success: true,
                action: "mint",
                agentId,
                tokenType,
                amount,
                txResult: result?.toString() || "submitted"
                }
            };
        } catch (err: any) {
            console.error("[MINT ERROR]", err);
            throw new Error(`Mint failed: ${err.message}`);
        }
    }

  /**
   * Transfer tokens from one agent to another.
   */
    async transfer(from: string, to: string, tokenType: string, amount: string) {
        try {
            const result = await submit("Transfer", [from, to, tokenType, amount]);
            console.log(`[TRANSFER] ${from} → ${to} (${amount} ${tokenType})`);
            return {
                status: "ok",
                result: {
                success: true,
                action: "transfer",
                from,
                to,
                tokenType,
                amount,
                txResult: result?.toString() || "submitted"
                }
            };
        } catch (err: any) {
            console.error("[TRANSFER ERROR]", err);
            throw new Error(`Transfer failed: ${err.message}`);
        }
    }
}
