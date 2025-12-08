# API Produtos - POC Fiber + GORM

API RESTful de produtos desenvolvida em Go utilizando **Fiber** como framework web e **GORM** como ORM, demonstrando a viabilidade e facilidade de uso desses frameworks.

## 🚀 Tecnologias Utilizadas

- **Go 1.23** - Linguagem de programação
- **Fiber v2** - Framework web extremamente rápido
- **GORM** - ORM para Go
- **PostgreSQL** - Banco de dados relacional
- **Prometheus** - Métricas e monitoramento
- **Loki** - Agregação de logs (Grafana)
- **Logrus** - Logging estruturado
- **Validator v10** - Validação de dados
- **Swagger** - Documentação da API
- **Docker** - Containerização

## 📁 Estrutura do Projeto

```
api_fibergorm/
├── cmd/
│   └── api/
│       └── main.go              # Ponto de entrada da aplicação
├── internal/
│   ├── config/
│   │   └── config.go            # Configurações e logger
│   ├── database/
│   │   ├── database.go          # Conexão e migrations
│   │   └── seed.go              # Carga inicial de dados
│   ├── dto/
│   │   ├── categoria_dto.go     # DTOs de Categoria
│   │   └── produto_dto.go       # DTOs de Produto
│   ├── handler/
│   │   ├── categoria_handler.go # Controller de Categorias
│   │   └── produto_handler.go   # Controller de Produtos
│   ├── middleware/
│   │   └── middleware.go        # Middlewares da aplicação
│   ├── models/
│   │   ├── categoria.go         # Entidade Categoria
│   │   └── produto.go           # Entidade Produto
│   ├── repository/
│   │   ├── categoria_repository.go
│   │   └── produto_repository.go
│   ├── routes/
│   │   └── routes.go            # Configuração de rotas
│   ├── service/
│   │   ├── categoria_service.go # Regras de negócio
│   │   └── produto_service.go
│   └── validator/
│       └── validator.go         # Validador customizado
├── docs/                        # Documentação Swagger
├── docker-compose.yml
├── Dockerfile
├── go.mod
└── README.md
```

## ⚙️ Variáveis de Ambiente

Todas as variáveis são **opcionais** e possuem valores padrão:

### Servidor

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `SERVER_PORT` | Porta do servidor HTTP | `3000` |
| `SERVER_READ_TIMEOUT` | Timeout de leitura (segundos) | `10` |
| `SERVER_WRITE_TIMEOUT` | Timeout de escrita (segundos) | `10` |

### Banco de Dados PostgreSQL

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `DB_HOST` | Host do PostgreSQL | `localhost` |
| `DB_PORT` | Porta do PostgreSQL | `5432` |
| `DB_USER` | Usuário do banco | `postgres` |
| `DB_PASSWORD` | Senha do banco | `postgres` |
| `DB_NAME` | Nome do banco de dados | `produtos_db` |
| `DB_SSLMODE` | Modo SSL (disable, require, verify-ca, verify-full) | `disable` |
| `DB_MAX_OPEN_CONNS` | Máximo de conexões abertas | `10` |
| `DB_MAX_IDLE_CONNS` | Máximo de conexões ociosas | `5` |
| `DB_CONN_MAX_LIFETIME` | Tempo de vida da conexão (minutos) | `30` |

### Logging

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `LOG_LEVEL` | Nível de log (debug, info, warn, error) | `debug` |
| `LOG_FORMAT` | Formato do log (json, text) | `json` |

### Loki (Observabilidade)

| Variável | Descrição | Padrão |
|----------|-----------|--------|
| `LOKI_ENABLED` | Habilita integração com Loki | `true` |
| `LOKI_URL` | URL do endpoint push do Loki | `http://10.110.0.239:3100/loki/api/v1/push` |
| `LOKI_SERVICE_NAME` | Nome do job para identificação no Grafana | `ARQUITETURA_FIBER_GORM` |
| `LOKI_BATCH_SIZE` | Quantidade de logs por batch | `10` |
| `LOKI_BATCH_WAIT_SECONDS` | Intervalo máximo entre envios (segundos) | `5` |
| `LOKI_TIMEOUT_SECONDS` | Timeout das requisições HTTP (segundos) | `10` |
| `ENVIRONMENT` | Ambiente da aplicação (label no Loki) | `development` |

