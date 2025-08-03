#!/bin/bash

# Script para executar migração no Neon via psql

echo "🚀 Executando migração no Neon..."

# Verificar se psql está instalado
if ! command -v psql &> /dev/null; then
    echo "❌ psql não está instalado. Use a opção do console web."
    echo "📋 Acesse: https://console.neon.tech"
    echo "📋 Execute o script em migrations/001_create_users_table.sql"
    exit 1
fi

# URL do banco (do arquivo .env)
DATABASE_URL="postgresql://neondb_owner:npg_rgLTU4Haxzh3@ep-calm-waterfall-acb0d7qm-pooler.sa-east-1.aws.neon.tech/neondb?sslmode=require&channel_binding=require"

echo "📊 Conectando ao banco Neon..."
psql "$DATABASE_URL" -f migrations/001_create_users_table.sql

if [ $? -eq 0 ]; then
    echo "✅ Migração executada com sucesso!"
    echo "🎉 A tabela 'users' foi criada no Neon"
else
    echo "❌ Erro ao executar migração"
    echo "📋 Use a opção do console web: https://console.neon.tech"
fi 