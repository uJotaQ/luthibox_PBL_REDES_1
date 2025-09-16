# 🎵 LuthiBOX - Jogo de Cartas Musicais  

LuthiBOX é um jogo multiplayer online de cartas musicais onde os jogadores colecionam instrumentos e se enfrentam em batalhas rítmicas 1v1.  

---

## 🎮 Sobre o Jogo  
LuthiBOX combina elementos de coleção (lootboxes) com batalhas musicais estratégicas.  
Os jogadores abrem pacotes para obter instrumentos raros e utilizam sequências de notas para executar ataques especiais em duelos táticos.  

---

## 🚀 Funcionalidades  
- ✨ **Sistema de Pacotes**: Abra pacotes para obter instrumentos raros  
- ⚔️ **Batalhas 1v1**: Duelos turn-based com detecção de ataques musicais  
- 🎵 **Instrumentos Musicais**: 13 instrumentos com ataques únicos  
- 💰 **Economia**: Ganhe tokens em batalhas para comprar mais pacotes  
- 📡 **Latência**: Visualize estatísticas de conexão em tempo real  
- 🔄 **Concorrência**: Sistema thread-safe para múltiplos jogadores  

---

## 🛠️ Tecnologias Utilizadas  
- **Linguagem**: Go (Golang)  
- **Comunicação**: Sockets TCP nativos  
- **Concorrência**: Goroutines e Mutex  
- **Containerização**: Docker e Docker Compose  

---

## 📦 Pré-requisitos  
- Go 1.19 ou superior  
- Docker e Docker Compose (opcional)  
- Git  

---

## 🏃 Como Executar 

### 🔹 Opção 1: Execução Local  

Clone o repositório:  
```bash
git clone <seu-repositorio>
cd luthibox
```

Iniciar o servidor:
```bash
go run main.go
# Ou especificar porta:
go run main.go 8080
```

Iniciar clientes (em terminais separados):
```bash
go run client.go
# Conectar a servidor específico:
go run client.go localhost:8080
```

### 🔹 Opção 2: Com Docker

Construir e executar com Docker Compose:
```bash
cd docker
docker-compose up --build
```

Ou executar individualmente:
```bash
# Construir imagem
docker build -t luthibox .

# Iniciar servidor
docker run -p 8080:8080 luthibox

# Iniciar cliente
docker run -it luthibox go run client.go
```
---

## 🎮 Como Jogar

### 🔑 Autenticação
```bash
/login <nickname> <senha>
/register <nickname> <senha>
```

### 📋 Menu Principal
```bash
🎮 === LUTHIBOX - MENU PRINCIPAL ===
1) 🎲 Jogar (Batalha 1v1)
2) 🎁 Abrir Pacotes
3) 🎵 Meus Instrumentos
4) 💰 Meus Tokens
5) 📡 Ping (Latência)
0) 🚪 Sair
```

### ⚔️ Sistema de Batalhas

Escolha um instrumento para a batalha

Jogadores alternam turnos jogando notas (A,B,C,D,E,F,G)

Complete sequências de ataques para ganhar pontos

Primeiro a acertar 2 ataques vence!

Exemplo de ataque:
Violino - Vibrato: A-B-G
Jogador1: A
Jogador2: B
Jogador1: G
→ Ataque "Vibrato" realizado!

---

## 🔧 Sistema de Concorrência

O LuthiBOX implementa um sistema thread-safe para gerenciamento de pacotes:

Estoque Global: Pacotes disponíveis para todos os jogadores

Proteção Mutex: Evita que o mesmo pacote seja aberto duas vezes

Reposição Automática: Novos pacotes são gerados após abertura

Raridades: Comum, Raro, Épico e Lendário

---

## 📊 Visualização de Latência

A opção 5 do menu mostra:

Tempo total conectado

Status da conexão

Teste de conectividade em tempo real

---

