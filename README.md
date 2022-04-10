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
