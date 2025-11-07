import { getContract } from "./gateway";

export async function evaluate(fn: string, args: string[]) {
  const contract = await getContract();
  const result = await contract.evaluateTransaction(fn, ...args);
  return result.toString();
}

export async function submit(fn: string, args: string[]) {
  const contract = await getContract();
  const result = await contract.submitTransaction(fn, ...args);
  return result.toString();
}
