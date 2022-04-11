# ExchengeChalenger
Desenvolvimento em GO de APIs e front-end


Para inicializar o desenvolvimento das API's instalar a biblioteca:

```console
~$ go get github.com/gin-gonic/gin
```

Database utilizado foi o Postgres e para criar o server foi utilizado:

```console
  $ sudo -u postgres psql
  postgres$: CREATE USER userwallet SUPERUSER;
  postgres$: \password userwallet
  postgres$: Enter new password for user "userwallet": postgres
```
Caso queria mudar alguma configuração do DATABASE é só mudar o arquivo dbconfig/dbconfig.go. e apartir dai ele criará automaticamente as tabelas e adicionará alguns dados básicos.

Pós isso, seguimos para a execução:

```console
  $ go run main.go
```

Assim, ele abrirá a porta do localhost:8080 para o desenvolvimento das API's.

## Deposit
A primeira API temos a de deposit, nela teremos 3 querys:
#### User
Aqui teremos o nome da pessoa que iremos depositar o token na carteira do usuário ('ben', 'eric', 'evandro': nomes básicos)

#### Currency
Neste a sigla do token que irá depositar o token na carteira do usuário(btc, eth, ada, xrp ou doge: tokens já pré selecionandos)

#### Amount
E por fim, a quantidade do token que irá depositar o token na carteira do usuário(qualquer número)

Por exemplo:
```browser
  http://localhost:8080/deposit?user=ben&currency=btc&amount=0.5
```

## Withdraw
A segunda API temos a de withdraw, nela teremos 3 querys:
#### User
Aqui teremos o nome da pessoa que iremos retirar o token na carteira do usuário ('ben', 'eric', 'evandro': nomes básicos)

#### Currency
Neste a sigla do token que irá retirar o token na carteira do usuário(btc, eth, ada, xrp ou doge: tokens já pré selecionandos)

#### Amount
E por fim, a quantidade do token que irá retirar o token na carteira do usuário(qualquer número)

Por exemplo:
```browser
  http://localhost:8080/withdraw?user=ben&currency=btc&amount=0.3
```

## Balance
A terceira API iremos ter um GET da carteira comppleta do usuário. E, por isso, irá ter somente uma query:

#### User
Adicionando a query com o user é onde pegaremos o nome do usuário para mostrar sua carteira completa ('ben', 'eric', 'evandro': nomes básicos)

Por exemplo:
```browser
  http://localhost:8080/balance?user=ben
```

## Balance
E, por fim, verificamos o histórico de transação em que o usuário realizaou. E nela terá 1 query:

#### User
Adicionando a query com o user é onde pegaremos o histórico de transaações do mesmo. ('ben', 'eric', 'evandro': nomes básicos)

Por exemplo:
```browser
  http://localhost:8080/history?user=ben
```

**Nesta API não consegui pegar o histórico pelo periodo desejado.