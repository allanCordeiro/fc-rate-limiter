# FC Rate Limiter

## Objetivo
Desafio Pos Tech Go Expert para desenvolver um rate limiter.

## Escopo

O rate limiter trabalha limitando por IP, ou por tokens. No caso do IP, o limite de tentativas e tempo de expiração são definidos em variáveis de ambiente:

`RATE_LIMITER_IP_BLOCK_TIME` - Tempo de expiração da requisição em *segundos*

A variável de ambiente `RATE_LIMITER_LIMIT` - define o throughput.

Caso a varíavel de ambiente não esteja definido, as mesmas serão injetadas do arquivo .env na raiz do projeto.

Para trabalhar com Tokens, é necessário inclui-lo na estrutura de arquivo `tokens.json` na raiz do projeto, estrutura similar a abaixo:

```
    {
        "token": "your-beloved-token",
        "expiresIn": 10 
    },
```

sendo:
 - Token: a chave, do tipo string
 - expiresIn: tempo de expiração da requisição em segundos.

O troughput usa a variável de ambiente `RATE_LIMITER_LIMIT`.

## Uso 
O rate limiter é um *Package*, que contém o middleware. Seu exemplo injetando-o em um server HTTP consta no main desse projeto. Para testes, após subir o serviço, executar uma chamada ao `http://localhost:8080`. Este acesso fara o controle do throughput por IP. Para realizar via token, é necessário injetar o paramêtro `API_KEY` no header, inserindo um token válido (presente no arquivo tokens.json). Caso o token inserindo não esteja validado, o que passa a limitar o throughput é o endereço IP da requisição.

## Iniciando os serviços
Executar o seguinte comando:

`make up`

O servidor HTTP de exemplo, assim como Redis, irão subir. O http deverá estar disponível na porta 8080.

## Testes
Os testes automatizados podem ser executados com o seguinte comando:

`make test`

## Redis

O uso do redis foi parte do requisito do projeto. A configuração do mesmo se encontra nas variáveis de ambiente:

`REDIS_HOST`

`REDIS_PORT`

Para implementação outro serviço de cache, basta criar uma implementação dentro do pacote ratelimiter, em cache, que atenda sua interface. Após gerar a implementação, realizar os ajustes no código cliente.S