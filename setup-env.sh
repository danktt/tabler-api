#!/bin/bash

# Script para configurar o ambiente de desenvolvimento

echo "Configurando ambiente de desenvolvimento..."

# Criar arquivo .env se não existir
if [ ! -f .env ]; then
    echo "Criando arquivo .env..."
    cat > .env << EOF
# Database Configuration
DATABASE_URL='postgresql://neondb_owner:npg_rgLTU4Haxzh3@ep-calm-waterfall-acb0d7qm-pooler.sa-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require'

# Server Configuration
SERVER_PORT=8080
SERVER_HOST=localhost

# Environment
ENV=development
EOF
    echo "✅ Arquivo .env criado com sucesso!"
else
    echo "⚠️  Arquivo .env já existe. Verifique se a DATABASE_URL está configurada corretamente."
fi

echo ""
echo "📋 Próximos passos:"
echo "1. Execute a migração SQL no seu banco Neon:"
echo "   - Acesse o console do Neon"
echo "   - Execute o script em migrations/001_create_users_table.sql"
echo ""
echo "2. Execute a aplicação:"
echo "   make run"
echo "   ou"
echo "   go run cmd/api/main.go"
echo ""
echo "3. Teste a API:"
echo "   curl http://localhost:8080/health" 