import { evaluate, submit } from "../fabric/invoke";

export interface CreateSupplyContractDTO {
  id: string;
  sellerID: string;
  buyerID: string;
  energyTotal: number;        // kWh
  pricePerKWh: number;        // ENGT/kWh
  startDate: string;          // ISO string
  endDate: string;            // ISO string
  settlementFreq: "DAILY" | "WEEKLY" | "MONTHLY";
  sellerCollateralECR?: number;   // default 0
  buyerCollateralENGT?: number;   // default 0
}

export class ContractsService {
  static async create(dto: CreateSupplyContractDTO) {
    const {
      id, sellerID, buyerID, energyTotal, pricePerKWh,
      startDate, endDate, settlementFreq,
      sellerCollateralECR = 0, buyerCollateralENGT = 0
    } = dto;

    const args = [
      id,
      sellerID,
      buyerID,
      String(energyTotal),
      String(pricePerKWh),
      startDate,
      endDate,
      settlementFreq,
      String(sellerCollateralECR),
      String(buyerCollateralENGT),
    ];

    const result = await submit("CreateSupplyContract", args);
    return { success: true, txResult: result ?? "submitted" };
  }

  static async getById(id: string) {
    const res = await evaluate("GetContract", [id]);
    return JSON.parse(res);
  }

  static async listAll() {
    const res = await evaluate("GetAllContracts", []);
    return JSON.parse(res);
  }

  static async reportDelivery(id: string, deliveredKWh: number) {
    const res = await submit("ReportDelivery", [id, String(deliveredKWh)]);
    return { success: true, txResult: res ?? "submitted" };
  }

  static async settlePeriod(id: string, kwhToSettle: number) {
    const res = await submit("SettleContractPeriod", [id, String(kwhToSettle)]);
    return { success: true, txResult: res ?? "submitted" };
  }

  static async close(id: string) {
    const res = await submit("CloseContract", [id]);
    return { success: true, txResult: res ?? "submitted" };
  }
}
