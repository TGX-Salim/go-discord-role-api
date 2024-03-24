package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"github.com/gorilla/mux"
	"github.com/bwmarrin/discordgo"
)

type Response struct {
	Message string `json:"message"`
	Error bool `json:"error"`
}

func main() {
	apiSecret := os.Getenv("API_SECRET")

	r := mux.NewRouter()
	r.Use(apiKeyMiddleware(apiSecret))

	r.HandleFunc("/api/role/add/{userID}/{roleID}", addRoleHandler)

	r.HandleFunc("/api/role/remove/{userID}/{roleID}", removeRoleHandler)

	r.NotFoundHandler = http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		response := Response{Message: "404 Not Found", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusNotFound)
		w.Write(jsonResponse)
	})

	log.Println("Role API by @akatiggerx04 | Starting server on :2209...")
	log.Fatal(http.ListenAndServe(":2209", r))
}

func apiKeyMiddleware(apiSecret string) mux.MiddlewareFunc {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			apiKey := r.Header.Get("X-API-Key")
			w.Header().Set("X-Powered-By", "@akatiggerx04")
			if apiKey != apiSecret {
				log.Printf("Unauthorized access, invalid API key: %s\n", apiKey)
				response := Response{Message: "Unauthorized access, invalid API key!", Error: true}
				jsonResponse, _ := json.Marshal(response)
				w.Header().Set("Content-Type", "application/json")
				w.WriteHeader(http.StatusUnauthorized)
				w.Write(jsonResponse)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}

func addRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	roleID := vars["roleID"]

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Printf("Error creating Discord session: %s\n", err)
		response := Response{Message: "Error creating Discord session", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}
	defer session.Close()

	err = session.Open()
	if err != nil {
		log.Printf("Error opening connection to Discord: %s\n", err)
		response := Response{Message: "Error opening connection to Discord", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	err = session.GuildMemberRoleAdd(os.Getenv("GUILD_ID"), userID, roleID)
	if err != nil {
		log.Printf("Error adding role to user: %s\n", err)
		response := Response{Message: "Error adding role to user", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	log.Printf("Added role %s from user ID: %s\n", roleID, userID)
	response := Response{Message: "Role added successfully", Error: false}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}

func removeRoleHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	userID := vars["userID"]
	roleID := vars["roleID"]

	session, err := discordgo.New("Bot " + os.Getenv("DISCORD_TOKEN"))
	if err != nil {
		log.Printf("Error creating Discord session: %s\n", err)
		response := Response{Message: "Error creating Discord session", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}
	defer session.Close()

	err = session.Open()
	if err != nil {
		log.Printf("Error opening connection to Discord: %s\n", err)
		response := Response{Message: "Error opening connection to Discord", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	err = session.GuildMemberRoleRemove(os.Getenv("GUILD_ID"), userID, roleID)
	if err != nil {
		log.Printf("Error removing role from user: %s\n", err)
		response := Response{Message: "Error removing role from user", Error: true}
		jsonResponse, _ := json.Marshal(response)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	log.Printf("Removed role %s from user ID: %s\n", roleID, userID)
	response := Response{Message: "Role removed successfully", Error: false}
	jsonResponse, _ := json.Marshal(response)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(jsonResponse)
}