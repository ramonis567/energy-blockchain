import express from "express";
import dotenv from "dotenv";
import agentRoutes from "./routes/agents.routes";

dotenv.config();
const app = express();
app.use(express.json());

app.get("/", (_, res) => res.send("Fabric Backend running ✅"));
app.use("/agents", agentRoutes);

export default app;
