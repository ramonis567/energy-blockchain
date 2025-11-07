import express from "express";
import dotenv from "dotenv";
import bodyParser from "body-parser"
import agentRoutes from "./routes/agents.routes";
import tokensRoutes from "./routes/tokens.routes"

dotenv.config();
const app = express();

app.use(express.json());
app.use(bodyParser.json())

app.get("/", (_, res) => res.send("Fabric Backend running ✅"));
app.use("/agents", agentRoutes);
app.use("/tokens", tokensRoutes); 

export default app;
