# Configuração com Neon PostgreSQL

Este guia mostra como configurar o projeto para usar o Neon PostgreSQL.

## 🚀 Configuração Rápida

### 1. Ambiente já configurado
O arquivo `.env` já foi criado com sua URL do Neon:
```env
DATABASE_URL='postgresql://neondb_owner:npg_rgLTU4Haxzh3@ep-calm-waterfall-acb0d7qm-pooler.sa-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require'
```

### 2. Executar migração no Neon
1. Acesse o [console do Neon](https://console.neon.tech)
2. Vá para o seu projeto `neondb`
3. Clique em "SQL Editor"
4. Execute o seguinte script:

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

### 3. Executar a aplicação
```bash
make run
# ou
go run cmd/api/main.go
```

### 4. Testar a API
```bash
# Health check
curl http://localhost:8080/health

# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "name": "João Silva",
    "email": "joao@example.com"
  }'
```

## 🔧 Detalhes da Configuração

### URL de Conexão
A URL do Neon inclui todos os parâmetros necessários:
- `sslmode=require`: SSL obrigatório (requerido pelo Neon)
- `channel_binding=require`: Autenticação segura
- Pooler: Usa o pooler do Neon para melhor performance

### Configuração Automática
O projeto foi configurado para:
1. **Detectar automaticamente** se `DATABASE_URL` está definida
2. **Usar a URL completa** do Neon quando disponível
3. **Fallback** para variáveis individuais se necessário

### Pool de Conexões
O projeto usa `pgxpool` com configurações otimizadas:
- Máximo de 10 conexões
- Mínimo de 2 conexões
- Timeout de 30 segundos para conexões ociosas

## 🧪 Testando a Conexão

### 1. Verificar logs de conexão
Quando você executar a aplicação, deve ver:
```
Successfully connected to PostgreSQL database
Starting server on localhost:8080
```

### 2. Testar endpoints
```bash
# Health check
curl http://localhost:8080/health

# Listar usuários (deve retornar array vazio inicialmente)
curl http://localhost:8080/api/v1/users

# Criar usuário
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "Teste", "email": "teste@example.com"}'
```

## 🔍 Troubleshooting

### Erro de Conexão
Se você ver erro de conexão:
1. Verifique se o banco Neon está ativo
2. Confirme se a URL está correta
3. Verifique se a migração foi executada

### Erro de SSL
O Neon requer SSL. A configuração já inclui:
- `sslmode=require`
- `channel_binding=require`

### Erro de Tabela
Se você ver erro de tabela não encontrada:
1. Execute a migração SQL no console do Neon
2. Verifique se a tabela `users` foi criada

## 📊 Monitoramento

### Console do Neon
- Acesse o console para ver métricas de conexão
- Monitore o uso de recursos
- Verifique logs de queries

### Logs da Aplicação
A aplicação loga:
- Conexão bem-sucedida com o banco
- Início do servidor
- Erros de conexão (se houver)

## 🚀 Próximos Passos

1. **Execute a migração** no console do Neon
2. **Inicie a aplicação**: `make run`
3. **Teste os endpoints** usando os exemplos acima
4. **Explore a API** usando o arquivo `examples/api_test.http`

A API está pronta para uso com o Neon! 🎉 