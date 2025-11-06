import { Request, Response } from "express";
import { AgentService } from "../services/agents.service";

const service = new AgentService();

export const AgentController = {
  async getAll(req: Request, res: Response) {
    try {
      const agents = await service.getAllAgents();
      res.json(agents);
    } catch (e: any) {
      res.status(500).json({ error: e.message });
    }
  },

  async getById(req: Request, res: Response) {
    try {
      const data = await service.getAgent(req.params.id);
      res.json(data);
    } catch (e: any) {
      res.status(500).json({ error: e.message });
    }
  },

  async count(req: Request, res: Response) {
    try {
      const n = await service.getCount();
      res.json({ count: n });
    } catch (e: any) {
      res.status(500).json({ error: e.message });
    }
  },

  async register(req: Request, res: Response) {
    try {
      const { id, type, name, address } = req.body;
      const result = await service.registerAgent(id, type, name, address);
      res.json({ status: "ok", result });
    } catch (e: any) {
      res.status(500).json({ error: e.message });
    }
  }
};
