package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strings"
)

func main() {
	reader := bufio.NewReader(os.Stdin)

	for {
		fmt.Println("===========================================================")
		fmt.Println("Welcome to the Car-Pooling Service Console Application!")
		fmt.Println("===========================================================")
		fmt.Println("\nPlease select an option:")
		fmt.Println("1. Sign Up (Create an account)")
		fmt.Println("2. Log In (Existing users)")
		fmt.Println("3. Exit")
		fmt.Print("\nEnter your choice (1-3): ")

		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			signUp(reader)
		case "2":
			logIn(reader)
		case "3":
			fmt.Println("Exiting the application. Goodbye!")
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}

func signUp(reader *bufio.Reader) {
	fmt.Println("\n==========================")
	fmt.Println("Sign Up for New Account")
	fmt.Println("==========================\n")

	fmt.Print("Please enter your first name: ")
	firstName, _ := reader.ReadString('\n')
	firstName = strings.TrimSpace(firstName)

	fmt.Print("Please enter your last name: ")
	lastName, _ := reader.ReadString('\n')
	lastName = strings.TrimSpace(lastName)

	fmt.Print("Please enter your mobile number: ")
	mobile, _ := reader.ReadString('\n')
	mobile = strings.TrimSpace(mobile)

	fmt.Print("Please enter your email address: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Please enter your password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	user := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
		"mobile":     mobile,
		"email":      email,
		"password":   password,
	}

	userData, err := json.Marshal(user)
	if err != nil {
		fmt.Println("Error encoding user data:", err)
		return
	}

	response, err := http.Post("http://localhost:5000/users", "application/json", strings.NewReader(string(userData)))
	if err != nil {
		fmt.Println("Error calling signup service:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusCreated {
		fmt.Println("\nSignup successful!")
	} else {
		fmt.Println("\nSignup failed!")
	}

	fmt.Println("\nPress 'Enter' to return to the main menu...")
	reader.ReadString('\n')
}

func logIn(reader *bufio.Reader) {
	// Logic to handle user login
}
