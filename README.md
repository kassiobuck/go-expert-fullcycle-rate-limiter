# go-expert-fullcycle-rate-limiter

### Objetivo:
Desenvolver um rate limiter em Go que possa ser configurado para limitar o número máximo de requisições por segundo com base em um endereço IP específico ou em um token de acesso.

### Requisitos:
- [Go SDK](https://golang.org/dl/): Linguagem de programação Go.
- [Docker](https://docs.docker.com/get-docker/): Plataforma de conteinerização.

### Executando o projeto:
- [Donwload do projeto](http://github.com/kassiobuck/go-expert-fullcycle-rate-limiter)
- Execute o comando `docker compose up -d`

### Configurando o projeto:
As váriaveis de ambiente estão disponiveis no arquivo .env
- É possivel alterar o numero de requisições por IP
- Tempo de bloqueio por IP
- Portas padrões de configuração
- JWT Secret key
- Tokens são configurados individualmente via requisição http (para testes).

### Exemplos:
Acessando api/server.http é possivel verificar as URL's disponiveis.
- É possivel gerar tokens (API_KEY) de acesso com limite maximo de requisições por segundo (max) e tempo de bloqueio (interval).
`http://localhost:8080/genToken?interval=10&max=2 `

- É possivel verificar se o token é válido e verificar sua configuração
`http://localhost:8080/decodeToken`

- Testar o rate limiter:
`GET http://localhost:8080`