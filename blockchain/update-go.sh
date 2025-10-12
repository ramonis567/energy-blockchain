#!/bin/bash
set -e

GO_VERSION="1.21.5"
GO_TAR="go${GO_VERSION}.linux-amd64.tar.gz"
GO_URL="https://go.dev/dl/${GO_TAR}"

echo "=========================================="
echo " 🚀 Atualizando Go para versão ${GO_VERSION}"
echo "=========================================="

# 1. Remover versões antigas
echo "🧹 Removendo versões antigas..."
sudo rm -rf /usr/local/go
sudo apt remove -y golang-go || true

# 2. Baixar nova versão
echo "⬇️  Baixando Go ${GO_VERSION}..."
wget -q ${GO_URL} -O /tmp/${GO_TAR}

# 3. Instalar
echo "📦 Instalando Go em /usr/local..."
sudo tar -C /usr/local -xzf /tmp/${GO_TAR}

# 4. Atualizar variáveis de ambiente
echo "⚙️  Atualizando PATH e GOPATH..."
if ! grep -q "/usr/local/go/bin" ~/.bashrc; then
    echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
fi

if ! grep -q "export GOPATH" ~/.bashrc; then
    echo 'export GOPATH=$HOME/go' >> ~/.bashrc
    echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
fi

# 5. Recarregar ambiente
source ~/.bashrc

# 6. Verificar versão
echo "✅ Verificando instalação..."
go version || { echo "❌ Falha ao instalar Go"; exit 1; }

echo "=========================================="
echo " 🎉 Go ${GO_VERSION} instalado com sucesso!"
echo " Agora você pode compilar seus chaincodes."
echo "=========================================="
