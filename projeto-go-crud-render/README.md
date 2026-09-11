# GoUser - CRUD de Usuarios

Projeto acadêmico em Go (Golang) com HTML/CSS e SQLite.

## Funcionalidades
- Inserir usuários
- Consultar usuários
- Alterar usuários
- Excluir usuários
- Validação de campos
- Bloqueio de e-mail duplicado

## Tecnologias
- Go (Golang)
- HTML
- CSS
- SQLite
- modernc.org/sqlite

## Executar localmente

No PowerShell, dentro da pasta do projeto:

```powershell
go version
go mod tidy
go run .
```

Depois abra:

```text
http://localhost:8080
```

O projeto usa a variável de ambiente `PORT`. Localmente, se ela não existir, a porta 8080 é usada automaticamente.

## Publicar no Render

O projeto já está preparado para o Render. O arquivo `render.yaml` contém a configuração do serviço.

### Opção manual
1. Envie os arquivos deste projeto para um repositório no GitHub.
2. Acesse o Render e escolha **New > Web Service**.
3. Conecte o repositório.
4. Selecione **Go** como linguagem.
5. Plano: **Free**.
6. Build Command:

```text
go build -tags netgo -ldflags "-s -w" -o app
```

7. Start Command:

```text
./app
```

8. Clique em **Create Web Service**.

O Render fornecerá um endereço público `onrender.com`.

### Atenção sobre o SQLite

O plano gratuito do Render possui sistema de arquivos temporário. Portanto, alterações feitas no arquivo `database.db` podem ser perdidas quando o serviço reiniciar, ficar ocioso ou receber um novo deploy. Para uma apresentação escolar isso pode ser suficiente, mas não é adequado para guardar dados importantes.

## Observação
A senha está armazenada de forma simples apenas para fins didáticos. Em um sistema real, ela deve ser armazenada usando hash.

O arquivo `executar.bat` continua disponível para facilitar a execução no Windows.
