import { Request, Response } from "express";
import { ContractsService, CreateSupplyContractDTO } from "../services/contracts.service";

export class ContractsController {
  static async create(req: Request, res: Response) {
    try {
      const dto = req.body as CreateSupplyContractDTO;

      // validações básicas e pragmáticas
      if (!dto.id || !dto.sellerID || !dto.buyerID) {
        return res.status(400).json({ message: "id, sellerID e buyerID são obrigatórios" });
      }
      if (!dto.energyTotal || !dto.pricePerKWh) {
        return res.status(400).json({ message: "energyTotal e pricePerKWh são obrigatórios" });
      }
      if (!dto.startDate || !dto.endDate || !dto.settlementFreq) {
        return res.status(400).json({ message: "startDate, endDate e settlementFreq são obrigatórios" });
      }

      const result = await ContractsService.create(dto);
      return res.json({ status: "ok", result: { ...result, id: dto.id } });
    } catch (err: any) {
      console.error("[CONTRACT CREATE ERROR]", err);
      return res.status(500).json({ message: err?.message || "CreateSupplyContract failed" });
    }
  }

  static async getById(req: Request, res: Response) {
    try {
      const { id } = req.params;
      const contract = await ContractsService.getById(id);
      return res.json(contract);
    } catch (err: any) {
      console.error("[CONTRACT GET ERROR]", err);
      return res.status(500).json({ message: err?.message || "GetContract failed" });
    }
  }

  static async list(req: Request, res: Response) {
    try {
      const contracts = await ContractsService.listAll();
      return res.json(contracts);
    } catch (err: any) {
      console.error("[CONTRACT LIST ERROR]", err);
      return res.status(500).json({ message: err?.message || "GetAllContracts failed" });
    }
  }

  static async reportDelivery(req: Request, res: Response) {
    try {
      const { id, deliveredKWh } = req.body as { id: string; deliveredKWh: number };
      if (!id || !deliveredKWh) {
        return res.status(400).json({ message: "id e deliveredKWh são obrigatórios" });
      }
      const result = await ContractsService.reportDelivery(id, deliveredKWh);
      return res.json({ status: "ok", result: { ...result, id, deliveredKWh } });
    } catch (err: any) {
      console.error("[CONTRACT DELIVERY ERROR]", err);
      return res.status(500).json({ message: err?.message || "ReportDelivery failed" });
    }
  }

  static async settle(req: Request, res: Response) {
    try {
      const { id, kwhToSettle } = req.body as { id: string; kwhToSettle: number };
      if (!id || !kwhToSettle) {
        return res.status(400).json({ message: "id e kwhToSettle são obrigatórios" });
      }
      const result = await ContractsService.settlePeriod(id, kwhToSettle);
      return res.json({ status: "ok", result: { ...result, id, kwhToSettle } });
    } catch (err: any) {
      console.error("[CONTRACT SETTLE ERROR]", err);
      return res.status(500).json({ message: err?.message || "SettleContractPeriod failed" });
    }
  }

  static async close(req: Request, res: Response) {
    try {
      const { id } = req.body as { id: string };
      if (!id) return res.status(400).json({ message: "id é obrigatório" });
      const result = await ContractsService.close(id);
      return res.json({ status: "ok", result: { ...result, id } });
    } catch (err: any) {
      console.error("[CONTRACT CLOSE ERROR]", err);
      return res.status(500).json({ message: err?.message || "CloseContract failed" });
    }
  }
}
