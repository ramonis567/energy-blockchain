import { Gateway, Wallets, Contract } from "fabric-network";
import * as fs from "fs";
import dotenv from "dotenv";

dotenv.config();

let contractPromise: Promise<Contract> | undefined;

export async function getContract(): Promise<Contract> {
  if (contractPromise) {
    return contractPromise;
  }

  contractPromise = (async () => {
    try {
      // 👉 Load the connection profile (JSON)
      const ccpPath = process.env.CONNECTION_PROFILE!;
      const ccp = JSON.parse(await fs.promises.readFile(ccpPath, "utf8"));

      // Load the local wallet (Admin@org1)
      const wallet = await Wallets.newFileSystemWallet("./src/fabric/wallet");
      const identity = await wallet.get("Admin@org1");
      if (!identity) throw new Error("Admin@org1 not found in wallet");

      // Create gateway and connect using the JSON
      const gateway = new Gateway();
      await gateway.connect(ccp, {
        wallet,
        identity: "Admin@org1",
        discovery: { enabled: true, asLocalhost: true },
      });

      console.log("✅ Connected to Fabric network via connection-org1.json");

      // Access channel and chaincode
      const network = await gateway.getNetwork(process.env.CHANNEL_NAME!);
      const contract = network.getContract(process.env.CHAINCODE_NAME!);

      return contract;
    } catch (error) {
      contractPromise = undefined;
      throw error;
    }
  })();

  return contractPromise;
}
