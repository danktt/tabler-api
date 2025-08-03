# 🚀 Como Executar a Migração no Neon

## 📋 Passo a Passo

### 1. Acesse o Console do Neon
- Vá para [console.neon.tech](https://console.neon.tech)
- Faça login na sua conta
- Selecione seu projeto

### 2. Abra o SQL Editor
- No menu lateral, clique em **"SQL Editor"**
- Ou clique no botão **"New Query"**

### 3. Execute o Script
Copie e cole este código no SQL Editor:

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

### 4. Execute a Query
- Clique no botão **"Run"** 
- Ou pressione **`Ctrl+Enter`** (Windows/Linux) ou **`Cmd+Enter`** (Mac)

### 5. Verifique se Funcionou
Você deve ver uma mensagem de sucesso. Para confirmar, execute:

```sql
-- Verificar se a tabela foi criada
SELECT table_name 
FROM information_schema.tables 
WHERE table_schema = 'public' 
AND table_name = 'users';
```

## ✅ Próximos Passos

Após executar a migração:

1. **Execute a aplicação**:
```bash
make run
```

2. **Teste a API**:
```bash
curl http://localhost:8080/health
```

3. **Crie um usuário**:
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{"name": "João Silva", "email": "joao@example.com"}'
```

## 🔍 Troubleshooting

### Erro: "relation 'users' does not exist"
- Verifique se executou o script completo
- Confirme se está no banco correto

### Erro de Permissão
- Verifique se você tem permissão para criar tabelas
- Use o console web do Neon (mais seguro)

### Erro de Conexão
- Verifique se o banco Neon está ativo
- Confirme se a DATABASE_URL está correta

---

**🎉 Pronto!** Agora você pode usar a API com o banco Neon configurado! 