# Tabler API

Uma API REST em Go usando Gin e Neon PostgreSQL com arquitetura limpa.

## Estrutura do Projeto

```
tabler-api/
├── cmd/api/main.go              # Entry point da aplicação
├── config/                      # Configurações e variáveis de ambiente
├── internal/
│   ├── handler/                 # Handlers HTTP
│   ├── service/                 # Regras de negócio
│   ├── repository/              # Operações com banco de dados
│   ├── model/                   # Structs de entidades e DTOs
│   └── router/                  # Definição das rotas
├── pkg/db/                      # Conexão com banco de dados
├── migrations/                  # Scripts SQL
└── README.md
```

## Pré-requisitos

- Go 1.21+
- Conta no Neon PostgreSQL
- Git

## Configuração Rápida

1. **Clone o repositório**:
```bash
git clone <repository-url>
cd tabler-api
```

2. **Configure o ambiente**:
```bash
cp env.example .env
# Edite o .env com sua DATABASE_URL do Neon
```

3. **Instale as dependências**:
```bash
go mod tidy
```

4. **Execute a migração no Neon**:
   - Acesse o [console do Neon](https://console.neon.tech)
   - Vá para o seu projeto
   - Execute o script SQL em `migrations/001_create_users_table.sql`

5. **Execute a aplicação**:
```bash
make run
# ou
go run cmd/api/main.go
```

## Configuração do Neon

### 1. Obter DATABASE_URL
1. Acesse [console.neon.tech](https://console.neon.tech)
2. Crie um novo projeto ou use um existente
3. Copie a URL de conexão fornecida

### 2. Configurar .env
```env
DATABASE_URL='postgresql://username:password@host:port/database?sslmode=require&channel_binding=require'
SERVER_PORT=8080
SERVER_HOST=localhost
ENV=development
```

### 3. Executar Migração
No console SQL do Neon, execute:
```sql
-- Create users table
CREATE TABLE IF NOT EXISTS users (
    id UUID PRIMARY KEY DEFAULT gen_random_uuid(),
    name VARCHAR(255) NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE DEFAULT CURRENT_TIMESTAMP
);

-- Create index on email for faster lookups
CREATE INDEX IF NOT EXISTS idx_users_email ON users(email);

-- Create index on created_at for sorting
CREATE INDEX IF NOT EXISTS idx_users_created_at ON users(created_at DESC);
```

## Executando a Aplicação

```bash
go run cmd/api/main.go
```

A API estará disponível em `http://localhost:8080`

## Autenticação

Este projeto implementa autenticação JWT usando Auth0. Para mais detalhes, consulte o arquivo [AUTHENTICATION.md](AUTHENTICATION.md).

### Configuração Rápida do Auth0

1. Configure as variáveis de ambiente no `.env`:
```env
AUTH0_DOMAIN=your-tenant.auth0.com
AUTH0_AUDIENCE=your-api-identifier
```

2. Teste a autenticação:
```bash
./scripts/test_auth.sh
```

## Endpoints da API

### Health Check
- `GET /health` - Verifica se o servidor está funcionando

### Rotas Públicas (Usuários)

#### Criar usuário
```bash
POST /api/v1/users
Content-Type: application/json

{
  "name": "João Silva",
  "email": "joao@example.com"
}
```

#### Listar todos os usuários
```bash
GET /api/v1/users
```

#### Buscar usuário por ID
```bash
GET /api/v1/users/{id}
```

#### Atualizar usuário
```bash
PUT /api/v1/users/{id}
Content-Type: application/json

{
  "name": "João Silva Atualizado",
  "email": "joao.novo@example.com"
}
```

#### Deletar usuário
```bash
DELETE /api/v1/users/{id}
```

### Rotas Protegidas (Requerem Autenticação)

#### Obter perfil do usuário autenticado
```bash
GET /api/v1/profile
Authorization: Bearer YOUR_JWT_TOKEN
```

**Resposta:**
```json
{
  "message": "Profile retrieved successfully",
  "data": {
    "sub": "auth0|1234567890",
    "email": "user@example.com",
    "name": "John Doe",
    "nickname": "johndoe",
    "picture": "https://example.com/picture.jpg",
    "updated_at": "2023-01-01T00:00:00.000Z",
    "issuer": "https://your-tenant.auth0.com/",
    "audience": "your-api-identifier",
    "expires_at": 1634571490,
    "issued_at": 1634567890
  }
}
```

## Exemplos de Uso

### Criar um usuário
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Maria Santos",
    "email": "maria@example.com"
  }'
```

### Listar usuários
```bash
curl http://localhost:8080/api/v1/users
```

### Buscar usuário específico
```bash
curl http://localhost:8080/api/v1/users/{user-id}
```

## Arquitetura

O projeto segue os princípios da Clean Architecture:

- **Handlers**: Responsáveis por receber requisições HTTP e retornar respostas
- **Services**: Contêm a lógica de negócio
- **Repositories**: Responsáveis pelas operações com banco de dados
- **Models**: Definem as estruturas de dados e DTOs

## Tecnologias Utilizadas

- **Gin**: Framework web para Go
- **pgx**: Driver PostgreSQL para Go
- **UUID**: Geração de IDs únicos
- **godotenv**: Carregamento de variáveis de ambiente
- **Neon**: PostgreSQL como serviço
- **Auth0**: Provedor de identidade e autenticação JWT
- **lestrrat-go/jwx**: Biblioteca para validação de tokens JWT e JWK

## Comandos Úteis

```bash
make run        # Executar a aplicação
make build      # Compilar a aplicação
make test       # Executar testes
make deps       # Instalar dependências
make help       # Ver todos os comandos
```

## Desenvolvimento

Para executar em modo de desenvolvimento:

```bash
ENV=development go run cmd/api/main.go
```

## Build

Para criar um executável:

```bash
go build -o bin/api cmd/api/main.go
```

## Testes

Para executar os testes:

```bash
go test ./...
```

## Troubleshooting

### Erro: "DATABASE_URL is required"
- Verifique se o arquivo `.env` existe
- Confirme se a `DATABASE_URL` está configurada corretamente

### Erro de Conexão com o Neon
- Verifique se a `DATABASE_URL` está correta
- Confirme se o banco Neon está ativo
- Verifique se a migração foi executada

### Erro: "relation 'users' does not exist"
- Execute a migração no console do Neon
- Verifique se está no banco correto

### Erro de Compilação
```bash
go mod tidy
go build cmd/api/main.go
```

### Logs de Debug
Para ver logs detalhados, execute:
```bash
ENV=development go run cmd/api/main.go
``` 