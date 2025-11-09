import { Request, Response } from "express";
import { OfferService } from "../services/offers.service";

const service = new OfferService();

export class OfferController {
  static async create(req: Request, res: Response) {
    try {
      const { id, sellerId, energyAmount, pricePerKWh } = req.body;
      if (!id || !sellerId) return res.status(400).json({ error: "Missing id or sellerId" });
      const result = await service.create(id, sellerId, energyAmount, pricePerKWh);
      res.json(result);
    } catch (err: any) {
      console.error("[CREATE OFFER ERROR]", err);
      res.status(500).json({ error: err.message });
    }
  }

  static async accept(req: Request, res: Response) {
    try {
      const { id, buyerId } = req.body;
      if (!id || !buyerId) return res.status(400).json({ error: "Missing id or buyerId" });
      const result = await service.accept(id, buyerId);
      res.json(result);
    } catch (err: any) {
      console.error("[ACCEPT OFFER ERROR]", err);
      res.status(500).json({ error: err.message });
    }
  }

  static async getAll(_req: Request, res: Response) {
    try {
      const offers = await service.getAll();
      res.json(offers);
    } catch (err: any) {
      res.status(500).json({ error: err.message });
    }
  }
}
