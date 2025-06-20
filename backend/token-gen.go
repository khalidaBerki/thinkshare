// Commande pour générer un token JWT de test
// Exécuter avec: go run token-gen.go
package main

import (
	"backend/internal/auth"
	"fmt"
	"os"
	"strconv"
)

func main() {
	// UserID par défaut (1)
	userID := 1

	// Si un argument est fourni, on l'utilise comme userID
	if len(os.Args) > 1 {
		id, err := strconv.Atoi(os.Args[1])
		if err == nil && id > 0 {
			userID = id
		}
	}

	// Génération du token
	token, err := auth.GenerateJWT(userID)
	if err != nil {
		fmt.Printf("Erreur: %v\n", err)
		os.Exit(1)
	}

	fmt.Println("====================================")
	fmt.Printf("TOKEN JWT POUR UTILISATEUR ID %d:\n", userID)
	fmt.Println("====================================")
	fmt.Println(token)
	fmt.Println("====================================")
	fmt.Println("Pour utiliser avec curl:")
	fmt.Printf("curl -H \"Authorization: Bearer %s\" http://localhost:8080/api/test-auth\n", token)
	fmt.Println("====================================")
	fmt.Println("Pour Postman - Header:")
	fmt.Println("Authorization: Bearer " + token)
}
