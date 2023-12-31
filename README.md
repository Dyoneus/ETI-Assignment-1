# Car-Pooling Service Platform
### Front-end
![Home Page](web-frontend/images/home.jpg)

### Console Application
![Platform on Console](web-frontend/images/consolehome.png)

## Introduction
This project implements a car-pooling service platform using microservices architecture. The application facilitates users to either offer car-pooling services (as car owners) or enroll in available trips (as passengers). The backend is developed in Go, interfacing with a MySQL database, and the front-end is a web application using HTML, CSS, and JavaScript. 

## Design Considerations

### Microservices
- **User Service**: Handles user authentication, profile management, and user-type upgrades.
- **Trip Service**: Manages trip creation, enrollment, and trip history.
- **Technologies**: Go for backend services, MySQL for persistent storage.
- **Communication**: RESTful APIs for microservice interaction.

### Frontend
- **Web Application**: Developed using HTML, CSS, and JavaScript.
- **Interaction**: Communicates with backend microservices via AJAX calls.

## Key Features

- **Password Hashing**: Implements bcrypt hashing for secure password storage, ensuring user credentials are protected.
- **Automated Trip Management**: The system automatically monitors trip schedules. Trips past their start time are automatically flagged and removed from active listings, ensuring data accuracy and relevancy.
- **User Type Flexibility**: Users can sign up as passengers and later upgrade to car owners.

## Architecture Diagram
![Architecture Diagram](web-frontend/images/architectureDiagram.jpg)


## Database Configuration
- Create an account named user with the password ‘password’ is created. This account is granted all permission on the database
- Ensure MySQL is running and configured by creating the 2 database (`carpool` and `carpool_trips`)
- Backend services will auto create the table in the 2 database created.
- SQL statement script is in the same directory as README.md

## Setup and Running Instructions
Note: Make sure to setup Database Configuration before continuing the setup here.

1. **Backend Services**:
   - Navigate to each service directory (`user-service` and `trip-service`) in Terminal/CMD.
   - Run `go run main.go` to start each service.
   - For e.g. (Terminal 1: `\user-service> go run main.go`) & (Terminal 2: `\trip-service> go run main.go`)
2. **Frontend**:
   - Open the `index.html` from the `web-frontend` directory in a browser.
   - Ensure backend services (`user-service\main.go` and `trip-service\main.go`) are running for full functionality.
3. **Console-App**:
   - Navigate to directory (`console-app`) in Terminal/CMD.
   - Run `go run main.go` to start the console application UI.
   - Ensure backend services are running for full functionality.
   - For e.g. (Terminal 1: `\console-app> go run main.go`), (Terminal 2: `\user-service> go run main.go`) & (Terminal 3: `\trip-service> go run main.go`)



---

By [Ong Jia Yuan]








# Car-Pooling Service Platform - Task List

1. User Account Creation
- Status: Completed
- Details: Implementation allows both passengers and car owners to create accounts.

2. Default Passenger Profile Creation
- Status: Completed
- Details: Users provide first name, last name, mobile number, and email address.

3. Car Owner Profile Enhancement
- Status: Completed
- Details: Passenger profiles can be upgraded to car owner profiles with additional details like driver’s license and car plate number.

4. Account Information Update
- Status: Completed
- Details: Users can update their account information.

5. Account Deletion Post One Year
- Status: Completed
- Details: Users can delete their accounts after 1 year, adhering to data retention policies.

6. Trip Publishing by Car Owners
- Status: Completed
- Details: Car owners can publish trips with detailed information including pick-up, alternative pick-up locations, start time, destination, and passenger capacity.

7. Trip Enrollment by Passengers
- Status: Completed
- Details: Passengers can browse, search, and enroll in available trips with seat availability and no schedule conflicts.

8. Trip Management by Car Owners
- Status: Completed
- Details: Car owners can only start 30 minutes in the future of scheduled time

10. Retrieve Past Trips
- Status: Completed
- Details: Users can access their past trips in reverse chronological order.

## Additional Features

1. Password Hashing
- Status: Completed
- Details: Passwords are securely hashed during account creation and login.

2. Automatic Trip Deletion
- Status: Completed
- Details: Past trips are automatically deleted from the system after their scheduled start time.

3. Cancel Enrollment (Front-end)
- Status: Completed
- Details: Passengers can cancel their enrollment trip but will not be able to enroll in the same trip.

4. User Authentication (Front-end)
- Status: Completed
- Details: User who tries to enter the main menu page manually in the URL will automatically redirect to login page.

This document will be updated as the project progresses, reflecting new implementations and enhancements.


