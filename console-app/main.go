// /console-app/main.go
package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/mail" // Parsing email address
	"os"
	"strings"
	"unicode"
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
		if session.Email == "" {
			// If there is no email in the session, show the login/signup options
			switch loginOrSignUp(reader) {
			case "login":
				session = logIn(reader)
			case "signup":
				signUp(reader)
			case "exit":
				fmt.Println("Exiting the application. Goodbye!")
				return
			}
		} else {
			// If there is an email in the session, the user is logged in, show the main menu
			showMainMenu(reader, &session)
		}
	}
}

func loginOrSignUp(reader *bufio.Reader) string {
	// Show options to login, sign up, or exit
	// Return the user's choice
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
			return "signup"
		case "2":
			return "login"
		case "3":
			return "exit"
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
			fmt.Println("6. Sign up to become a Car Owner and Publish your own Trips!")
			fmt.Println("7. Log Out")
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
				fmt.Println("\nLogging out.")
				*session = AppSession{} // Clear the session data
				return                  // Exit the showMainMenu function
			} else {
				fmt.Println("\nGoing to update passenger profile.")
				updateUserProfile(reader, session)
			}
		case "6":
			if session.UserType == "passenger" {
				becomeCarOwner(reader, session)
			}
		case "7":
			if session.UserType == "passenger" {
				fmt.Println("\nLogging out.")
				*session = AppSession{} // Clear the session data
				return                  // Exit the showMainMenu function
			}
		default:
			fmt.Println("\nInvalid choice, please try again.")
			fmt.Println("Press 'Enter' to continue...")
			reader.ReadString('\n')
		}
	}
}

