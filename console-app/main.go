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
	"strconv"
	"strings"
	"time"
	"trip-service/models"
	"unicode"
)

// Simulate the idea of a session retaining the user's information
type AppSession struct {
	UserID    uint
	Email     string
	UserType  string
	FirstName string
	LastName  string
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
			UserID    uint   `json:"userID"`
			UserType  string `json:"userType"`
			FirstName string `json:"first_name"`
			LastName  string `json:"last_name"`
		}
		// Unmarshal the JSON response into the loginResponse struct
		err = json.Unmarshal(responseBody, &loginResponse)
		if err != nil {
			fmt.Println("Error unmarshaling response:", err)
			return
		}

		// Store user's email and type in the session
		session = AppSession{
			UserID:    loginResponse.UserID, // Store the UserID in the session
			Email:     email,
			UserType:  loginResponse.UserType,
			FirstName: loginResponse.FirstName,
			LastName:  loginResponse.LastName,
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
				publishTrip(reader, session)
			} else {
				//browseTrips(reader)
			}
		case "2":
			if session.UserType == "car_owner" {
				manageTrips(reader, session)
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
		fmt.Println("4. Driver's License Number")
		fmt.Println("5. Car Plate Number")
		fmt.Println("6. Delete Account")
		fmt.Println("7. Return to Main Menu")

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
			updateDriversLicense(reader, session)
		case "5":
			updateCarPlate(reader, session)
		case "6":
			deleteAccount(reader, session)
			return // Return to main menu after deleting the account
		case "7":
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

func updateDriversLicense(reader *bufio.Reader, session *AppSession) {
	fmt.Print("\nPlease enter your new driver's license number: ")
	newLicense, _ := reader.ReadString('\n')
	newLicense = strings.TrimSpace(newLicense)

	// Add input validation here

	// Prepare the request data
	requestData := map[string]string{
		"email":           session.Email,
		"drivers_license": newLicense,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling request data:", err)
		return
	}

	// Send the update request to the server
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/updateDriversLicense", bytes.NewBuffer(jsonData))
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

	// Handle the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Driver's license number updated successfully.")
	} else {
		fmt.Printf("Failed to update driver's license number. Status code: %d\n", resp.StatusCode)
	}

	// Read the 'Enter' key to pause
	fmt.Println("Press 'Enter' to continue...")
	reader.ReadString('\n')
}

func updateCarPlate(reader *bufio.Reader, session *AppSession) {
	fmt.Print("\nPlease enter your new car plate number: ")
	newPlate, _ := reader.ReadString('\n')
	newPlate = strings.TrimSpace(newPlate)

	// Add input validation here

	// Prepare the request data to match the JSON fields expected by the server
	requestData := map[string]string{
		"email":            session.Email,
		"car_plate_number": newPlate,
	}

	jsonData, err := json.Marshal(requestData)
	if err != nil {
		fmt.Println("Error marshaling request data:", err)
		return
	}

	// Send the update request to the server
	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5000/updateCarPlate", bytes.NewBuffer(jsonData))
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

	// Handle the response
	if resp.StatusCode == http.StatusOK {
		fmt.Println("Car plate number updated successfully.")
	} else {
		var response map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&response); err != nil {
			fmt.Println("Error decoding response data:", err)
		} else {
			fmt.Printf("Failed to update car plate number. Status code: %d, Message: %s\n", resp.StatusCode, response["message"])
		}
	}

	// Read the 'Enter' key to pause
	fmt.Println("Press 'Enter' to continue...")
	reader.ReadString('\n')
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