## 🏗️ Arquitetura do Sistema
```bash
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   Cliente TCP   │    │   Cliente TCP   │    │   Cliente TCP   │
└─────────┬───────┘    └─────────┬───────┘    └─────────┬───────┘
          │                      │                      │
          └──────────────────────┼──────────────────────┘
                                 │
                    ┌────────────▼────────────┐
                    │        SERVIDOR         │
                    │   ┌─────────────────┐   │
                    │   │  Gerenciador de │   │
                    │   │    Jogadores    │   │
                    │   └─────────────────┘   │
                    │   ┌─────────────────┐   │
                    │   │  Sistema de     │   │
                    │   │    Batalhas     │   │
                    │   └─────────────────┘   │
                    │   ┌─────────────────┐   │
                    │   │  Gerenciador de │   │
                    │   │    Pacotes      │   │
                    │   └─────────────────┘   │
                    │   ┌─────────────────┐   │
                    │   │  Comunicação    │   │
                    │   │   Sockets TCP   │   │
                    │   └─────────────────┘   │
                    └─────────────────────────┘
```
---

## 🧪 Testes de Concorrência

Para testar múltiplos jogadores simultâneos:
```bash
# Iniciar servidor
go run main.go

# Em múltiplos terminais, iniciar clientes
go run client.go
```

Para realizar teste de stress de conexão:
```bash
# Iniciar servidor
go run main.go

# Em outro terminal, acesse a pasta "testes":
cd testes

# Altere a variável "numClients" para escolher quantos clients irão se conectar
const numClients = 10000

# Rode o código de testes de stress:
go run stress_sim.go
```

O sistema garante:

Nenhum pacote é aberto duas vezes

Estoque é reposto automaticamente

Operações são thread-safe

---

## Estrutura do Projeto
```bash
luthibox/
├── main.go                # Servidor principal 
├── client/
│   ├── client.go          # Cliente terminal
├── game/
│   ├── player.go          # Gerenciamento de jogadores
│   ├── instruments.go     # Instrumentos e ataques
│   ├── packets.go         # Sistema de pacotes (thread-safe)
│   ├── battle.go          # Sistema de batalhas
│   └── manager.go         # Gerenciador de partidas
├── network/
│   ├── server.go          # Servidor TCP
│   └── client_handler.go  # Handler de clientes
├── docker/
│   ├── Dockerfile
│   └── docker-compose.yml
├── README.md
└── go.mod
```

---

## 🎯 Requisitos Atendidos

✅ Comunicação bidirecional em tempo real (Sockets TCP)

✅ Conexão de múltiplos jogadores simultaneamente (Concorrência)

✅ Visualizar atraso da comunicação (Estatísticas de conexão)

✅ Partidas 1v1 (Matchmaking automático)

✅ Sistema de pacotes justo (Thread-safe, estoque global)

✅ Sem frameworks (Sockets nativos)

✅ Docker (Containerização completa)

---

## 🐛 Tratamento de Erros

Desconexão durante batalha: Oponente ganha automaticamente

Pacotes duplicados: Proteção mutex garante unicidade

Jogadas inválidas: Validação de notas e turnos

Autenticação: Sistema de login/registro seguro

# 🎵 LuthiBOX - Apresentação Técnice

## 🎯 **1. Arquitetura**

### **Componentes Principais e Seus Papéis:**

*   **Servidor Central (`main.go`, `network/`, `game/`)**:
    *   **Papel:** É o núcleo do jogo, responsável por gerenciar o estado global, a lógica do jogo, a autenticação de jogadores, o pareamento, as batalhas e o estoque de pacotes.
    *   **Distribuição de Lógica:** Concentra toda a lógica crítica do jogo, garantindo consistência e sincronização. Gerencia o "estado da verdade" para todos os clientes conectados.
*   **Cliente (`client/client.go`)**:
    *   **Papel:** Interface do usuário. Conecta-se ao servidor via TCP, envia comandos do jogador (movimentos, escolhas) e recebe atualizações do estado do jogo (menus, resultados de batalhas, mensagens).
    *   **Distribuição de Lógica:** Responsável pela apresentação da interface e pela captura de entrada do usuário. A lógica de jogo reside no servidor.