// SECTION 5: UPDATE PROFILE MENU
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
			updateUserMobile(reader, session)
		case "3":
			updateUserEmail(reader, session)
		case "4":
			//updatePassengerCapacity(reader, session)
		case "5":
			//updateDriversLicense(reader, session)
		case "6":
			//updateCarPlate(reader, session)
		case "7":
			deleteAccount(reader, session)
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
			updateUserMobile(reader, session)
		case "3":
			updateUserEmail(reader, session)
		case "4":
			deleteAccount(reader, session)
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
		fmt.Println("\nFirst and last name cannot be empty.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
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

func updateUserMobile(reader *bufio.Reader, session *AppSession) {
	fmt.Print("\nPlease enter your new mobile number: ")
	newMobile, _ := reader.ReadString('\n')
	newMobile = strings.TrimSpace(newMobile)

	// Validate the new mobile number
	if newMobile == "" {
		fmt.Println("\nNew mobile number cannot be empty.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
		return
	} else if !isValidMobileNumber(newMobile) {
		fmt.Println("\nInvalid mobile number. Please enter only numbers.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
		return
	}

	// Prepare the request payload
	updateData := map[string]string{
		"email":  session.Email, // Use the email from the session for user identification
		"mobile": newMobile,
	}

	// Include the mobile number in the update data
	updateData["mobile"] = newMobile

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		fmt.Println("Error marshaling update data:", err)
		return
	}

	// Create a new PATCH request to update the user's mobile number
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/updateMobile", bytes.NewBuffer(jsonData)) // Adjust the URL as needed
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the PATCH request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response status code
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to update mobile number. Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("\nMobile number updated successfully.")
	fmt.Println("Press 'Enter' to return to the main menu...")
	reader.ReadString('\n')
}

// isValidMobileNumber func is a mobile number validation logic
func isValidMobileNumber(number string) bool {
	for _, r := range number {
		if !unicode.IsDigit(r) {
			return false
		}
	}

	return true
}

func updateUserEmail(reader *bufio.Reader, session *AppSession) {
	fmt.Println("\nUpdate Email Address:")
	fmt.Print("\nPlease enter your new email address: ")
	newEmail, _ := reader.ReadString('\n')
	newEmail = strings.TrimSpace(newEmail)

	// Validate the new email
	if !isValidEmail(newEmail) {
		fmt.Println("\nInvalid email format. Please enter a valid email address.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
		return
	}

	// Prepare the request payload
	updateData := map[string]string{
		"old_email": session.Email,
		"new_email": newEmail,
	}

	jsonData, err := json.Marshal(updateData)
	if err != nil {
		fmt.Println("Error marshaling update data:", err)
		return
	}

	// Create the request
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/updateEmail", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	// Execute the request
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request:", err)
		return
	}
	defer resp.Body.Close()

	// Check the response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to update email. Status code: %d\n", resp.StatusCode)
		return
	}

	// Update the session with the new email
	session.Email = newEmail
	fmt.Println("\nEmail address updated successfully.")

	// Prompt to return to main menu
	fmt.Println("Press 'Enter' to return to the main menu...")
	reader.ReadString('\n')
}

// isValidEmail func is a email validation logic.
func isValidEmail(email string) bool {
	_, err := mail.ParseAddress(email)
	return err == nil
}

func updatePassengerCapacity(reader *bufio.Reader, session *AppSession) {
	fmt.Print("Please enter the new number of passengers your car can accommodate: ")
}

func updateDriversLicense(reader *bufio.Reader, session *AppSession) {
	/* ... */
}
func updateCarPlate(reader *bufio.Reader, session *AppSession) {
	/* ... */
}

func deleteAccount(reader *bufio.Reader, session *AppSession) {
	// Confirm account deletion
	fmt.Println("\nAre you sure you want to de-activate your account? This action cannot be undone. (y/n)")
	response, _ := reader.ReadString('\n')
	response = strings.TrimSpace(response)

	// Check if the user confirmed the deletion
	if strings.ToLower(response) != "y" {
		fmt.Println("\nAccount de-activation cancelled.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
		return
	}

	// Send a PATCH request to mark the account as inactive
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/deleteAccount", nil)
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	// Assuming you're sending the email as an identification of the user to be deleted
	q := req.URL.Query()
	q.Add("email", session.Email)
	req.URL.RawQuery = q.Encode()

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request to delete account:", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Failed to delete account. Status code: %d\n", resp.StatusCode)
		return
	}

	fmt.Println("\nAccount marked as inactive successfully.")
	fmt.Println("Thank you for your time with us! We hope to see you again soon")

	// Clear the session data
	*session = AppSession{}

	fmt.Println("Press 'Enter' to return to login menu...")
	reader.ReadString('\n')
	loginOrSignUp(reader)
}

// SECTION 6: Upgrade to Car Owner
func becomeCarOwner(reader *bufio.Reader, session *AppSession) {
	fmt.Print("Please enter your driver's license number: ")
	driversLicense, _ := reader.ReadString('\n')
	driversLicense = strings.TrimSpace(driversLicense)

	fmt.Print("Please enter your car plate number: ")
	carPlateNumber, _ := reader.ReadString('\n')
	carPlateNumber = strings.TrimSpace(carPlateNumber)

	// Validate the input
	if driversLicense == "" || carPlateNumber == "" {
		fmt.Println("\nDriver's license number and car plate number cannot be empty.")
		fmt.Println("Press 'Enter' to return to the main menu...")
		reader.ReadString('\n')
		return
	}

	// Prepare the request payload
	upgradeData := map[string]string{
		"email":            session.Email,
		"drivers_license":  driversLicense,
		"car_plate_number": carPlateNumber,
	}

	jsonData, err := json.Marshal(upgradeData)
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}

	resp, err := http.Post("http://localhost:5000/upgradeToCarOwner", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode == http.StatusOK {
		fmt.Println("\nSuccessfully upgraded to car owner.")
		session.UserType = "car_owner"
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("\nFailed to upgrade to car owner. Status code: %d, Message: %s\n", resp.StatusCode, string(body))
	}

	fmt.Println("Press 'Enter' to return to the main menu...")
	reader.ReadString('\n')
	showMainMenu(reader, session)
}
