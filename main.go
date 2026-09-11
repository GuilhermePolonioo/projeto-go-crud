package main

import (
	"database/sql"
	"html/template"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	_ "modernc.org/sqlite"
)

type User struct {
	ID    int
	Nome  string
	Email string
}

var db *sql.DB

func main() {
	var err error
	db, err = sql.Open("sqlite", "database.db")
	if err != nil {
		log.Fatal("erro ao abrir banco: ", err)
	}
	defer db.Close()

	if err = criarTabela(); err != nil {
		log.Fatal("erro ao criar tabela: ", err)
	}

	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.Dir("static"))))
	http.HandleFunc("/", home)
	http.HandleFunc("/usuarios", listarUsuarios)
	http.HandleFunc("/cadastro", cadastro)
	http.HandleFunc("/editar", editar)
	http.HandleFunc("/excluir", excluir)

	log.Println("======================================")
	log.Println(" GoUser - CRUD de Usuarios")
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Println(" Servidor iniciado na porta " + port)
	log.Println(" Acesso local: http://localhost:" + port)
	log.Println("======================================")
	log.Fatal(http.ListenAndServe("0.0.0.0:"+port, nil))
}

func criarTabela() error {
	_, err := db.Exec(`CREATE TABLE IF NOT EXISTS usuarios (
		id INTEGER PRIMARY KEY AUTOINCREMENT,
		nome TEXT NOT NULL,
		email TEXT NOT NULL UNIQUE,
		senha TEXT NOT NULL
	)`)
	return err
}

func render(w http.ResponseWriter, arquivo string, dados any) {
	tmpl, err := template.ParseFiles("templates/" + arquivo)
	if err != nil {
		http.Error(w, "Erro ao carregar pagina.", 500)
		return
	}
	if err := tmpl.Execute(w, dados); err != nil {
		http.Error(w, "Erro ao renderizar pagina.", 500)
	}
}

func home(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/" {
		http.NotFound(w, r)
		return
	}
	render(w, "index.html", nil)
}

func listarUsuarios(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Metodo nao permitido.", 405)
		return
	}
	rows, err := db.Query("SELECT id, nome, email FROM usuarios ORDER BY id DESC")
	if err != nil {
		http.Error(w, "Erro ao consultar usuarios.", 500)
		return
	}
	defer rows.Close()

	usuarios := []User{}
	for rows.Next() {
		var u User
		if err := rows.Scan(&u.ID, &u.Nome, &u.Email); err != nil {
			http.Error(w, "Erro ao ler usuarios.", 500)
			return
		}
		usuarios = append(usuarios, u)
	}
	render(w, "usuarios.html", usuarios)
}

func cadastro(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		render(w, "cadastro.html", nil)
		return
	}
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido.", 405)
		return
	}

	nome := strings.TrimSpace(r.FormValue("nome"))
	email := strings.TrimSpace(r.FormValue("email"))
	senha := r.FormValue("senha")

	if nome == "" || email == "" || senha == "" {
		render(w, "cadastro.html", map[string]string{"Erro": "Preencha todos os campos."})
		return
	}

	_, err := db.Exec("INSERT INTO usuarios (nome, email, senha) VALUES (?, ?, ?)", nome, email, senha)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			render(w, "cadastro.html", map[string]string{"Erro": "Este e-mail ja esta cadastrado."})
			return
		}
		http.Error(w, "Erro ao cadastrar usuario.", 500)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func editar(w http.ResponseWriter, r *http.Request) {
	id, err := strconv.Atoi(r.URL.Query().Get("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID invalido.", 400)
		return
	}

	if r.Method == http.MethodGet {
		var u User
		err := db.QueryRow("SELECT id, nome, email FROM usuarios WHERE id = ?", id).Scan(&u.ID, &u.Nome, &u.Email)
		if err == sql.ErrNoRows {
			http.NotFound(w, r)
			return
		}
		if err != nil {
			http.Error(w, "Erro ao consultar usuario.", 500)
			return
		}
		render(w, "editar.html", u)
		return
	}

	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido.", 405)
		return
	}

	nome := strings.TrimSpace(r.FormValue("nome"))
	email := strings.TrimSpace(r.FormValue("email"))
	senha := r.FormValue("senha")
	if nome == "" || email == "" {
		http.Error(w, "Nome e e-mail sao obrigatorios.", 400)
		return
	}

	if senha != "" {
		_, err = db.Exec("UPDATE usuarios SET nome=?, email=?, senha=? WHERE id=?", nome, email, senha, id)
	} else {
		_, err = db.Exec("UPDATE usuarios SET nome=?, email=? WHERE id=?", nome, email, id)
	}
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "unique") {
			http.Error(w, "Este e-mail ja esta cadastrado.", 409)
			return
		}
		http.Error(w, "Erro ao atualizar usuario.", 500)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}

func excluir(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Metodo nao permitido.", 405)
		return
	}
	id, err := strconv.Atoi(r.FormValue("id"))
	if err != nil || id <= 0 {
		http.Error(w, "ID invalido.", 400)
		return
	}
	if _, err := db.Exec("DELETE FROM usuarios WHERE id=?", id); err != nil {
		http.Error(w, "Erro ao excluir usuario.", 500)
		return
	}
	http.Redirect(w, r, "/usuarios", http.StatusSeeOther)
}
