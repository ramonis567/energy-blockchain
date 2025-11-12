import express from "express";
import dotenv from "dotenv";
import bodyParser from "body-parser";
import cors from "cors";

import agentRoutes from "./routes/agents.routes";
import tokensRoutes from "./routes/tokens.routes";
import offersRoutes from "./routes/offers.routes";
import contractRoutes from "./routes/contracts.routes"

dotenv.config();
const app = express();

app.use(
    cors({
        origin: ["http://localhost:5175", "http://127.0.0.1:5175"],
        methods: ["GET", "POST", "PUT", "DELETE", "OPTIONS"],
        credentials: true,
    })
)

app.use((req, res, next) => {
  res.setHeader("Access-Control-Allow-Origin", "*");
  res.setHeader("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS");
  res.setHeader("Access-Control-Allow-Headers", "Content-Type, Authorization");
  console.log(`CORS OK for ${req.method} ${req.url}`);
  next();
});


app.use(express.json());
app.use(bodyParser.json())

app.get("/", (_, res) => res.send("Fabric Backend running ✅"));
app.use("/agents", agentRoutes);
app.use("/tokens", tokensRoutes); 
app.use("/offers", offersRoutes);
app.use("/contracts", contractRoutes);

export default app;
