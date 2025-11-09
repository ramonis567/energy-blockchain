import { submit, evaluate } from "../fabric/invoke";

export class OfferService {
  async create(id: string, sellerId: string, energyAmount: number, pricePerKWh: number) {
    const args = [
      id,
      sellerId,
      energyAmount.toString(),
      pricePerKWh.toString()
    ];
    const result = await submit("CreateOffer", args);
    return { status: "ok", txResult: result.toString() };
  }

  async accept(id: string, buyerId: string) {
    const args = [id, buyerId];
    const result = await submit("AcceptOffer", args);
    return { status: "ok", txResult: result.toString() };
  }

  async getAll() {
    const data = await evaluate("GetAllOffers", []);
    return JSON.parse(data);
  }
}
