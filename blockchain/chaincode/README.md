Verificar caminho dos chaincodes

echo 'export PATH=$PATH:/usr/local/go/bin' >> ~/.bashrc
echo 'export GOPATH=$HOME/go' >> ~/.bashrc
echo 'export PATH=$PATH:$GOPATH/bin' >> ~/.bashrc
source ~/.bashrc

Iniciar dependências e gerar .sum
go mod init creditmarket
go mod tidy
