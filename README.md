# 📌 **README.md**

# 🚀 API Backend - Módulos Independentes

## 📖 Sobre o Projeto
Este projeto implementa uma API utilizando **Golang** e **Gin Framework** para fornecer funcionalidades financeiras, incluindo extração de dados e cálculo de **PPR** (Participação nos Lucros e Resultados).  

A API também possui autenticação via **JWT/Bearer Token** e segue um design modular, facilitando a escalabilidade.

---

## 🛠️ **Tecnologias Utilizadas**
- **Go** (Golang)
- **Gin Framework** (Para roteamento e middleware)
- **MariaDB** (Banco de dados relacional)
- **JWT** (Autenticação segura)
- **dotenv** (Gerenciamento de variáveis de ambiente)
- **Logging** (Log estruturado para depuração)

---

## 📂 **Estrutura do Projeto**
```
/app
├── db
│   └── db.go
├── go.mod
├── go.sum
├── LICENSE
├── logger
│   └── logger.go
├── logs
│   └── app.log
├── main.go
├── README.md
├── server
│   ├── healthCheck.go
│   ├── import.go
│   ├── modules
│   │   ├── finance
│   │   │   ├── files
│   │   │   ├── routes.go
│   │   │   └── services.go
│   │   ├── login
│   │   │   ├── auth_utils
│   │   │   │   └── jwt_utils.go
│   │   │   ├── controllers
│   │   │   │   ├── auth_controller.go
│   │   │   │   └── users_controller.go    # BlankFile
│   │   │   ├── middleware
│   │   │   │   ├── auth_middleware.go
│   │   │   │   └── traefik_auth.go
│   │   │   ├── models
│   │   │   │   └── user_model.go
│   │   │   ├── routes.go
│   │   │   └── services
│   │   │       └── auth_service.go
│   │   └── ppr
│   │       ├── routes.go
│   │       └── service.go
│   ├── server.go
│   └── static
│       └── favicon.ico
├── tmp
│   ├── build-errors.log
│   └── main
└── utils
    └── utils.go
```

---

## 🔧 **Configuração do Ambiente**

1. **Clone o repositório**
   ```sh
   git clone git@github.com:vandinho10/api_backend.git
   cd api_backend
   ```

2. **Configure as variáveis de ambiente**  
   Renomeie o arquivo `.env.example` para `.env` e preencha de acordo com seu ambiente:
   ```ini
   # Database Config
   DB_USER=           # Usuário do banco de dados
   DB_PASS=           # Senha do banco de dados
   DB_HOST=           # Endereço do servidor do banco de dados (ex: localhost, IP)
   DB_PORT=           # Porta do banco de dados (ex: 3306 para MySQL/MariaDB)
   DB_NAME=           # Nome do banco de dados

   # JWT Config
   JWT_SECRET=        # Chave secreta para assinatura do JWT
   JWT_EXPIRE=        # Tempo de expiração do JWT (ex: 3600s para 1 hora)

   # GIN MODES: release, debug, test
   GIN_MODE=          # Modo de execução do Gin (release, debug ou test)

   # LOG LEVELS: DEBUG, INFO, WARN, ERROR, FATAL, PANIC
   LOG_LEVEL=         # Nível de log para o sistema (ex: DEBUG)

   # Server Config - Default Port 80 (executando em docker, com gerenciamento de Certificados e Redirecionamento para HTTPS pelo Traefik)
   PORT_HTTP=         # Porta HTTP do servidor (ex: 8080)

   # Domain Config - Para Configuração do CORS
   DOMAIN_NAME=       # Nome do domínio para CORS (ex: http://localhost ou http://meudominio.com)

   # Bearer Protected Paths
   BEARER_PROTECTED_PATHS=  # Caminhos da API que requerem autenticação com Bearer Token (ex: /finance, /ppr)

   # Finances API
   FINANCE_PATH=      # Caminho base para o módulo financeiro (ex: /finance)
   FINANCE_CSV=       # Rota para extração de extrato financeiro (ex: /extract)
   FINANCE_CSV_DB=    # Rota para extração de extrato financeiro do banco de dados (ex: /extract_db)
   ```

3. **Instale as dependências**
   ```sh
   go mod tidy
   ```

4. **Inicie o servidor**
   ```sh
   go run main.go
   ```

---

## 🔑 **Autenticação**
As rotas protegidas utilizam **JWT/Bearer Token**.  
Para acessar endpoints protegidos, envie o token no cabeçalho `Authorization`:
```
Authorization: Bearer SEU_TOKEN_AQUI
```

---

## 🚀 **Endpoints da API**

### 🏦 **Módulo Financeiro**
> Rota de acordo com o informado no .env
- `POST /finance/extract` → Retorna um arquivo CSV de extrato financeiro.
- `POST /finance/extract_db` → Retorna um arquivo CSV de extrato financeiro do banco de dados.

### 📊 **Módulo PPR**
- `GET /ppr/calculate?salary={valor}&ppr_value={valor}&months_worked={valor}`  
  - Calcula a participação nos lucros com base no salário e no tempo de trabalho.

### 🔒 **Módulo de Login**
- **(Em desenvolvimento)**

  - **POST /login**
    - **Descrição**: Realiza o login do usuário e retorna um token JWT.
    - **Corpo da Requisição**: `{ "email": "user@example.com", "password": "password123" }`
    - **Resposta**:
      - 200 OK: `{ "token": "jwt_token", "user": { ...user_data... }, "time_remaining": "3600s" }`
      - 400 Bad Request: Se as credenciais forem inválidas.

  - **POST /register**
    - **Descrição**: Cria um novo usuário no sistema.
    - **Corpo da Requisição**: `{ "name": "John Doe", "email": "user@example.com", "password": "password123" }`
    - **Resposta**:
      - 201 Created: Confirmação de sucesso no registro do usuário.
      - 400 Bad Request: Se o e-mail ou nome de usuário já estiverem em uso.

  - **POST /logout**
    - **Descrição**: Realiza o logout do usuário e revoga o token.
    - **Corpo da Requisição**: `{ "token": "jwt_token" }`
    - **Resposta**:
      - 200 OK: Confirmação de que o token foi revogado.
      - 400 Bad Request: Se o token não for válido ou não estiver na blacklist.

---

## 📜 **Licença**

Este projeto está licenciado sob a **MIT License**. Consulte o arquivo [LICENSE](./LICENSE) para mais detalhes.

A licença MIT é uma das licenças mais permissivas e amplamente utilizadas, permitindo que os desenvolvedores usem, modifiquem e distribuam o software de maneira livre, com a única exigência de incluir o aviso de copyright e a licença em todas as cópias do software.

---

## 📞 **Contato**
Dúvidas ou sugestões? Entre em contato!  
📧 **Email:** contato@vandinho.com.br
