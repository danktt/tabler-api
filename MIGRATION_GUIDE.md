# 🚀 Guia de Migrações - Neon PostgreSQL

## 📋 **Opção 1: Console Web do Neon (Recomendado)**

### Passo a Passo:

1. **Acesse o Console do Neon**:
   - Vá para [console.neon.tech](https://console.neon.tech)
   - Faça login na sua conta
   - Selecione seu projeto `neondb`

2. **Abra o SQL Editor**:
   - No menu lateral, clique em "SQL Editor"
   - Ou clique no botão "New Query"

3. **Execute a Migração**:
   - Copie e cole o script abaixo
   - Clique em "Run" ou pressione `Ctrl+Enter`

### Script de Migração:

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

## 🔧 **Opção 2: Instalar psql (macOS)**

Se você quiser usar linha de comando:

### Instalar PostgreSQL (inclui psql):
```bash
# Usando Homebrew
brew install postgresql

# Ou baixar do site oficial
# https://www.postgresql.org/download/macosx/
```

### Executar migração:
```bash
# Tornar o script executável
chmod +x run-migration.sh

# Executar migração
./run-migration.sh
```

## 🐳 **Opção 3: Usando Docker**

Se você tiver Docker instalado:

```bash
# Executar psql via Docker
docker run --rm -it postgres:15 psql "postgresql://neondb_owner:npg_rgLTU4Haxzh3@ep-calm-waterfall-acb0d7qm-pooler.sa-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require" -f migrations/001_create_users_table.sql
```

## ✅ **Verificar se a Migração Funcionou**

### No Console do Neon:
1. Vá para "Tables" no menu lateral
2. Você deve ver a tabela `users` listada

### Ou execute esta query no SQL Editor:
```sql
-- Verificar se a tabela foi criada
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'users';

-- Verificar estrutura da tabela
\d users;
```

## 🚀 **Após a Migração**

1. **Execute a aplicação**:
```bash
make run
# ou
go run cmd/api/main.go
```

2. **Teste a API**:
```bash
# Health check
curl http://localhost:8080/health

# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "João Silva", "email": "joao@example.com"}'
```

## 🔍 **Troubleshooting**

### Erro: "relation 'users' does not exist"
- Verifique se a migração foi executada
- Confirme se você está no banco correto

### Erro de Conexão
- Verifique se a URL do Neon está correta
- Confirme se o banco está ativo

### Erro de Permissão
- Verifique se o usuário tem permissão para criar tabelas
- Use o console web do Neon (mais seguro)

## 📊 **Próximas Migrações**

Para futuras migrações, crie novos arquivos:
```
migrations/
├── 001_create_users_table.sql
├── 002_add_user_roles.sql
└── 003_create_posts_table.sql
```

E execute no console do Neon ou usando os scripts fornecidos.

---

**💡 Dica**: Use sempre o console web do Neon para migrações - é mais seguro e confiável! 