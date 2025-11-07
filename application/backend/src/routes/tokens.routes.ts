import { Router } from "express";
import { TokenController } from "../controllers/tokens.controller";

const router = Router();
const controller = new TokenController();

// POST /tokens/mint
router.post("/mint", (req, res) => controller.mint(req, res));
// POST /tokens/transfer
router.post("/transfer", (req, res) => controller.transfer(req, res));

export default router;
