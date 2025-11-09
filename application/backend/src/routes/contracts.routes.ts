import { Router } from "express";
import { ContractsController } from "../controllers/contracts.controller";

const router = Router();

// CRUD-like / market ops
router.post("/create", ContractsController.create);
router.get("/", ContractsController.list);
router.get("/:id", ContractsController.getById);
router.post("/deliver", ContractsController.reportDelivery);
router.post("/settle", ContractsController.settle);
router.post("/close", ContractsController.close);

export default router;
