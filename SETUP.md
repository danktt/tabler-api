# 🚀 Setup para Desenvolvedores

## 📋 Pré-requisitos

- Go 1.21+
- Git
- Acesso ao banco Neon PostgreSQL

## 🔧 Configuração Inicial

### 1. Clone o repositório
```bash
git clone <repository-url>
cd tabler-api
```

### 2. Configure o ambiente
```bash
# Opção 1: Configuração automática
./setup-env.sh

# Opção 2: Configuração manual
cp env.example .env
# Edite o arquivo .env com suas credenciais
```

### 3. Instale as dependências
```bash
go mod tidy
```

### 4. Execute a migração no Neon
Veja o guia completo em [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md)

**Resumo rápido:**
1. Acesse [console.neon.tech](https://console.neon.tech)
2. Vá para SQL Editor
3. Execute o script em `migrations/001_create_users_table.sql`

### 5. Execute a aplicação
```bash
make run
# ou
go run cmd/api/main.go
```

### 6. Teste a API
```bash
curl http://localhost:8080/health
```

## 🔒 Segurança

### Arquivos Sensíveis
- ✅ `.env` - **NÃO** está no repositório (contém credenciais)
- ✅ `env.example` - Está no repositório (template)
- ✅ `.gitignore` - Configurado para excluir arquivos sensíveis

### Configuração de Ambiente
O projeto suporta duas formas de configuração:

1. **DATABASE_URL** (recomendado para Neon):
```env
DATABASE_URL='postgresql://username:password@host:port/database?sslmode=require&channel_binding=require'
```

2. **Variáveis individuais**:
```env
DB_HOST=your-host
DB_PORT=5432
DB_USER=your-user
DB_PASSWORD=your-password
DB_NAME=your-database
DB_SSL_MODE=require
```

## 🛠️ Comandos Úteis

```bash
make run        # Executar aplicação
make build      # Compilar
make test       # Executar testes
make deps       # Instalar dependências
make help       # Ver todos os comandos
```

## 📁 Estrutura do Projeto

```
tabler-api/
├── cmd/api/main.go              # Entry point
├── config/                      # Configurações
├── internal/                    # Código interno
│   ├── handler/                 # Handlers HTTP
│   ├── service/                 # Regras de negócio
│   ├── repository/              # Operações com banco
│   ├── model/                   # Modelos e DTOs
│   └── router/                  # Rotas
├── pkg/db/                      # Conexão com banco
├── migrations/                  # Scripts SQL
├── examples/                    # Exemplos de uso
├── .env                         # Variáveis de ambiente (local)
├── env.example                  # Template de variáveis
├── setup-env.sh                 # Script de configuração
└── README.md                    # Documentação principal
```

## 🔍 Troubleshooting

### Erro: "module not found"
```bash
go mod tidy
```

### Erro: "connection refused"
- Verifique se o banco Neon está ativo
- Confirme se a DATABASE_URL está correta

### Erro: "table does not exist"
- Execute a migração no console do Neon
- Verifique se está no banco correto

### Erro: "permission denied"
- Verifique se o arquivo `.env` existe
- Confirme se as credenciais estão corretas

## 📚 Documentação Adicional

- [README.md](./README.md) - Documentação principal
- [MIGRATION_GUIDE.md](./MIGRATION_GUIDE.md) - Guia de migrações
- [NEON_SETUP.md](./NEON_SETUP.md) - Configuração específica do Neon

## 🤝 Contribuindo

1. **Nunca commite** arquivos `.env` com credenciais
2. Use sempre `env.example` como template
3. Teste suas mudanças antes de commitar
4. Siga as convenções de código do projeto

---

**🎉 Pronto!** Agora você pode desenvolver com segurança! 