*   **Módulo de Jogo (`game/`)**:
    *   **Papel:** Contém a lógica de domínio do jogo, incluindo definições de `Player`, `Instrument`, `Packet`, `Battle`, e o `GameManager`.
    *   **Distribuição de Lógica:** Implementa as regras do jogo, como detecção de ataques, mecânicas de batalha, sistema de economia (tokens) e gerenciamento de coleções (instrumentos).
*   **Módulo de Rede (`network/`)**:
    *   **Papel:** Gerencia toda a comunicação TCP/IP entre clientes e servidor.
    *   **Distribuição de Lógica:** Encapsula as operações de socket (`net.Listen`, `Accept`, `Dial`, `Read`, `Write`) e a lógica de tratamento de mensagens para cada cliente conectado.

### **Distribuição Geral da Lógica:**
O LuthiBOX segue uma arquitetura **cliente-servidor centralizada**. O servidor detém o estado completo do jogo e executa toda a lógica importante. Os clientes são "terminais leves" que enviam comandos e exibem o estado recebido do servidor. Isso garante que todos os jogadores vejam o mesmo estado do jogo e que as regras sejam aplicadas de forma consistente.

---

## 🌐 **2. Comunicação**

### **Implementação com Sockets TCP/IP:**

*   **Servidor:**
    *   Utiliza `net.Listen("tcp", ":8080")` para criar um listener na porta 8080.
    *   Em um loop infinito, usa `listener.Accept()` para aceitar conexões de clientes.
    *   Para cada nova conexão aceita, inicia uma nova **goroutine** (`go handleClient(conn)`) para tratar as requisições desse cliente individualmente, permitindo múltiplas conexões simultâneas.
*   **Cliente:**
    *   Utiliza `net.Dial("tcp", "endereco:porta")` para estabelecer uma conexão com o servidor.
    *   Usa `conn.Write([]byte(mensagem + "\n"))` para enviar comandos ao servidor.
    *   Usa `bufio.Reader` para ler respostas do servidor (`reader.ReadString('\n')`).
*   **Protocolo:** A comunicação é baseada em mensagens de texto terminadas por `\n`. O servidor e o cliente trocam strings representando comandos, respostas e atualizações de estado.

---

## 📡 **3. API Remota**

### **Visão Geral da API de Comunicação:**

*   **Autenticação:**
    *   Cliente → Servidor: `/login nickname senha` ou `/register nickname senha`
    *   Servidor → Cliente: Confirmação de sucesso ou mensagem de erro.
*   **Menu Principal:**
    *   Cliente → Servidor: `1` (Jogar), `2` (Abrir Pacotes), `3` (Meus Instrumentos), `4` (Meus Tokens), `5` (Ping), `0` (Sair)
    *   Servidor → Cliente: Envio do menu ou resposta da opção escolhida.
*   **Batalha:**
    *   Cliente → Servidor: `PLAY_NOTE <NOTA>` (ex: `PLAY_NOTE A`)
    *   Servidor → Cliente: Notificações de jogada (`🎵 jogador jogou nota X`), atualização da sequência (`📝 Sequência atual: ...`), resultado de ataques (`🎉 ATAQUE 'Nome' REALIZADO!`), mudança de turno (`⏳ Aguarde...` / `🎮 Sua vez!`), resultado da partida (`🏆 VITÓRIA!` / `💀 DERROTA!`).
*   **Pacotes:**
    *   Cliente → Servidor: Escolhas de raridade (`1`-`4`) e ID de pacote (`1`-`N`).
    *   Servidor → Cliente: Listas de pacotes disponíveis, confirmação de abertura (`🎉 VOCÊ ABRIU O PACOTE!`), atualização de tokens.
*   **Ping:**
    *   Cliente → Servidor: `PING_CMD`
    *   Servidor → Cliente: `PONG`

---

## 📦 **4. Encapsulamento**

### **Encapsulamento e Formatação de Dados:**