## 🏃‍♂️ Como Executar

### Com Docker (Recomendado)

```bash
# Inicia todos os serviços
docker-compose up -d

# Verifica os logs
docker-compose logs -f api
```

### Sem Docker

```bash
# A aplicação cria automaticamente o banco de dados se não existir!
go mod download
go run cmd/api/main.go
```

Ou com variáveis personalizadas:

```bash
DB_HOST=meuhost DB_PASSWORD=minhasenha go run cmd/api/main.go
```

## 📚 Endpoints da API

### Categorias

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/categorias` | Criar categoria |
| GET | `/api/v1/categorias` | Listar categorias (paginado) |
| GET | `/api/v1/categorias/ativas` | Listar apenas ativas |
| GET | `/api/v1/categorias/:id` | Buscar por ID |
| GET | `/api/v1/categorias/:id/produtos` | Categoria com seus produtos |
| PUT | `/api/v1/categorias/:id` | Atualizar categoria |
| DELETE | `/api/v1/categorias/:id` | Excluir categoria |

### Produtos

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| POST | `/api/v1/produtos` | Criar produto |
| GET | `/api/v1/produtos` | Listar produtos (paginado) |
| GET | `/api/v1/produtos/categoria/:id` | Produtos por categoria |
| GET | `/api/v1/produtos/:id` | Buscar por ID |
| PUT | `/api/v1/produtos/:id` | Atualizar produto |
| DELETE | `/api/v1/produtos/:id` | Excluir produto |

### Outros

| Método | Endpoint | Descrição |
|--------|----------|-----------|
| GET | `/health` | Health check |
| GET | `/metrics` | Métricas Prometheus |
| GET | `/swagger/*` | Documentação Swagger |

## 📖 Documentação Swagger

Acesse a documentação interativa em: `http://localhost:3000/swagger/`

## 🔍 Exemplos de Requisições

### Criar Categoria
```bash
curl -X POST http://localhost:3000/api/v1/categorias \
  -H "Content-Type: application/json" \
  -d '{
    "nome": "Eletrônicos",
    "descricao": "Produtos eletrônicos em geral"
  }'
```

### Criar Produto
```bash
curl -X POST http://localhost:3000/api/v1/produtos \
  -H "Content-Type: application/json" \
  -d '{
    "codigo": "PROD001",
    "descricao": "Notebook Dell Inspiron",
    "preco": 3599.90,
    "categoria_id": 1
  }'
```

### Listar Produtos com Categoria
```bash
curl http://localhost:3000/api/v1/produtos?page=1&page_size=10
```

### Buscar Categoria com Produtos
```bash
curl http://localhost:3000/api/v1/categorias/1/produtos
```

## 🔗 Relacionamentos (GORM)

```
┌─────────────┐       ┌─────────────┐
│  Categoria  │ 1───N │   Produto   │
├─────────────┤       ├─────────────┤
│ ID          │       │ ID          │
│ Nome        │       │ Codigo      │
│ Descricao   │       │ Descricao   │
│ Ativo       │       │ Preco       │
│             │       │ CategoriaID │◄── FK obrigatória
└─────────────┘       └─────────────┘
```

Recursos do GORM demonstrados:
- **`foreignKey`** - Define chave estrangeira
- **`Preload()`** - Eager loading de relacionamentos
- **`AutoMigrate`** - Criação automática de tabelas e FKs

## ✅ Validações de Negócio

### Categorias
- Nome único e obrigatório (mín. 2 caracteres)
- Não é possível excluir categoria com produtos

### Produtos
- Código único e obrigatório
- Descrição mínima de 3 caracteres
- Preço deve ser maior que zero
- Categoria obrigatória e deve estar ativa

## 🌱 Seed de Dados

Na primeira execução, a aplicação:
1. Cria o banco de dados automaticamente
2. Executa as migrations (criação de tabelas)
3. Cria uma categoria padrão "Geral"
4. Atualiza produtos órfãos para a categoria padrão

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

## 📊 Métricas Prometheus

A aplicação expõe métricas no endpoint `/metrics` para monitoramento com Prometheus.

### Métricas Disponíveis

| Métrica | Tipo | Descrição |
|---------|------|-----------|
| `http_requests_total` | Counter | Total de requisições HTTP recebidas |
| `http_request_duration_seconds` | Histogram | Duração das requisições HTTP em segundos |
| `http_requests_in_flight` | Gauge | Número de requisições em processamento |
| `http_response_size_bytes` | Histogram | Tamanho das respostas HTTP em bytes |
| `database_queries_total` | Counter | Total de queries executadas no banco |
| `database_query_duration_seconds` | Histogram | Duração das queries em segundos |

### Labels das Métricas HTTP

- `method`: Método HTTP (GET, POST, PUT, DELETE)
- `path`: Padrão da rota (ex: `/api/v1/produtos/:id`)
- `status`: Código de status HTTP (200, 201, 400, 404, 500)

### Exemplo de Consulta PromQL

```promql
# Taxa de requisições por segundo nos últimos 5 minutos
rate(http_requests_total[5m])

# Latência P95 das requisições
histogram_quantile(0.95, rate(http_request_duration_seconds_bucket[5m]))

# Requisições com erro (status >= 400)
sum(rate(http_requests_total{status=~"4..|5.."}[5m]))
```

### Configuração no Prometheus

Adicione ao seu `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'api_fibergorm'
    static_configs:
      - targets: ['localhost:3000']
    metrics_path: /metrics
    scrape_interval: 15s
```

## 📋 Integração com Loki/Grafana

A aplicação envia logs automaticamente para o Loki, permitindo visualização e consulta no Grafana.

### Configuração Padrão

- **URL do Loki**: `http://10.110.0.239:3100/loki/api/v1/push`
- **Job Name**: `ARQUITETURA_FIBER_GORM`

### Labels dos Logs

Cada log enviado ao Loki contém os seguintes labels:

| Label | Descrição |
|-------|-----------|
| `job` | Nome do serviço (`ARQUITETURA_FIBER_GORM`) |
| `level` | Nível do log (debug, info, warn, error) |
| `hostname` | Nome do host da aplicação |
| `app` | Nome da aplicação (`api_fibergorm`) |
| `environment` | Ambiente (development, production) |

### Campos dos Logs

Cada entrada de log contém os seguintes campos JSON:

```json
{
  "level": "info",
  "msg": "Requisição HTTP",
  "time": "2024-01-15T10:30:00.123456789Z",
  "hostname": "api-server-01",
  "service": "ARQUITETURA_FIBER_GORM",
  "method": "GET",
  "path": "/api/v1/produtos",
  "status": 200,
  "latency": "5.123ms"
}
```

### Consultas LogQL no Grafana

```logql
# Todos os logs do serviço
{job="ARQUITETURA_FIBER_GORM"}

# Apenas erros
{job="ARQUITETURA_FIBER_GORM", level="error"}

# Logs de requisições HTTP
{job="ARQUITETURA_FIBER_GORM"} |= "Requisição HTTP"

# Logs com latência > 1s
{job="ARQUITETURA_FIBER_GORM"} | json | latency > 1s

# Erros nos últimos 5 minutos
{job="ARQUITETURA_FIBER_GORM", level=~"error|warn"} | json
```

### Dashboard Grafana

Para criar um dashboard, use as seguintes queries:

1. **Taxa de logs por nível**:
```logql
sum by (level) (rate({job="ARQUITETURA_FIBER_GORM"}[5m]))
```

2. **Contagem de erros**:
```logql
count_over_time({job="ARQUITETURA_FIBER_GORM", level="error"}[1h])
```

3. **Logs de requisições lentas**:
```logql
{job="ARQUITETURA_FIBER_GORM"} | json | latency > 500ms
```

### Desabilitar Loki

Para desabilitar o envio de logs ao Loki:

```bash
LOKI_ENABLED=false go run cmd/api/main.go
```

## 📈 Benefícios Demonstrados

- **Fiber**: Alta performance, sintaxe familiar (Express-like)
- **GORM**: ORM maduro, migrations automáticas, relacionamentos
- **Prometheus**: Métricas detalhadas para observabilidade
- **Loki/Grafana**: Agregação e visualização de logs centralizada
- **Arquitetura limpa**: Fácil manutenção e escalabilidade
- **Logs estruturados**: Facilita debugging e monitoramento
- **Swagger**: Documentação automática e interativa
- **Configuração flexível**: Variáveis de ambiente opcionais com defaults sensatos
