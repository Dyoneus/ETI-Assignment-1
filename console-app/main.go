// /console-app/main.go
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
	Email     string
	UserType  string
	FirstName string
	LastName  string
	Password  string
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
				updateCarOwnerProfile(reader, session)
			} else {
				//viewEnrolledTrips(reader)
			}
		case "4":
			//viewPastTrips(reader)
		case "5":
			if session.UserType == "car_owner" {
				fmt.Println("\nReturning to the main menu.")
				return
			} else {
				fmt.Println("\nGoing to update passenger profile.")
				updateUserProfile(reader, session)
			}
		case "6":
			return
		default:
			fmt.Println("\nInvalid choice, please try again.")
			fmt.Println("Press 'Enter' to continue...")
			reader.ReadString('\n')
		}
	}
}

// SECTION 1: UPDATE PROFILE MENU
func updateCarOwnerProfile(reader *bufio.Reader, session *AppSession) {
	for {
		fmt.Println("\nWhat would you like to update?")
		fmt.Println("1. Name")
		fmt.Println("2. Mobile Number")
		fmt.Println("3. Email Address")
		fmt.Println("4. Number of Passengers")
		fmt.Println("5. Driver's License Number")
		fmt.Println("6. Car Plate Number")
		fmt.Println("7. Delete Account")
		fmt.Println("8. Return to Main Menu")

		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			updateUserName(reader, session)
		case "2":
			//updateUserMobile(reader, session)
		case "3":
			//updateUserEmail(reader, session)
		case "4":
			//updatePassengerCapacity(reader, session)
		case "5":
			//updateDriversLicense(reader, session)
		case "6":
			//updateCarPlate(reader, session)
		case "7":
			//deleteAccount(reader, session)
			return // Return to main menu after deleting the account
		case "8":
			return // Return to main menu
		default:
			fmt.Println("\nInvalid choice, please try again.")
			fmt.Println("Press 'Enter' to continue...")
			reader.ReadString('\n')
		}
	}
}

func updateUserProfile(reader *bufio.Reader, session *AppSession) {
	for {
		fmt.Println("\nWhat would you like to update?")
		fmt.Println("1. Name")
		fmt.Println("2. Mobile Number")
		fmt.Println("3. Email Address")
		fmt.Println("4. Delete Account")
		fmt.Println("5. Return to Main Menu")

		fmt.Print("\nEnter your choice: ")
		choice, _ := reader.ReadString('\n')
		choice = strings.TrimSpace(choice)

		switch choice {
		case "1":
			updateUserName(reader, session)
		case "2":
			//updateUserMobile(reader, session)
		case "3":
			//updateUserEmail(reader, session)
		case "4":
			//deleteAccount(reader, session)
			return // Return to main menu after deleting the account
		case "5":
			return // Return to main menu
		default:
			fmt.Println("\nInvalid choice, please try again.")
			fmt.Println("Press 'Enter' to continue...")
			reader.ReadString('\n')
		}
	}
}

// Functions for updating the profile
func updateUserName(reader *bufio.Reader, session *AppSession) {
	// Prompt the user for new first name
	fmt.Print("Please enter your new first name: ")
	newFirstName, _ := reader.ReadString('\n')
	newFirstName = strings.TrimSpace(newFirstName)

	// Prompt the user for new last name
	fmt.Print("Please enter your new last name: ")
	newLastName, _ := reader.ReadString('\n')
	newLastName = strings.TrimSpace(newLastName)

	// Validate the input
	if newFirstName == "" || newLastName == "" {
		fmt.Println("First and last name cannot be empty.")
		return
	}

	updateData := map[string]string{
		"email":      session.Email, // assuming session.Email contains the user's email
		"first_name": newFirstName,
		"last_name":  newLastName,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		fmt.Println("Error marshaling update data:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/users", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to update name. Status code: %d\n", resp.StatusCode)
		return
	}

	// Successfully updated the name
	fmt.Println("Your name has been updated successfully.")

	// Update the session with the new name
	session.FirstName = newFirstName
	session.LastName = newLastName

	fmt.Println("\nPress 'Enter' to return to the main menu...")
	reader.ReadString('\n')
}

func updateUserMobile(reader *bufio.Reader, session *AppSession) { /* ... */ }
func updateUserEmail(reader *bufio.Reader, session *AppSession)  { /* ... */ }

func updatePassengerCapacity(reader *bufio.Reader, session *AppSession) {
	fmt.Print("Please enter the new number of passengers your car can accommodate: ")
}

func updateDriversLicense(reader *bufio.Reader, session *AppSession) { /* ... */ }
func updateCarPlate(reader *bufio.Reader, session *AppSession)       { /* ... */ }

func deleteAccount(reader *bufio.Reader, session *AppSession) {
	// Confirm deletion and send a DELETE request to the user-service
}
