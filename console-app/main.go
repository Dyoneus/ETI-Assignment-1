package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
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

	// Add logic to call the microservice to create an account
	userData := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
		"mobile":     mobile,
		"email":      email,
		"password":   password,
	}

	jsonData, err := json.Marshal(userData)
	if err != nil {
		fmt.Println("Error marshaling user data:", err)
		return
	}

	resp, err := http.Post("http://localhost:5000/users", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error calling signup service:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusCreated {
		fmt.Println("\nAccount created successfully!")
	} else {
		fmt.Println("\nFailed to create account. Status code:", resp.StatusCode)
	}

	fmt.Println("\nPress 'Enter' to return to the main menu...")
	reader.ReadString('\n')
}

func logIn(reader *bufio.Reader) {
	fmt.Println("\n============")
	fmt.Println("Log In")
	fmt.Println("============\n")

	fmt.Print("Please enter your email address: ")
	email, _ := reader.ReadString('\n')
	email = strings.TrimSpace(email)

	fmt.Print("Please enter your password: ")
	password, _ := reader.ReadString('\n')
	password = strings.TrimSpace(password)

	// Add logic to call the microservice to create an account
	loginData := map[string]string{
		"email":    email,
		"password": password,
	}

	jsonData, err := json.Marshal(loginData)
	if err != nil {
		fmt.Println("Error marshaling login data:", err)
		return
	}

	resp, err := http.Post("http://localhost:5000/login", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error calling login service:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("\nLogin successful!")
		fmt.Println("\nPress 'Enter' to proceed to the main menu...")

		// Read the response body
		responseBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var loginResponse struct {
			UserType string `json:"userType"`
		}

		// Unmarshal the JSON response into the loginResponse struct
		err = json.Unmarshal(responseBody, &loginResponse)
		if err != nil {
			fmt.Println("Error unmarshaling response:", err)
			return
		}

		showMainMenu(reader, loginResponse.UserType)
	} else {
		fmt.Println("\nLogin failed. Please check your credentials and try again.")
		fmt.Println("\nPress 'Enter' to return to the login screen...")
	}
	// Prompt to press 'Enter' and read it to pause the program
	fmt.Println("\nPress 'Enter' to continue...")
	reader.ReadString('\n')
}

func showMainMenu(reader *bufio.Reader, userType string) {
	for {
		if userType == "car_owner" {
			fmt.Println("\nCar Owner Menu:")
			fmt.Println("1. Publish a Trip")
			fmt.Println("2. Manage Trips")
			fmt.Println("3. Update Profile")
			fmt.Println("4. View Past Trips")
			fmt.Println("5. Log Out")
		} else { // Default to passenger menu
			fmt.Println("\nPassenger Menu:")
			fmt.Println("1. Browse Trips")
			fmt.Println("2. Enroll in a Trip")
			fmt.Println("3. View Enrolled Trips")
			fmt.Println("4. View Past Trips")
			fmt.Println("5. Log Out")
		}

		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if userType == "car_owner" {
				//publishTrip(reader)
			} else {
				//browseTrips(reader)
			}
		case "2":
			if userType == "car_owner" {
				//manageTrips(reader)
			} else {
				//enrollInTrip(reader)
			}
		case "3":
			//updateProfile(reader)
		case "4":
			//viewPastTrips(reader)
		case "5":
			return
		default:
			fmt.Println("Invalid choice, please try again.")
		}
	}
}
