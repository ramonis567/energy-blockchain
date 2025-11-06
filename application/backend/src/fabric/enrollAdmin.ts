import { Wallets } from "fabric-network";
import * as fs from "fs";
import dotenv from "dotenv";
dotenv.config();

async function main() {
  // Create or open the wallet folder
  const wallet = await Wallets.newFileSystemWallet("./src/fabric/wallet");

  // Load cert and key from .env paths
  const cert = fs.readFileSync(process.env.CERT_PATH!, "utf8");
  const key = fs.readFileSync(process.env.KEY_PATH!, "utf8");

  // Fabric 2.2 format: use "credentials"
  const identity = {
    credentials: {
      certificate: cert,
      privateKey: key,
    },
    mspId: process.env.MSP_ID!,
    type: "X.509",
  };

  await wallet.put("Admin@org1", identity);
  console.log("✅ Admin@org1 successfully enrolled in wallet!");
}

main().catch((err) => {
  console.error("❌ Enrollment failed:", err);
  process.exit(1);
});
