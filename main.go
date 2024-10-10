package main

import (
	"html/template"
	"net/http"
	"regexp"
	"strconv"
	"sync"
)

type Etudiant struct {
	Nom    string
	Prenom string
	Age    int
	Sexe   string
}

type Classe struct {
	Nom             string
	Filiere         string
	Niveau          string
	NombreEtudiants int
	Etudiants       []Etudiant
}

type UserInfo struct {
	Nom           string
	Prenom        string
	DateNaissance string
	Sexe          string
}

var (
	counter         int
	mu              sync.Mutex
	userInfo        UserInfo
	templatePromo   *template.Template
	templateChange  *template.Template
	templateForm    *template.Template
	templateDisplay *template.Template
	templateError   *template.Template
)

func main() {
	// Chargement des templates
	templatePromo = template.Must(template.ParseFiles("templates/promo.html"))
	templateChange = template.Must(template.ParseFiles("templates/change.html"))
	templateForm = template.Must(template.ParseFiles("templates/form.html"))
	templateDisplay = template.Must(template.ParseFiles("templates/display.html"))
	templateError = template.Must(template.ParseFiles("templates/error.html"))

	// Routes
	http.HandleFunc("/promo", handlePromo)
	http.HandleFunc("/change", handleChange)
	http.HandleFunc("/user/form", handleForm)
	http.HandleFunc("/user/treatment", handleTreatment)
	http.HandleFunc("/user/display", handleDisplay)
	http.HandleFunc("/user/error", handleError)

	// Lancer le serveur
	http.ListenAndServe(":8080", nil)
}

func handlePromo(w http.ResponseWriter, r *http.Request) {
	classe := Classe{
		Nom:             "B1 Informatique",
		Filiere:         "Informatique",
		Niveau:          "Bachelor 1",
		NombreEtudiants: 3,
		Etudiants: []Etudiant{
			{"Dupont", "Jean", 20, "male"},
			{"Martin", "Alice", 22, "female"},
			{"Durand", "Marc", 21, "male"},
		},
	}

	err := templatePromo.Execute(w, classe)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleChange(w http.ResponseWriter, r *http.Request) {
	mu.Lock()
	counter++
	isEven := counter%2 == 0
	message := ""
	if isEven {
		message = "Le nombre de vues est pair : " + strconv.Itoa(counter)
	} else {
		message = "Le nombre de vues est impair : " + strconv.Itoa(counter)
	}
	mu.Unlock()

	data := struct {
		Message string
		IsEven  bool
	}{
		Message: message,
		IsEven:  isEven,
	}

	err := templateChange.Execute(w, data)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleForm(w http.ResponseWriter, r *http.Request) {
	err := templateForm.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleTreatment(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodPost {
		nom := r.FormValue("nom")
		prenom := r.FormValue("prenom")
		dateNaissance := r.FormValue("date_naissance")
		sexe := r.FormValue("sexe")

		// Validation
		if validateInput(nom, prenom) {
			userInfo = UserInfo{
				Nom:           nom,
				Prenom:        prenom,
				DateNaissance: dateNaissance,
				Sexe:          sexe,
			}
			http.Redirect(w, r, "/user/display", http.StatusSeeOther)
			return
		}
		http.Redirect(w, r, "/user/error", http.StatusSeeOther)
	}
	http.Redirect(w, r, "/user/form", http.StatusSeeOther)
}

func handleDisplay(w http.ResponseWriter, r *http.Request) {
	if userInfo.Nom == "" || userInfo.Prenom == "" || userInfo.DateNaissance == "" || userInfo.Sexe == "" {
		http.Redirect(w, r, "/user/error", http.StatusSeeOther)
		return
	}

	err := templateDisplay.Execute(w, userInfo)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func handleError(w http.ResponseWriter, r *http.Request) {
	err := templateError.Execute(w, nil)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}
}

func validateInput(nom, prenom string) bool {
	nameRegex := regexp.MustCompile("^[a-zA-ZÀ-ÿ '-]{1,32}$")
	return nameRegex.MatchString(nom) && nameRegex.MatchString(prenom)
}