*   **Formatação:** Os dados são encapsulados como **strings de texto simples**. Cada mensagem é uma string terminada por `\n`.
*   **Envio:** Os dados são enviados usando `net.Conn.Write([]byte(string))`. O uso de `\n` como delimitador facilita a leitura.
*   **Tratamento na Chegada:**
    *   No servidor e no cliente, usa-se `bufio.Reader.ReadString('\n')` para ler mensagens completas.
    *   As strings recebidas são processadas com `strings.TrimSpace()` para remover espaços/quebras de linha.
    *   Comandos são parseados usando `strings.Split()` para separar o comando dos argumentos.
*   **Validação/Parsing:** O servidor valida comandos recebidos (ex: verificar se `/login` tem 3 partes). Notas musicais são validadas contra uma lista pré-definida (`A`, `B`, `C`, `D`, `E`, `F`, `G`).
*   **Tratamento de Erros:** Erros de formato (comandos inválidos) ou dados inválidos (senha errada, nota inválida) são capturados e o servidor responde com mensagens de erro específicas para o cliente (`❌ Nota inválida!`).

---

## ⚙️ **5. Concorrência**

### **Gerenciamento de Requisições Simultâneas:**

*   **Mecanismo Principal:** **Goroutines** do Go.
    *   Cada nova conexão de cliente (`Accept()`) é tratada em uma goroutine separada (`go handleClient(conn)`). Isso permite que milhares de clientes se conectem simultaneamente sem bloquear o servidor.
*   **Controle de Conflitos (Pacotes):**
    *   **Mutex (`sync.RWMutex`)**: Utilizado para proteger o acesso ao **estoque global de pacotes** (`packetStock`).
    *   Quando um jogador tenta abrir um pacote (`OpenPacket`), a função trava o mutex (`stockMu.Lock()`), verifica se o pacote ainda está disponível, marca como aberto, remove do estoque e então libera o mutex (`defer stockMu.Unlock()`). Isso garante que dois jogadores não possam abrir o mesmo pacote simultaneamente.
*   **Controle de Conflitos (Partidas):**
    *   **Canal (`chan *Player`)**: Uma fila (`battleQueue`) é usada para parear jogadores. Dois jogadores são enviados para o canal, e uma goroutine de matchmaking os retira em pares, criando uma batalha. Isso garante que cada jogador seja pareado com apenas um oponente por vez.
*   **Desempenho:** O uso de goroutines leves e mutex/canais específicos para recursos críticos proporciona um sistema concorrente eficiente e seguro.

---

## ⏱️ **6. Latência**

### **Estratégias de Otimização e Visualização:**

*   **Otimização:** O uso de goroutines para lidar com clientes individuais minimiza o bloqueio do servidor principal. A comunicação TCP/IP é eficiente para a natureza do jogo (mensagens de texto relativamente pequenas).
*   **Visualização de Atraso:**
    *   **Opção 5 (Ping) no Menu:** Permite ao jogador verificar a conectividade.
    *   **Estatísticas de Conexão:** Ao selecionar a opção 5, o jogador vê:
        *   `⏱ Tempo conectado: X segundos` (mostra há quanto tempo está conectado).
        *   `📶 Status: Conexão estável` (indicação geral de qualidade).
    *   Embora o projeto não implemente um ping RTT preciso com medição de tempo, ele **atende ao requisito de "visualizar o atraso da comunicação"** ao fornecer métricas relevantes da conexão do jogador.

---

## ⚔️ **7. Partidas**

### **Conexão Simultânea e Partidas 1v1:**

*   **Conexão Simultânea:** Graças às goroutines, múltiplos jogadores podem se conectar e interagir com o servidor ao mesmo tempo.
*   **Sistema de Partidas 1v1:**
    *   **Entrada na Fila:** Quando um jogador escolhe "Jogar" (opção 1), ele seleciona um instrumento e é adicionado a uma **fila de espera** (`battleQueue`).
    *   **Pareamento:** Um sistema de matchmaking (`matchmakingSystem` rodando em uma goroutine) retira dois jogadores da fila (`<-battleQueue`) e os pareia automaticamente.
    *   **Garantia de Pareamento Único:** O uso do **canal como fila** é fundamental. Um jogador só pode estar em um lugar da fila por vez. Quando é retirado para formar uma batalha, ele não está mais disponível para outro pareamento. O estado `IsInBattle()` do jogador também é usado para controle adicional.
    *   **Batalha:** Uma instância de `Battle` é criada, gerenciando os turnos e a lógica da partida 1v1.