// SECTION 7: Publish Trips (For Car Owner)
func publishTrip(reader *bufio.Reader, session *AppSession) {
	if session.UserType != "car_owner" {
		fmt.Println("\nOnly car owners can publish trips.")
		return
	}

	fmt.Println("\nEnter trip details to publish your trip:")

	// Capture pick-up location
	fmt.Print("Pick Up Location: ")
	pickUpLocation, _ := reader.ReadString('\n')
	pickUpLocation = strings.TrimSpace(pickUpLocation)

	// Capture alternative pick-up location
	fmt.Print("Alternative Pick Up Location (if any, press enter to skip): ")
	alternativePickUp, _ := reader.ReadString('\n')
	alternativePickUp = strings.TrimSpace(alternativePickUp)

	// Capture start traveling time
	fmt.Print("Start Traveling Time (format DD-MM-YYYY HH:MM): ")
	travelStartTime, _ := reader.ReadString('\n')
	travelStartTime = strings.TrimSpace(travelStartTime)
	// Define the layout string to match input format
	const layout = "02-01-2006 15:04"
	// Parse the travelStartTime using the specified layout
	loc, _ := time.LoadLocation("Asia/Singapore")
	parsedTime, err := time.ParseInLocation(layout, travelStartTime, loc)
	if err != nil {
		fmt.Println("Invalid travel start time format:", err)
		return
	}
	if time.Until(parsedTime) < 30*time.Minute {
		fmt.Println("Travel start time must be at least 30 minutes in the future.")
		return
	}

	// Capture destination address
	fmt.Print("Destination Address: ")
	destinationAddress, _ := reader.ReadString('\n')
	destinationAddress = strings.TrimSpace(destinationAddress)

	// Capture number of passengers car can accommodate
	fmt.Print("Number of Passengers Car Can Accommodate: ")
	availableSeatsStr, _ := reader.ReadString('\n')
	availableSeatsStr = strings.TrimSpace(availableSeatsStr)
	// Convert availableSeats from string to int
	availableSeats, err := strconv.Atoi(availableSeatsStr)
	if err != nil {
		fmt.Println("Invalid number of seats. Please enter a valid number.")
		return
	}

	// Validate input data
	if !isValidTripInput(pickUpLocation, alternativePickUp, travelStartTime, destinationAddress, availableSeats) {
		fmt.Println("Invalid input. Please try again.")
		return
	}

	// Convert car_owner_id to uint before adding it to the payload
	carOwnerID := uint(session.UserID)
	// Prepare the request payload
	tripData := map[string]interface{}{
		"car_owner_id":        carOwnerID,
		"car_owner_name":      session.FirstName + " " + session.LastName,
		"pick_up_location":    pickUpLocation,
		"alternative_pick_up": alternativePickUp,
		"travel_start_time":   parsedTime,
		"destination_address": destinationAddress,
		"available_seats":     availableSeats,
	}

	jsonData, err := json.Marshal(tripData)
	if err != nil {
		fmt.Println("Error marshaling trip data:", err)
		return
	}

	// Send the POST request
	resp, err := http.Post("http://localhost:5001/trips", "application/json", bytes.NewBuffer(jsonData))
	if err != nil {
		fmt.Println("Error sending trip publication request:", err)
		return
	}
	defer resp.Body.Close()

	// Handle the response
	if resp.StatusCode == http.StatusCreated {
		fmt.Println("Trip published successfully.")
	} else {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to publish trip. Status code: %d, Message: %s\n", resp.StatusCode, string(body))
	}
}

// isValidTripInput checks if the provided trip details are valid.
func isValidTripInput(pickUpLocation, alternativePickUp, travelStartTime, destinationAddress string, availableSeats int) bool {
	// Check if the required strings are not empty
	if pickUpLocation == "" || destinationAddress == "" {
		return false
	}

	// If an alternative pick-up is provided, ensure it's not empty
	if alternativePickUp != "" && strings.TrimSpace(alternativePickUp) == "" {
		return false
	}

	/*
		// Validate travel start time
		travelTime, err := time.Parse("2006-01-02 15:04:05", travelStartTime)
		if err != nil {
			fmt.Println("Invalid travel start time format:", err)
			return false
		}

		if time.Until(travelTime) < 30*time.Minute {
			fmt.Println("Travel start time must be at least 30 minutes in the future.")
			return false
		}
	*/

	// Validate available seats
	if availableSeats <= 0 {
		fmt.Println("Number of available seats must be a positive number.")
		return false
	}

	return true
}

// SECTION 8: Manage Trips (For Car Owner)
func manageTrips(reader *bufio.Reader, session *AppSession) {
	fmt.Println("\nManage Your Trips:")
	trips := listTrips(session.UserID)
	if trips == nil {
		fmt.Println("No trips found or there was an error retrieving your trips.")
		return // Go back to the main menu
	}

	// 1. List Trips
	for i, trip := range trips {
		fmt.Printf("%d: %s to %s at %s\n", i+1, trip.PickUpLocation, trip.DestinationAddress, trip.TravelStartTime.Format("02-01-2006 15:04"))
	}

	// 2. Ask for which trip to manage
	fmt.Print("\nEnter the number of the trip you want to manage or '0' to return: ")
	input, _ := reader.ReadString('\n')
	input = strings.TrimSpace(input)
	choice, _ := strconv.Atoi(input)

	if choice == 0 {
		return // Go back to the main menu
	}

	selectedTrip := trips[choice-1]

	// 3. Edit or Delete Trip
	fmt.Println("\n1. Edit Trip")
	fmt.Println("2. Delete Trip")
	fmt.Print("\nEnter your choice: ")
	action, _ := reader.ReadString('\n')
	action = strings.TrimSpace(action)

	switch action {
	case "1":
		editTrip(reader, session, selectedTrip) // Implement this function to edit the trip
	case "2":
		deleteTrip(selectedTrip.ID) // Implement this function to delete the trip
	default:
		fmt.Println("Invalid choice.")
	}
}

