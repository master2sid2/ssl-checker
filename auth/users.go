package auth

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"regexp"

	"golang.org/x/crypto/bcrypt"
)

type User struct {
	Username string `json:"username"`
	Password string `json:"password"`
	Role     string `json:"role"`
}

var Users = make(map[string]User)

func ValidateUsername(username string) bool {
	// Разрешены только буквы, цифры, дефисы, подчеркивания и точки
	validUsername := regexp.MustCompile(`^[a-zA-Z0-9._-]+$`)
	return validUsername.MatchString(username)
}

// Проверка, существует ли пользователь
func UserExists(username string) bool {
	_, exists := Users[username]
	return exists
}

func LoadUsers() {
	data, err := ioutil.ReadFile("data/users.json")
	if err != nil {
		log.Println("Error reading users file:", err)
		CreateDefaultUsers()
		SaveUsers()
		return
	}
	if len(data) == 0 {
		CreateDefaultUsers()
		SaveUsers()
	} else {
		err = json.Unmarshal(data, &Users)
		if err != nil {
			log.Println("Error unmarshalling users:", err)
		}
	}
}

func SaveUsers() {
	data, err := json.MarshalIndent(Users, "", "    ")
	if err != nil {
		log.Println("Error marshalling users:", err)
		return
	}
	err = ioutil.WriteFile("data/users.json", data, 0644)
	if err != nil {
		log.Println("Error writing users file:", err)
	}
}

func CreateDefaultUsers() {
	if _, exists := Users["admin"]; !exists {
		Users["admin"] = User{
			Username: "admin",
			Password: string(HashPassword("admin")),
			Role:     "admin",
		}
	}
}

func HashPassword(password string) []byte {
	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return hashedPassword
}

func CheckPassword(username, password string) bool {
	user, exists := Users[username]
	if !exists {
		log.Println("User does not exist:", username)
		return false
	}
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil {
		log.Println("Password mismatch for user:", username)
	}
	return err == nil
}
