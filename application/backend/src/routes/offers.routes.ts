import { Router } from "express";
import { OfferController } from "../controllers/offers.controller";

const router = Router();

router.get("/", OfferController.getAll);
router.post("/create", OfferController.create);
router.post("/accept", OfferController.accept);

export default router;
