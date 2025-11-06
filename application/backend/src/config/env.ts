import dotenv from "dotenv";
dotenv.config();

export const ENV = {
  PORT: process.env.PORT || 3000,
  FABRIC_CHANNEL: process.env.FABRIC_CHANNEL || "mychannel",
  FABRIC_CONTRACT: process.env.FABRIC_CONTRACT || "main_cc",
  CONNECTION_PROFILE: process.env.CONNECTION_PROFILE || "",
  WALLET_PATH: process.env.WALLET_PATH || "",
  USER_ID: process.env.USER_ID || "appUser",
};