func listTrips(userID uint) []models.Trip {
	// Log the URL for debugging purposes
	url := fmt.Sprintf("http://localhost:5001/trips?carOwnerID=%d", userID)
	//1fmt.Printf("Fetching trips from URL: %s\n", url)

	// Make an HTTP GET request to the server's ListTrips endpoint
	resp, err := http.Get(url)
	if err != nil {
		fmt.Printf("Error fetching trips: %v\n", err)
		return nil // Return nil to indicate an error occurred
	}
	defer resp.Body.Close()

	// Check if the status code is not 'StatusOK'
	if resp.StatusCode != http.StatusOK {
		fmt.Printf("Error fetching trips, server responded with status code: %d\n", resp.StatusCode)
		return nil // Return nil to indicate an error occurred
	}

	// Decode the JSON response into a slice of Trip structs
	var trips []models.Trip
	if err := json.NewDecoder(resp.Body).Decode(&trips); err != nil {
		fmt.Printf("Error decoding trips response: %v\n", err)
		return nil // Return nil to indicate an error occurred
	}

	return trips
}

func editTrip(reader *bufio.Reader, session *AppSession, trip models.Trip) {
	fmt.Println("Editing trip:", trip.ID)
	fmt.Println("Leave input empty if you do not wish to change the value.")

	// Capture pick-up location
	fmt.Print("New Pick Up Location (current: " + trip.PickUpLocation + "): ")
	pickUpLocation, _ := reader.ReadString('\n')
	pickUpLocation = strings.TrimSpace(pickUpLocation)

	if pickUpLocation != "" {
		trip.PickUpLocation = pickUpLocation
	}

	// Capture alternative pick-up location
	fmt.Print("New Alternative Pick Up Location (current: " + trip.AlternativePickUp + "): ")
	alternativePickUp, _ := reader.ReadString('\n')
	alternativePickUp = strings.TrimSpace(alternativePickUp)

	if alternativePickUp != "" {
		trip.AlternativePickUp = alternativePickUp
	}

	// Capture start traveling time
	fmt.Print("New Start Traveling Time (format DD-MM-YYYY HH:MM, current: " + trip.TravelStartTime.Format("02-01-2006 15:04") + "): ")
	travelStartTime, _ := reader.ReadString('\n')
	travelStartTime = strings.TrimSpace(travelStartTime)

	if travelStartTime != "" {
		const layout = "02-01-2006 15:04"
		parsedTime, err := time.Parse(layout, travelStartTime)
		if err != nil {
			fmt.Println("Invalid travel start time format:", err)
			return
		}
		trip.TravelStartTime = parsedTime
	}

	// Capture destination address
	fmt.Print("New Destination Address (current: " + trip.DestinationAddress + "): ")
	destinationAddress, _ := reader.ReadString('\n')
	destinationAddress = strings.TrimSpace(destinationAddress)

	if destinationAddress != "" {
		trip.DestinationAddress = destinationAddress
	}

	// Capture number of passengers car can accommodate
	fmt.Print("New Number of Passengers Car Can Accommodate (current: " + strconv.Itoa(trip.AvailableSeats) + "): ")
	availableSeatsStr, _ := reader.ReadString('\n')
	availableSeatsStr = strings.TrimSpace(availableSeatsStr)

	if availableSeatsStr != "" {
		availableSeats, err := strconv.Atoi(availableSeatsStr)
		if err != nil {
			fmt.Println("Invalid number of seats. Please enter a valid number.")
			return
		}
		trip.AvailableSeats = availableSeats
	}

	// Send the updated trip data to the backend
	updatedTripData, err := json.Marshal(trip)
	if err != nil {
		fmt.Println("Error marshaling updated trip data:", err)
		return
	}

	req, err := http.NewRequest(http.MethodPatch, "http://localhost:5001/trips/"+strconv.Itoa(int(trip.ID)), bytes.NewBuffer(updatedTripData))
	if err != nil {
		fmt.Println("Error creating request to update trip:", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println("Error sending request to update trip:", err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		body, _ := ioutil.ReadAll(resp.Body)
		fmt.Printf("Failed to update trip. Status code: %d, Message: %s\n", resp.StatusCode, string(body))
		return
	}

	fmt.Println("\nTrip updated successfully.")
}

func deleteTrip(tripID uint) {
	// HTTP DELETE request to the server's DeleteTrip endpoint
	req, err := http.NewRequest("DELETE", fmt.Sprintf("http://localhost:5001/trips/%d", tripID), nil)
	if err != nil {
		// handle error
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		// handle error
	}

	fmt.Println("\nTrip deleted successfully.")
}