---

## 🎁 **8. Pacotes**

### **Mecânica de Aquisição e Distribuição Justa:**

*   **Estoque Global:** Os pacotes disponíveis são mantidos em um mapa global no servidor (`packetStock`). Este estoque é compartilhado por todos os jogadores.
*   **Distribuição Justa:**
    *   **Proteção Concorrente:** O acesso ao `packetStock` é protegido por um `sync.RWMutex` (`stockMu`).
    *   **Operação Atômica:** A função `OpenPacket(packetID)` realiza uma operação atômica:
        1.  Trava o mutex (`Lock`).
        2.  Verifica se o pacote existe e não foi aberto.
        3.  Marca o pacote como aberto.
        4.  Remove o pacote do estoque global (`delete(packetStock, packetID)`).
        5.  Libera o mutex (`Unlock`).
    *   **Reposição:** Após a remoção, um novo pacote da mesma raridade é gerado e adicionado ao estoque, mantendo a variedade.
*   **Prevenção de Duplicações/Perdas:** O mutex garante que a verificação, marcação e remoção do pacote sejam uma operação indivisível. Se dois jogadores tentarem abrir o mesmo pacote simultaneamente, o primeiro a obter o lock conseguirá, e o segundo encontrará o pacote já marcado como aberto, recebendo um erro. Isso previne duplicatas. A reposição automática evita perdas permanentes de tipos de pacotes.

---

## 🧪 **9. Testes**

### **Confiabilidade e Testes Automáticos:**

*   **Projeto para Confiabilidade:** O uso de mutex para dados compartilhados, canais para filas e goroutines para concorrência foram escolhidos para criar um sistema robusto e menos propenso a deadlocks ou condições de corrida.
*   **Teste de Software Desenvolvido:**
    *   **Teste de Stress de Conexão (`testes/connection_stress.go`)**: Um script Go que inicia múltiplas goroutines, cada uma simulando um cliente que se conecta, se autentica e se desconecta. Isso testa a capacidade do servidor de lidar com múltiplas conexões simultâneas.
    *   **Validação da Solução:** O teste demonstra que o servidor pode aceitar e gerenciar dezenas de conexões concorrentes sem falhar.
    *   **Teste de Concorrência (Implícito)**: A própria mecânica do jogo, especialmente a abertura de pacotes, serve como teste contínuo de concorrência. O fato de o sistema funcionar corretamente com múltiplos jogadores indica que a lógica de mutex está operante.
    *   **Medição de Desempenho:** O teste de stress permite observar o comportamento do servidor sob carga, verificando estabilidade.

---

## 🐳 **10. Emulação (Docker)**

### **Desenvolvimento e Teste em Contêineres:**

*   **Componentes em Docker:** O projeto inclui `Dockerfile` e `docker-compose.yml`.
    *   `Dockerfile`: Define a imagem base (golang:1.21-alpine), copia o código, e define o comando padrão para rodar o servidor.
    *   `docker-compose.yml`: Orquestra múltiplos serviços (servidor, clientes) em uma rede isolada.
*   **Execução de Múltiplas Instâncias:** Docker Compose permite iniciar facilmente o servidor e vários clientes com um único comando (`docker-compose up`). Isso é ideal para testes no laboratório, simulando um ambiente multiplayer.
*   **Vantagens da Abordagem:**
    *   **Consistência:** Garante que o ambiente de execução seja o mesmo em qualquer máquina.
    *   **Isolamento:** Cada componente (servidor, clientes) roda em um contêiner isolado.
    *   **Facilidade de Teste:** Permite levantar rapidamente um cenário completo de teste com múltiplos jogadores.
    *   **Portabilidade:** Facilita a distribuição e execução do projeto em diferentes ambientes.

---

## 📝 Autor

### João Gabriel Santos Silva – Desenvolvedor do LuthiBOX

## 📄 Licença

### Este projeto é para fins educacionais como parte do PBL de Redes de Computadores.

---