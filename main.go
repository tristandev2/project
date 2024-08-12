package main

import (
	"database/sql"
	"log"
	"net/http"
	"text/template"
	"time"

	"golang.org/x/crypto/bcrypt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/go-sql-driver/mysql"
)

var db *sql.DB
var templates = template.Must(template.ParseGlob("templates/*"))
var jwtKey = []byte("your_secret_key")

type Claims struct {
	Email string `json:"email"`
	jwt.RegisteredClaims
}

func ConexionDB() (*sql.DB, error) {
	// Datos de conexión a la base de datos
	driver := "mysql"
	usuario := "root"
	contrasenia := "" 
	nombre := "facturacion"

	conexion, err := sql.Open(driver, usuario+":"+contrasenia+"@tcp(127.0.0.1)/"+nombre)
	if err != nil {
		return nil, err // Devolver el error
	}

	// Verificar si la conexión es válida
	err = conexion.Ping()
	if err != nil {
		return nil, err 
	}

	return conexion, nil  
}

func main() {
	var err error
	db, err = ConexionDB()
	if err != nil {
		log.Fatal("error al conectar a la base de datos:", err)
	}
	defer db.Close()

	fs := http.FileServer(http.Dir("./img"))
	http.Handle("/img/", http.StripPrefix("/img/", fs))

	http.HandleFunc("/", Login) 
	http.HandleFunc("/registro", Registro)
	http.HandleFunc("/insertar", Insertar)
	http.HandleFunc("/loguearse", Loguearse)
	http.HandleFunc("/iniciop", iniciop) 

	log.Println("servidor activo")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Fatal("listenAndServe: ", err)
	}
}

func Login(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "Login.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Registro(w http.ResponseWriter, r *http.Request) {
	err := templates.ExecuteTemplate(w, "Registro.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func Insertar(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		name := r.FormValue("name")
		apellido := r.FormValue("apellido")
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Aquí haciendo hashing a la contraseña antes de guardarla
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			http.Error(w, "error al encriptar la contraseña", http.StatusInternalServerError)
			return
		}

		// Insertar los datos en la base de datos
		_, err = db.Exec("INSERT INTO usuarios (name, apellido, email, password) VALUES (?, ?, ?, ?)",
			name, apellido, email, hashedPassword)
		if err != nil {
			http.Error(w, "error al insertar el usuario en la base de datos", http.StatusInternalServerError)
			return
		}

		// Enviar respuesta al cliente
		err = templates.ExecuteTemplate(w, "Login.html", nil)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
		}
	} else {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	}
}

func Loguearse(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		// Obtener los datos del formulario
		email := r.FormValue("email")
		password := r.FormValue("password")

		// Consultar la base de datos para verificar la existencia del usuario
		var contrasenaAlmacenada string
		err := db.QueryRow("SELECT password FROM usuarios WHERE email = ?", email).Scan(&contrasenaAlmacenada)
		if err != nil {
			if err == sql.ErrNoRows {
				http.Error(w, "usuario no encontrado", http.StatusUnauthorized)
			} else {
				http.Error(w, "error al consultar la base de datos", http.StatusInternalServerError)
			}
			return
		}

		// Verificar la contraseña
		err = bcrypt.CompareHashAndPassword([]byte(contrasenaAlmacenada), []byte(password))
		if err != nil {
			http.Error(w, "contraseña incorrecta", http.StatusUnauthorized)
			return
		}

		// Crear el token JWT
		expirationTime := time.Now().Add(24 * time.Hour)
		claims := &Claims{
			Email: email,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(expirationTime),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
		tokenString, err := token.SignedString(jwtKey)
		if err != nil {
			http.Error(w, "error al generar el token", http.StatusInternalServerError)
			return
		}

		// Establecer la cookie con el token JWT
		http.SetCookie(w, &http.Cookie{
			Name:    "token",
			Value:   tokenString,
			Expires: expirationTime,
			Path:    "/",
		})

		log.Println("Token JWT establecido y redirigiendo a /iniciop")
		http.Redirect(w, r, "/iniciop", http.StatusSeeOther)
	} else {
		http.Error(w, "método no permitido", http.StatusMethodNotAllowed)
	}
}

func iniciop(w http.ResponseWriter, r *http.Request) {
	// Aquí se mostraría la página principal después del inicio de sesión
	log.Println("Acceso a la página de inicio")
	err := templates.ExecuteTemplate(w, "index.html", nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

