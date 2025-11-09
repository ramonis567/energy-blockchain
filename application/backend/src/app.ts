import express from "express";
import dotenv from "dotenv";
import bodyParser from "body-parser"

import agentRoutes from "./routes/agents.routes";
import tokensRoutes from "./routes/tokens.routes";
import offersRoutes from "./routes/offers.routes";
import contractRoutes from "./routes/contracts.routes"

dotenv.config();
const app = express();

app.use(express.json());
app.use(bodyParser.json())

app.get("/", (_, res) => res.send("Fabric Backend running ✅"));
app.use("/agents", agentRoutes);
app.use("/tokens", tokensRoutes); 
app.use("/offers", offersRoutes);
app.use("/contracts", contractRoutes);

export default app;
