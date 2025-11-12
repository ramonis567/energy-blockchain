import { useEffect, useState } from "react";
import { getAgents, mintTokens, transferTokens } from "../api/api";

interface Agent {
  id: string;
  type: string;
  name: string;
  ecr_balance: number;
  engt_balance: number;
  address: string;
  registered_at: string;
}

interface DistributionItem {
  id: string;
  name: string;
  amount: number;
}

export default function Tokens() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [filteredAgents, setFilteredAgents] = useState<Agent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [refreshing, setRefreshing] = useState(false);

  const [selectedAgents, setSelectedAgents] = useState<string[]>([]);
  const [typeFilter, setTypeFilter] = useState("");

  // === Modais ===
  const [showMintModal, setShowMintModal] = useState(false);
  const [showTransferModal, setShowTransferModal] = useState(false);
  const [showRandomModal, setShowRandomModal] = useState(false);

  // === Mint ===
  const [mintToken, setMintToken] = useState("ECR");
  const [mintAmount, setMintAmount] = useState("");
  const [minting, setMinting] = useState(false);

  // === Transfer ===
  const [fromAgent, setFromAgent] = useState("");
  const [toAgent, setToAgent] = useState("");
  const [transferToken, setTransferToken] = useState("ECR");
  const [transferAmount, setTransferAmount] = useState("");
  const [transferring, setTransferring] = useState(false);

  // === Distribuição aleatória ===
  const [randomToken, setRandomToken] = useState("ECR");
  const [randomTotal, setRandomTotal] = useState("");
  const [randomCount, setRandomCount] = useState(5);
  const [distribution, setDistribution] = useState<DistributionItem[]>([]);
  const [distributing, setDistributing] = useState(false);

  // === Buscar agentes ===
  async function fetchAgents() {
    try {
      setLoading(true);
      const data = await getAgents();
      setAgents(data);
      setFilteredAgents(data);
      setSelectedAgents([]);
    } catch (err) {
      console.error(err);
      setError("Erro ao carregar dados de agentes");
    } finally {
      setLoading(false);
    }
  }

  useEffect(() => {
    fetchAgents();
  }, []);

  useEffect(() => {
    let result = [...agents];
    if (typeFilter) result = result.filter((a) => a.type === typeFilter);
    setFilteredAgents(result);
  }, [typeFilter, agents]);

  async function handleRefresh() {
    setRefreshing(true);
    await fetchAgents();
    setRefreshing(false);
  }

  // === Controle de seleção ===
  function toggleAgentSelection(id: string) {
    setSelectedAgents((prev) =>
      prev.includes(id) ? prev.filter((x) => x !== id) : [...prev, id]
    );
  }

  function toggleSelectAll() {
    if (selectedAgents.length === filteredAgents.length) {
      setSelectedAgents([]);
    } else {
      setSelectedAgents(filteredAgents.map((a) => a.id));
    }
  }

  // === Abrir modais ===
  function openMintModal() {
    setShowTransferModal(false);
    setShowRandomModal(false);
    setShowMintModal(true);
  }

  function openTransferModal() {
    setShowMintModal(false);
    setShowRandomModal(false);
    setShowTransferModal(true);
  }

  function openRandomModal() {
    setShowMintModal(false);
    setShowTransferModal(false);
    setShowRandomModal(true);
  }

  // === Mint ===
  async function handleMintTokens() {
    if (selectedAgents.length === 0 || !mintAmount) {
      alert("Selecione agentes e informe a quantidade.");
      return;
    }
    setMinting(true);
    try {
      for (const id of selectedAgents) {
        await mintTokens(id, mintToken, mintAmount);
      }
      alert(`✅ ${mintAmount} ${mintToken} emitidos para ${selectedAgents.length} agentes.`);
      setShowMintModal(false);
      setMintAmount("");
      fetchAgents();
    } catch (err) {
      console.error(err);
      alert("Erro ao emitir tokens.");
    } finally {
      setMinting(false);
    }
  }

  // === Transfer ===
  async function handleTransferTokens() {
    if (!fromAgent || !toAgent || !transferAmount) {
      alert("Preencha todos os campos.");
      return;
    }
    setTransferring(true);
    try {
      await transferTokens(fromAgent, toAgent, transferToken, transferAmount);
      alert("✅ Transferência realizada com sucesso!");
      setShowTransferModal(false);
      fetchAgents();
    } catch (err) {
      console.error(err);
      alert("Erro ao transferir tokens.");
    } finally {
      setTransferring(false);
    }
  }

  // === Distribuição aleatória ===
  function generateRandomDistribution() {
    const total = parseFloat(randomTotal);
    if (!total || randomCount <= 0) {
      alert("Preencha os campos corretamente.");
      return;
    }

    if (randomCount > agents.length) {
      alert("Número de agentes maior que o total disponível.");
      return;
    }

    // Seleciona agentes aleatórios únicos
    const shuffled = [...agents].sort(() => 0.5 - Math.random());
    const selected = shuffled.slice(0, randomCount);

    // Gera pesos aleatórios normalizados
    const weights = Array.from({ length: randomCount }, () => Math.random());
    const sum = weights.reduce((a, b) => a + b, 0);
    const normalized = weights.map((w) => w / sum);

    // Distribui valores
    const distributionResult = selected.map((a, i) => ({
      id: a.id,
      name: a.name,
      amount: parseFloat((normalized[i] * total).toFixed(2)),
    }));

    setDistribution(distributionResult);
  }

  async function confirmRandomDistribution() {
    if (distribution.length === 0) {
      alert("Gere a distribuição primeiro.");
      return;
    }

    setDistributing(true);
    try {
      for (const d of distribution) {
        await mintTokens(d.id, randomToken, d.amount.toString());
      }
      alert(`✅ ${randomTotal} ${randomToken} distribuídos entre ${distribution.length} agentes.`);
      setShowRandomModal(false);
      setDistribution([]);
      setRandomTotal("");
      fetchAgents();
    } catch (err) {
      console.error(err);
      alert("Erro ao distribuir tokens.");
    } finally {
      setDistributing(false);
    }
  }

  // === Totais ===
  const totalECR = filteredAgents.reduce((sum, a) => sum + a.ecr_balance, 0);
  const totalENGT = filteredAgents.reduce((sum, a) => sum + a.engt_balance, 0);

  if (loading) return <p className="text-gray-500">Carregando agentes...</p>;
  if (error) return <p className="text-red-600">{error}</p>;

  return (
    <div className="p-6 bg-gray-50 rounded-lg shadow-sm">
      <div className="flex flex-wrap items-center justify-between mb-6">
        <h2 className="text-2xl font-semibold text-[var(--blueColor)]">Gestão de Tokens</h2>

        <div className="text-right">
          <p className="text-sm text-gray-600">
            Agentes: <strong>{filteredAgents.length}</strong>
          </p>
          <p className="text-sm text-gray-700">
            💰 <strong className="text-green-700">ECR: {totalECR.toFixed(2)}</strong> &nbsp; | &nbsp;
            ⚡ <strong className="text-blue-700">ENGT: {totalENGT.toFixed(2)}</strong>
          </p>
        </div>
      </div>

      <div className="flex flex-wrap gap-4 items-center mb-6">
        <select
          value={typeFilter}
          onChange={(e) => setTypeFilter(e.target.value)}
          className="px-3 py-2 border rounded-md focus:ring-2 focus:ring-[var(--blueColor)]"
        >
          <option value="">Todos os tipos</option>
          <option value="producer">Produtor</option>
          <option value="consumer">Consumidor</option>
          <option value="prosumer">Prosumer</option>
          <option value="distributor">Distribuidor</option>
          <option value="battery">Bateria</option>
        </select>

        <div className="flex gap-3 ml-auto">
          <button
            onClick={handleRefresh}
            disabled={refreshing}
            className={`px-4 py-2 rounded-md text-white ${
              refreshing ? "bg-gray-400" : "bg-[var(--highlightColor)] hover:bg-sky-600"
            }`}
          >
            {refreshing ? "Atualizando..." : "Atualizar"}
          </button>

          <button
            onClick={openMintModal}
            disabled={selectedAgents.length === 0}
            className={`font-medium px-4 py-2 rounded-md transition text-white ${
              selectedAgents.length === 0
                ? "bg-gray-400 cursor-not-allowed"
                : "bg-[var(--greenColor)] hover:bg-green-700"
            }`}
          >
            Emitir ({selectedAgents.length})
          </button>

          <button
            onClick={openTransferModal}
            className="bg-[var(--blueColor)] hover:bg-blue-900 text-white font-medium px-4 py-2 rounded-md"
          >
            Transferir
          </button>

          <button
            onClick={openRandomModal}
            className="bg-purple-600 hover:bg-purple-800 text-white font-medium px-4 py-2 rounded-md"
          >
            Distribuição Aleatória
          </button>
        </div>
      </div>

      {/* === Tabela === */}
      <div className="overflow-x-auto bg-white shadow rounded-lg">
        <table className="min-w-full text-sm text-left text-gray-700">
          <thead className="bg-[var(--blueColor)] text-white">
            <tr>
              <th className="px-4 py-2 text-center">
                <input
                  type="checkbox"
                  checked={
                    selectedAgents.length > 0 &&
                    selectedAgents.length === filteredAgents.length
                  }
                  onChange={toggleSelectAll}
                />
              </th>
              <th className="px-4 py-2">ID</th>
              <th className="px-4 py-2">Nome</th>
              <th className="px-4 py-2">Tipo</th>
              <th className="px-4 py-2 text-right">ECR</th>
              <th className="px-4 py-2 text-right">ENGT</th>
            </tr>
          </thead>
          <tbody>
            {filteredAgents.map((a) => (
              <tr
                key={a.id}
                className={`border-b hover:bg-gray-50 transition-colors ${
                  selectedAgents.includes(a.id) ? "bg-green-50" : ""
                }`}
              >
                <td className="px-4 py-2 text-center">
                  <input
                    type="checkbox"
                    checked={selectedAgents.includes(a.id)}
                    onChange={() => toggleAgentSelection(a.id)}
                  />
                </td>
                <td className="px-4 py-2 font-medium">{a.id}</td>
                <td className="px-4 py-2">{a.name}</td>
                <td className="px-4 py-2 capitalize">{a.type}</td>
                <td className="px-4 py-2 text-right text-green-700 font-semibold">
                  {a.ecr_balance.toFixed(2)}
                </td>
                <td className="px-4 py-2 text-right text-blue-700 font-semibold">
                  {a.engt_balance.toFixed(2)}
                </td>
              </tr>
            ))}
          </tbody>
        </table>
      </div>

      {/* === Modal Mint === */}
      {showMintModal && (
        <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-[60]">
          <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
            <h3 className="text-xl font-semibold text-[var(--greenColor)] mb-4">
              Emitir Tokens
            </h3>

            <div className="flex flex-col gap-3">
              <label className="text-sm font-medium text-gray-700">Tipo de Token</label>
              <select
                value={mintToken}
                onChange={(e) => setMintToken(e.target.value)}
                className="border px-3 py-2 rounded-md"
              >
                <option value="ECR">ECR</option>
                <option value="ENGT">ENGT</option>
              </select>

              <label className="text-sm font-medium text-gray-700">Quantidade</label>
              <input
                type="number"
                min={0.01}
                step="0.01"
                value={mintAmount}
                onChange={(e) => setMintAmount(e.target.value)}
                placeholder="Ex: 100"
                className="border px-3 py-2 rounded-md"
              />
            </div>

            <div className="flex justify-end mt-6 gap-3">
              <button
                onClick={() => setShowMintModal(false)}
                className="px-4 py-2 rounded-md border text-gray-600 hover:bg-gray-100"
              >
                Cancelar
              </button>
              <button
                onClick={handleMintTokens}
                disabled={minting}
                className={`px-4 py-2 rounded-md text-white ${
                  minting ? "bg-gray-400" : "bg-[var(--greenColor)] hover:bg-green-700"
                }`}
              >
                {minting ? "Processando..." : "Confirmar"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* === Modal Transfer === */}
      {showTransferModal && (
        <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-[70]">
          <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
            <h3 className="text-xl font-semibold text-[var(--blueColor)] mb-4">
              Transferir Tokens
            </h3>

            <div className="flex flex-col gap-3">
              <label className="text-sm font-medium text-gray-700">De</label>
              <select
                value={fromAgent}
                onChange={(e) => setFromAgent(e.target.value)}
                className="border px-3 py-2 rounded-md"
              >
                <option value="">Selecione</option>
                {agents.map((a) => (
                  <option key={a.id} value={a.id}>{a.name}</option>
                ))}
              </select>

              <label className="text-sm font-medium text-gray-700">Para</label>
              <select
                value={toAgent}
                onChange={(e) => setToAgent(e.target.value)}
                className="border px-3 py-2 rounded-md"
              >
                <option value="">Selecione</option>
                {agents.map((a) => (
                  <option key={a.id} value={a.id}>{a.name}</option>
                ))}
              </select>

              <label className="text-sm font-medium text-gray-700">Token</label>
              <select
                value={transferToken}
                onChange={(e) => setTransferToken(e.target.value)}
                className="border px-3 py-2 rounded-md"
              >
                <option value="ECR">ECR</option>
                <option value="ENGT">ENGT</option>
              </select>

              <label className="text-sm font-medium text-gray-700">Quantidade</label>
              <input
                type="number"
                min={0.01}
                step="0.01"
                value={transferAmount}
                onChange={(e) => setTransferAmount(e.target.value)}
                placeholder="Ex: 50"
                className="border px-3 py-2 rounded-md"
              />
            </div>

            <div className="flex justify-end mt-6 gap-3">
              <button
                onClick={() => setShowTransferModal(false)}
                className="px-4 py-2 rounded-md border text-gray-600 hover:bg-gray-100"
              >
                Cancelar
              </button>
              <button
                onClick={handleTransferTokens}
                disabled={transferring}
                className={`px-4 py-2 rounded-md text-white ${
                  transferring
                    ? "bg-gray-400 cursor-not-allowed"
                    : "bg-[var(--blueColor)] hover:bg-blue-900"
                }`}
              >
                {transferring ? "Processando..." : "Confirmar"}
              </button>
            </div>
          </div>
        </div>
      )}

      {/* === Modal Distribuição Aleatória === */}
      {showRandomModal && (
        <div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-[80]">
          <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-lg">
            <h3 className="text-xl font-semibold text-purple-700 mb-4">
              Distribuir Tokens Aleatoriamente
            </h3>

            <div className="flex flex-col gap-3 mb-4">
              <label className="text-sm font-medium text-gray-700">Tipo de Token</label>
              <select
                value={randomToken}
                onChange={(e) => setRandomToken(e.target.value)}
                className="border px-3 py-2 rounded-md"
              >
                <option value="ECR">ECR</option>
                <option value="ENGT">ENGT</option>
              </select>

              <label className="text-sm font-medium text-gray-700">Valor Total</label>
              <input
                type="number"
                min={1}
                value={randomTotal}
                onChange={(e) => setRandomTotal(e.target.value)}
                placeholder="Ex: 1000"
                className="border px-3 py-2 rounded-md"
              />

              <label className="text-sm font-medium text-gray-700">Número de Agentes</label>
              <input
                type="number"
                min={1}
                max={agents.length}
                value={randomCount}
                onChange={(e) => setRandomCount(Number(e.target.value))}
                className="border px-3 py-2 rounded-md"
              />

              <button
                onClick={generateRandomDistribution}
                className="bg-purple-600 hover:bg-purple-800 text-white px-4 py-2 rounded-md mt-2"
              >
                Gerar Distribuição
              </button>
            </div>

            {distribution.length > 0 && (
              <div className="max-h-48 overflow-y-auto border-t border-gray-200 pt-3">
                <table className="min-w-full text-sm">
                  <thead>
                    <tr className="text-gray-600 border-b">
                      <th className="text-left px-2 py-1">Agente</th>
                      <th className="text-right px-2 py-1">Valor</th>
                    </tr>
                  </thead>
                  <tbody>
                    {distribution.map((d) => (
                      <tr key={d.id} className="border-b last:border-none">
                        <td className="px-2 py-1 text-gray-700">{d.name}</td>
                        <td className="px-2 py-1 text-right font-medium">
                          {d.amount.toFixed(2)}
                        </td>
                      </tr>
                    ))}
                  </tbody>
                </table>
              </div>
            )}

            <div className="flex justify-end mt-6 gap-3">
              <button
                onClick={() => setShowRandomModal(false)}
                className="px-4 py-2 rounded-md border text-gray-600 hover:bg-gray-100"
              >
                Cancelar
              </button>
              <button
                onClick={confirmRandomDistribution}
                disabled={distributing || distribution.length === 0}
                className={`px-4 py-2 rounded-md text-white ${
                  distributing
                    ? "bg-gray-400 cursor-not-allowed"
                    : "bg-purple-700 hover:bg-purple-900"
                }`}
              >
                {distributing ? "Processando..." : "Confirmar Distribuição"}
              </button>
            </div>
          </div>
        </div>
      )}
    </div>
  );
}
