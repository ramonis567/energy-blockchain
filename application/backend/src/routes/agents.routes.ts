import { Router } from "express";
import { AgentController } from "../controllers/agents.controller";

const router = Router();

router.get("/", AgentController.getAll);
router.get("/count", AgentController.count);
router.get("/:id", AgentController.getById);
router.post("/register", AgentController.register);

export default router;
