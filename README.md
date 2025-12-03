# API Produtos - POC Fiber + GORM

API RESTful de produtos desenvolvida em Go utilizando **Fiber** como framework web e **GORM** como ORM, demonstrando a viabilidade e facilidade de uso desses frameworks.

## 🚀 Tecnologias Utilizadas

- **Go 1.21** - Linguagem de programação
- **Fiber v2** - Framework web extremamente rápido
- **GORM** - ORM para Go
- **PostgreSQL** - Banco de dados relacional
- **Logrus** - Logging estruturado
- **Validator v10** - Validação de dados
- **Swagger** - Documentação da API
- **Docker** - Containerização

## 📁 Estrutura do Projeto

```
api_fibergorm/
├── cmd/
│   └── api/
│       └── main.go          # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go        # Configurações e logger
│   ├── database/
│   │   └── database.go      # Conexão com banco de dados
│   ├── dto/
│   │   └── produto_dto.go   # Data Transfer Objects
│   ├── handler/
│   │   └── produto_handler.go # Controllers/Handlers HTTP
│   ├── middleware/
│   │   └── middleware.go    # Middlewares da aplicação
│   ├── models/
│   │   └── produto.go       # Entidades/Models
│   ├── repository/
│   │   └── produto_repository.go # Camada de acesso a dados
│   ├── routes/
│   │   └── routes.go        # Configuração de rotas
│   ├── service/
│   │   └── produto_service.go # Regras de negócio
│   └── validator/
│       └── validator.go     # Validador customizado
├── docs/
│   ├── docs.go              # Documentação Swagger
│   └── swagger.json         # Especificação OpenAPI
├── docker-compose.yml       # Orquestração de containers
├── Dockerfile               # Build da aplicação
├── go.mod                   # Dependências Go
└── README.md
```

## 🏃‍♂️ Como Executar

### Pré-requisitos
- Go 1.21+
- Docker e Docker Compose (opcional)
- PostgreSQL (se não usar Docker)

### Com Docker (Recomendado)

```bash
# Inicia todos os serviços
docker-compose up -d

# Verifica os logs
docker-compose logs -f api
```

### Sem Docker

1. Configure o PostgreSQL e crie o banco de dados `produtos_db`

2. Configure as variáveis de ambiente (ou use os valores padrão):
```bash
export DB_HOST=localhost
export DB_PORT=5432
export DB_USER=postgres
export DB_PASSWORD=postgres
export DB_NAME=produtos_db
export SERVER_PORT=3000
```

3. Execute a aplicação:
```bash
go mod download
go run cmd/api/main.go
```

## 📚 Endpoints da API

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check |
| GET | `/swagger/*` | Documentação Swagger |
| POST | `/api/v1/produtos` | Criar produto |
| GET | `/api/v1/produtos` | Listar produtos (paginado) |
| GET | `/api/v1/produtos/:id` | Buscar produto por ID |
| PUT | `/api/v1/produtos/:id` | Atualizar produto |
| DELETE | `/api/v1/produtos/:id` | Excluir produto |

## 📖 Documentação Swagger

Acesse a documentação interativa em: `http://localhost:3000/swagger/`

## 🔍 Exemplos de Requisições

### Criar Produto
```bash
curl -X POST http://localhost:3000/api/v1/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "codigo": "PROD001",
    "descricao": "Notebook Dell Inspiron",
    "preco": 3599.90
  }'
```

### Listar Produtos
```bash
curl http://localhost:3000/api/v1/produtos?page=1&page_size=10
```

### Buscar por ID
```bash
curl http://localhost:3000/api/v1/produtos/1
```

### Atualizar Produto
```bash
curl -X PUT http://localhost:3000/api/v1/produtos/1 \
  -H "Content-Type: application/json" \
  -d '{
    "preco": 3299.90
  }'
```

### Excluir Produto
```bash
curl -X DELETE http://localhost:3000/api/v1/produtos/1
```

## ✅ Validações de Negócio

A API implementa as seguintes validações:

- **Código único**: Não permite duplicidade de códigos
- **Código obrigatório**: Campo código é obrigatório
- **Descrição mínima**: Mínimo de 3 caracteres
- **Preço positivo**: Preço deve ser maior que zero

## 🔒 Validações de Entrada (validator/v10)

- `codigo`: obrigatório, 1-50 caracteres
- `descricao`: obrigatório, 3-255 caracteres
- `preco`: obrigatório, maior que 0

## 📝 Logs

Os logs são estruturados em formato JSON usando Logrus:

```json
{
  "level": "info",
  "msg": "Requisição HTTP",
  "method": "POST",
  "path": "/api/v1/produtos",
  "status": 201,
  "latency": "5.123ms"
}
```

## 🏗️ Arquitetura em Camadas

1. **Handler/Controller**: Recebe requisições HTTP, valida entrada e retorna respostas
2. **Service**: Contém a lógica de negócio e validações
3. **Repository**: Abstrai o acesso ao banco de dados
4. **Model**: Representa as entidades do domínio
5. **DTO**: Objetos de transferência de dados entre camadas

## 🧪 Testando a POC

1. Inicie os containers: `docker-compose up -d`
2. Acesse o Swagger: `http://localhost:3000/swagger/`
3. Teste os endpoints através da interface Swagger ou curl

## 📈 Benefícios Demonstrados

- **Fiber**: Alta performance, sintaxe familiar (Express-like), excelente documentação
- **GORM**: ORM maduro, migrations automáticas, suporte a relacionamentos
- **Arquitetura limpa**: Fácil manutenção e escalabilidade
- **Logs estruturados**: Facilita debugging e monitoramento
- **Swagger**: Documentação automática e interativa

