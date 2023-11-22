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

// Simulate the idea of a session retaining the user's information
type AppSession struct {
	Email    string
	UserType string
}

func main() {
	reader := bufio.NewReader(os.Stdin)
	var session AppSession // A variable to hold the session information

	for {
		fmt.Println("\n===========================================================")
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
			session = logIn(reader) // Capture the session information after logging in
			if session.Email != "" {
				showMainMenu(reader, &session)
			}
		case "3":
			fmt.Println("\nExiting the application. Goodbye!")
			return
		default:
			fmt.Println("\nInvalid choice, please try again.")
			fmt.Println("Press 'Enter' to return to the main menu...")
			reader.ReadString('\n') // Pause the program
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

	// Before marshaling and sending the data, validate the input
	if firstName == "" || lastName == "" || mobile == "" || email == "" || password == "" {
		fmt.Println("\nAll fields are required and cannot be empty. Please try again.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n') // Pause the program
		return
	}

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

func logIn(reader *bufio.Reader) (session AppSession) {
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

		// Store user's email and type in the session
		session = AppSession{
			Email:    email,
			UserType: loginResponse.UserType,
		}
		showMainMenu(reader, &session)
	} else {
		fmt.Println("\nLogin failed. Please check your credentials and try again.")
	}

	// Prompt to press 'Enter' and read it to pause the program
	fmt.Println("\nPress 'Enter' to continue...")
	reader.ReadString('\n')

	return session
}

func showMainMenu(reader *bufio.Reader, session *AppSession) {
	for {
		if session.UserType == "car_owner" {
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
			fmt.Println("5. Update Profile")
			fmt.Println("6. Log Out")
		}

		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			if session.UserType == "car_owner" {
				//publishTrip(reader)
			} else {
				//browseTrips(reader)
			}
		case "2":
			if session.UserType == "car_owner" {
				//manageTrips(reader)
			} else {
				//enrollInTrip(reader)
			}
		case "3":
			if session.UserType == "car_owner" {
				updateCarOwnerProfile(reader)
			} else {
				//viewEnrolledTrips(reader)
			}
		case "4":
			//viewPastTrips(reader)
		case "5":
			if session.UserType == "car_owner" {
				return
			} else {
				updateUserProfile(reader)
			}
		case "6":

			return
		default:
			fmt.Println("\nInvalid choice, please try again.")
		}
	}
}

// SECTION 1: UPDATE PROFILE MENU
func updateCarOwnerProfile(reader *bufio.Reader) {
	fmt.Println("\nUpdate Car Owner Profile:")
	// Prompt the user to update number of passengers car can accomodate
}

func updateUserProfile(reader *bufio.Reader) {
	fmt.Println("\nUpdate User Profile:")
	// Prompt the user name, mobile number, and email changes.
	// An option to delete the account if it's older than 1 year.
}
