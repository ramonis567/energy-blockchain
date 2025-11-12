import { useEffect, useState } from "react";
import { getAgents, registerAgent, mintTokens } from "../api/api";
import { nanoid } from "nanoid";

interface Agent {
  id: string;
  type: string;
  name: string;
  address: string;
  ecr_balance: number;
  engt_balance: number;
  registered_at: string;
}

export default function Agents() {
  const [agents, setAgents] = useState<Agent[]>([]);
  const [filteredAgents, setFilteredAgents] = useState<Agent[]>([]);
  const [loading, setLoading] = useState(true);
  const [error, setError] = useState<string | null>(null);
  const [typeFilter, setTypeFilter] = useState("");
  const [nameFilter, setNameFilter] = useState("");

  // === Seleção múltipla ===
  const [selectedAgents, setSelectedAgents] = useState<string[]>([]);

  // === Modais ===
  const [showCreateModal, setShowCreateModal] = useState(false);
  const [showMintModal, setShowMintModal] = useState(false);

  // Cadastro de agentes
  const [formType, setFormType] = useState("consumer");
  const [formBaseName, setFormBaseName] = useState("");
  const [formAddress, setFormAddress] = useState("");
  const [formCount, setFormCount] = useState(1);
  const [saving, setSaving] = useState(false);

  // Mint de tokens
  const [mintToken, setMintToken] = useState("ECR");
  const [mintAmount, setMintAmount] = useState("");
  const [minting, setMinting] = useState(false);

  // Refresh controlador
  const [refreshing, setRefreshing] = useState(false);

  // === Carregar agentes ===
  async function fetchAgents() {
	try {
	  setLoading(true);
	  const data = await getAgents();
	  setAgents(data);
	  setFilteredAgents(data);
	  setSelectedAgents([]);
	} catch (err) {
	  console.error(err);
	  setError("Erro ao carregar agentes");
	} finally {
	  setLoading(false);
	}
  }

  useEffect(() => {
	fetchAgents();
  }, []);

  // === Aplicar filtros ===
  useEffect(() => {
	let result = [...agents];
	if (typeFilter) result = result.filter((a) => a.type === typeFilter);
	if (nameFilter)
	  result = result.filter((a) =>
		a.id.toLowerCase().includes(nameFilter.toLowerCase())
	  );
	setFilteredAgents(result);
  }, [typeFilter, nameFilter, agents]);

  // Atualizar dados após clique no botão
  async function handleRefresh() {
	setRefreshing(true);
	try {
		await fetchAgents();
	} catch (err) {
		console.error("Erro ao atualizar tabela:", err);
	} finally {
		setRefreshing(false);
	}
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

  // === Cadastro em lote ===
	async function handleRegisterBatch() {
		if (!formBaseName.trim() || !formAddress.trim()) {
			alert("Preencha todos os campos do formulário.");
			return;
		}

		setSaving(true);
		try {
			const createdAgents = [];

			// Define o prefixo conforme tipo
			const typePrefixMap: Record<string, string> = {
			producer: "prod",
			consumer: "cons",
			prosumer: "pros",
			distributor: "dist",
			battery: "bat",
			};

			const prefix = typePrefixMap[formType] || "agt";

			for (let i = 1; i <= formCount; i++) {
			const id = `${prefix}-${nanoid(7)}`; // padroniza o ID
			const agent = {
				id,
				type: formType,
				name: formCount > 1 ? `${formBaseName} ${i}` : formBaseName,
				address: formAddress,
			};
			const res = await registerAgent(agent);
			if (res?.status === "ok") createdAgents.push(agent);
			}

			alert(`${createdAgents.length} agente(s) criado(s) com sucesso!`);
			setShowCreateModal(false);
			setFormBaseName("");
			setFormAddress("");
			setFormCount(1);
			fetchAgents();
		} catch (err) {
			console.error(err);
			alert("Erro ao registrar agentes.");
		} finally {
			setSaving(false);
		}
	}


  // === Mint tokens ===
  async function handleMintTokens() {
	if (selectedAgents.length === 0) {
	  alert("Selecione ao menos um agente.");
	  return;
	}
	if (!mintToken || !mintAmount) {
	  alert("Preencha todos os campos para adicionar saldo.");
	  return;
	}

	setMinting(true);
	try {
	  for (const id of selectedAgents) {
		console.log(`Minting ${mintAmount} ${mintToken} → ${id}`);
		const res = await mintTokens(id, mintToken, mintAmount);
		if (res?.status !== "ok") console.warn("Erro em:", id, res);
	  }

	  alert(`✅ ${mintAmount} ${mintToken} emitidos para ${selectedAgents.length} agentes.`);
	  setShowMintModal(false);
	  setSelectedAgents([]);
	  setMintAmount("");
	  fetchAgents();
	} catch (err) {
	  console.error(err);
	  alert("Erro ao realizar mint.");
	} finally {
	  setMinting(false);
	}
  }

  // === Render ===
  if (loading) return <p className="text-gray-500">Carregando agentes...</p>;
  if (error) return <p className="text-red-600">{error}</p>;

  return (
	<div className="p-6 bg-gray-50 rounded-lg shadow-sm">
	  <div className="flex flex-wrap items-center justify-between mb-6">
		<h2 className="text-2xl font-semibold text-[var(--blueColor)]">
		  Agentes Registrados
		</h2>
		<span className="text-lg font-medium text-gray-700">
		  Total:{" "}
		  <span className="text-[var(--blueColor)] font-semibold">
			{filteredAgents.length}
		  </span>
		</span>
	  </div>

	  {/* Filtros */}
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

		<input
		  type="text"
		  placeholder="Buscar por ID..."
		  value={nameFilter}
		  onChange={(e) => setNameFilter(e.target.value)}
		  className="px-3 py-2 border rounded-md focus:ring-2 focus:ring-[var(--blueColor)] w-64"
		/>

		<div className="flex gap-3 ml-auto">
			{/* Botão de Refresh */}
			<button
				onClick={handleRefresh}
				disabled={refreshing}
				className={`flex items-center gap-2 px-4 py-2 rounded-md font-medium text-white transition ${
				refreshing
					? "bg-gray-400 cursor-not-allowed"
					: "bg-[var(--highlightColor)] hover:bg-sky-600"
				}`}
			>
				{refreshing ? (
				<>
					<svg
					className="animate-spin h-4 w-4 text-white"
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					>
					<circle
						className="opacity-25"
						cx="12"
						cy="12"
						r="10"
						stroke="currentColor"
						strokeWidth="4"
					></circle>
					<path
						className="opacity-75"
						fill="currentColor"
						d="M4 12a8 8 0 018-8v8H4z"
					></path>
					</svg>
					Atualizando...
				</>
				) : (
				<>
					<svg
					xmlns="http://www.w3.org/2000/svg"
					fill="none"
					viewBox="0 0 24 24"
					strokeWidth={2}
					stroke="currentColor"
					className="h-4 w-4"
					>
					<path
						strokeLinecap="round"
						strokeLinejoin="round"
						d="M4 4v6h6M20 20v-6h-6m6-4a8 8 0 00-15.9-1M4 12a8 8 0 0015.9 1"
					/>
					</svg>
					Atualizar
				</>
				)}
			</button>

			<button
				onClick={() => setShowCreateModal(true)}
				className="bg-[var(--blueColor)] hover:bg-blue-900 text-white font-medium px-4 py-2 rounded-md transition"
			>
				+ Novo Agente
			</button>

			<button
				onClick={() => setShowMintModal(true)}
				disabled={selectedAgents.length === 0}
				className={`font-medium px-4 py-2 rounded-md transition text-white ${
				selectedAgents.length === 0
					? "bg-gray-400 cursor-not-allowed"
					: "bg-[var(--greenColor)] hover:bg-green-700"
				}`}
			>
				Adicionar Saldo ({selectedAgents.length})
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
			  <th className="px-4 py-2">Tipo</th>
			  <th className="px-4 py-2">Nome</th>
			  <th className="px-4 py-2">Endereço</th>
			  <th className="px-4 py-2 text-right">ECR</th>
			  <th className="px-4 py-2 text-right">ENGT</th>
			  <th className="px-4 py-2 text-right">Registrado em</th>
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
				<td className="px-4 py-2 capitalize">{a.type}</td>
				<td className="px-4 py-2">{a.name}</td>
				<td className="px-4 py-2">{a.address}</td>
				<td className="px-4 py-2 text-right">{a.ecr_balance.toFixed(2)}</td>
				<td className="px-4 py-2 text-right">{a.engt_balance.toFixed(2)}</td>
				<td className="px-4 py-2 text-right">
				  {new Date(a.registered_at).toLocaleString("pt-BR")}
				</td>
			  </tr>
			))}
		  </tbody>
		</table>
	  </div>

	  {/* === Modal de Mint === */}
	  {showMintModal && (
		<div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
		  <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
			<h3 className="text-xl font-semibold text-[var(--greenColor)] mb-4">
			  Adicionar Saldo
			</h3>

			<p className="text-sm text-gray-600 mb-3">
			  {selectedAgents.length} agente(s) selecionado(s).
			</p>

			<div className="flex flex-col gap-3">
			  <label className="text-sm font-medium text-gray-700">Tipo de Token</label>
			  <select
				value={mintToken}
				onChange={(e) => setMintToken(e.target.value)}
				className="border px-3 py-2 rounded-md focus:ring-2 focus:ring-[var(--greenColor)]"
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
				  minting
					? "bg-gray-400 cursor-not-allowed"
					: "bg-[var(--greenColor)] hover:bg-green-700"
				}`}
			  >
				{minting ? "Processando..." : "Confirmar"}
			  </button>
			</div>
		  </div>
		</div>
	  )}

	  {/* === Modal de Cadastro (igual à versão anterior) === */}
	  {showCreateModal && (
		<div className="fixed inset-0 bg-black bg-opacity-40 flex items-center justify-center z-50">
		  <div className="bg-white p-6 rounded-lg shadow-lg w-full max-w-md">
			<h3 className="text-xl font-semibold text-[var(--blueColor)] mb-4">
			  Cadastrar Agente
			</h3>

			<div className="flex flex-col gap-3">
			  <label className="text-sm font-medium text-gray-700">Tipo</label>
			  <select
				value={formType}
				onChange={(e) => setFormType(e.target.value)}
				className="border px-3 py-2 rounded-md focus:ring-2 focus:ring-[var(--blueColor)]"
			  >
				<option value="producer">Produtor</option>
				<option value="consumer">Consumidor</option>
				<option value="prosumer">Prosumer</option>
				<option value="distributor">Distribuidor</option>
				<option value="battery">Bateria</option>
			  </select>

			  <label className="text-sm font-medium text-gray-700">Nome base</label>
			  <input
				type="text"
				value={formBaseName}
				onChange={(e) => setFormBaseName(e.target.value)}
				placeholder="Ex: Consumidor"
				className="border px-3 py-2 rounded-md"
			  />

			  <label className="text-sm font-medium text-gray-700">Endereço</label>
			  <input
				type="text"
				value={formAddress}
				onChange={(e) => setFormAddress(e.target.value)}
				placeholder="Rua Exemplo 123"
				className="border px-3 py-2 rounded-md"
			  />

			  <label className="text-sm font-medium text-gray-700">Quantidade</label>
			  <input
				type="number"
				min={1}
				max={50}
				value={formCount}
				onChange={(e) => setFormCount(Number(e.target.value))}
				className="border px-3 py-2 rounded-md w-24"
			  />
			</div>

			<div className="flex justify-end mt-6 gap-3">
			  <button
				onClick={() => setShowCreateModal(false)}
				className="px-4 py-2 rounded-md border text-gray-600 hover:bg-gray-100"
			  >
				Cancelar
			  </button>
			  <button
				onClick={handleRegisterBatch}
				disabled={saving}
				className={`px-4 py-2 rounded-md text-white ${
				  saving
					? "bg-gray-400 cursor-not-allowed"
					: "bg-[var(--blueColor)] hover:bg-blue-900"
				}`}
			  >
				{saving ? "Salvando..." : "Salvar"}
			  </button>
			</div>
		  </div>
		</div>
	  )}
	</div>
  );
}
