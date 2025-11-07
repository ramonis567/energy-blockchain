import { Request, Response } from "express";
import { TokenService } from "../services/tokens.service";

const service = new TokenService();

export class TokenController {
  async mint(req: Request, res: Response) {
    try {
      const { agentId, tokenType, amount } = req.body;
      if (!agentId || !tokenType || !amount)
        return res.status(400).json({ error: "Missing fields" });

      const result = await service.mint(agentId, tokenType, amount);
      res.status(200).json(result);
    } catch (err: any) {
      res.status(500).json({ error: err.message });
    }
  }

  async transfer(req: Request, res: Response) {
    try {
      const { from, to, tokenType, amount } = req.body;
      if (!from || !to || !tokenType || !amount)
        return res.status(400).json({ error: "Missing fields" });

      const result = await service.transfer(from, to, tokenType, amount);
      res.status(200).json(result);
    } catch (err: any) {
      res.status(500).json({ error: err.message });
    }
  }
}
