package main  
  
import (  
	"encoding/json"  
	"fmt"  
	"html/template"  
	"log"  
	"net/http"  
	"os"  
)  
  
type User struct {  
	Username string `json:"username"`  
	Password string `json:"password"`  
}  
  
var users []User  
  
func main() {  
	// Load existing users from users.json  
	loadUsers()  
  
	http.HandleFunc("/login", loginHandler)  
	http.HandleFunc("/register", registerHandler)  
	http.HandleFunc("/", homeHandler)  
  
	fmt.Println("Server is running on http://localhost:8080")  
	log.Fatal(http.ListenAndServe(":8080", nil))  
}  
  
func loadUsers() {  
	file, err := os.Open("users.json")  
	if err != nil {  
		if os.IsNotExist(err) {  
			// If the file does not exist, create an empty slice of users  
			users = []User{}  
			return  
		}  
		log.Fatalf("Failed to open users.json: %v", err)  
	}  
	defer file.Close()  
  
	decoder := json.NewDecoder(file)  
	if err := decoder.Decode(&users); err != nil {  
		log.Fatalf("Failed to decode users.json: %v", err)  
	}  
}  
  
func saveUsers() {  
	file, err := os.Create("users.json")  
	if err != nil {  
		log.Fatalf("Failed to create users.json: %v", err)  
	}  
	defer file.Close()  
  
	encoder := json.NewEncoder(file)  
	encoder.SetIndent("", "  ")  
	if err := encoder.Encode(users); err != nil {  
		log.Fatalf("Failed to encode users.json: %v", err)  
	}  
}  
  
func loginHandler(w http.ResponseWriter, r *http.Request) {  
	if r.Method == http.MethodGet {  
		tmpl := template.Must(template.ParseFiles("templates/login.html"))  
		tmpl.Execute(w, nil)  
		return  
	}  
  
	if r.Method == http.MethodPost {  
		username := r.FormValue("username")  
		password := r.FormValue("password")  
  
		for _, user := range users {  
			if user.Username == username && user.Password == password {  
				http.Redirect(w, r, "/", http.StatusSeeOther)  
				return  
			}  
		}  
  
		http.Error(w, "Invalid username or password", http.StatusUnauthorized)  
	}  
}  
  
func registerHandler(w http.ResponseWriter, r *http.Request) {  
	if r.Method == http.MethodGet {  
		tmpl := template.Must(template.ParseFiles("templates/register.html"))  
		tmpl.Execute(w, nil)  
		return  
	}  
  
	if r.Method == http.MethodPost {  
		username := r.FormValue("username")  
		password := r.FormValue("password")  
  
		if username == "" || password == "" {  
			http.Error(w, "Username and password are required", http.StatusBadRequest)  
			return  
		}  
  
		// Check for duplicate username  
		for _, user := range users {  
			if user.Username == username {  
				http.Error(w, "Username already exists", http.StatusBadRequest)  
				return  
			}  
		}  
  
		// Check for duplicate password  
		for _, user := range users {  
			if user.Password == password {  
				// Return a response with a JavaScript alert  
				alertMessage := fmt.Sprintf("alert('Maaf password yang Anda masukkan telah digunakan oleh %s');", user.Username)  
				w.Header().Set("Content-Type", "text/html")  
				fmt.Fprintf(w, "<script>%s window.location.href='/register';</script>", alertMessage)  
				return  
			}  
		}  
  
		newUser := User{Username: username, Password: password}  
		users = append(users, newUser)  
		saveUsers()  
  
		http.Redirect(w, r, "/login", http.StatusSeeOther)  
	}  
}  
  
func homeHandler(w http.ResponseWriter, r *http.Request) {  
	if r.URL.Path != "/" {  
		http.NotFound(w, r)  
		return  
	}  
  
	fmt.Fprintf(w, "Welcome to the Home Page!")  
}  
