//  /web-frontend/js/loadUserMenu.js

function loadUserMenu() {
    // Fetch the user type from session storage or an API call
    let userType = sessionStorage.getItem('userType'); // e.g., 'passenger' or 'car_owner'
    console.log('User type:', userType); // Debugging line

    let menu = document.getElementById('user-menu');
    menu.innerHTML = '<ul class="menu-list">'; // Clear the menu and start a list

    // Check if userType is null and handle the case
    if (!userType) {
        console.log('No user type found in sessionStorage. Redirecting to login.');
        window.location.href = 'login.html'; // Redirect to login page if userType is not set
        return;
    }
    
    // Based on user type, add relevant menu items
    if (userType === 'car_owner') {
        menu.innerHTML += `
            <li><button onclick="publishTrip()">Publish a Trip</button></li>
            <li><button onclick="manageTrips()">Manage Trips</button></li>
        `;
        // Add other car owner specific menu items here
    } else if (userType === 'passenger') {
        menu.innerHTML += `
            <li><button onclick="browseTrips()">Browse Trips</button></li>
            <li><button onclick="viewEnrolledTrips()">View Enrolled Trips</button></li>
        `;
        // Add other passenger specific menu items here
    }

    menu.innerHTML += '</ul>'; // Close the list
}


// CAR OWNER MENU
function publishTrip() {
    const mainContent = document.getElementById('main-content');
    
    mainContent.innerHTML = `
        <h2>Publish a Trip</h2>
        <form id="publishTripForm">
            <input type="text" id="pick_up_location" placeholder="Pick-up Location" required>
            <input type="text" id="alt_pick_up_location" placeholder="Alternative Pick-up Location">
            <input type="text" id="destination_address" placeholder="Destination Address" required>
            <input type="datetime-local" id="travel_start_time" placeholder="Start Time" required>
            <input type="number" id="number_of_seats" placeholder="Number of Seats" required min="1">
            <button type="submit">Publish</button>
        </form>
    `;

    document.getElementById('publishTripForm').addEventListener('submit', publishTripHandler);
}

function publishTripHandler(event) {
    event.preventDefault(); // Prevent the default form submission
    //console.log('Form submission triggered');

    // Retrieve user data from sessionStorage
    const user = JSON.parse(sessionStorage.getItem('user'));
    const carOwnerId = user.userID; 
    const carOwnerName = user.first_name + ' ' + user.last_name;

    // Collect trip data from form inputs
    var pickUpLocation = document.getElementById('pick_up_location').value;
    var alternativePickUp = document.getElementById('alt_pick_up_location').value || ""; // Optional field
    var destinationAddress = document.getElementById('destination_address').value;
    var travelStartTime = document.getElementById('travel_start_time').value; // Needs to be in ISO format
    var availableSeats = document.getElementById('number_of_seats').value;

    // Convert to Date object and check if it's at least 30 minutes ahead
    var startTime = new Date(travelStartTime);
    var currentTime = new Date();
    var thirtyMinutesLater = new Date(currentTime.getTime() + 30 * 60000);

    if (startTime < thirtyMinutesLater) {
        alert('Trip start time must be at least 30 minutes in the future.');
        return; // Stop execution if the time is not valid
    }
    
    // Create trip data object
    var tripData = {
        car_owner_id: carOwnerId,
        car_owner_name: carOwnerName,
        pick_up_location: pickUpLocation,
        alternative_pick_up: alternativePickUp,
        destination_address: destinationAddress,
        travel_start_time: new Date(travelStartTime).toISOString(), // Convert to ISO format
        available_seats: parseInt(availableSeats, 10) // Parse to int
    };
    //console.log('Sending trip data:', tripData);
    
    // Send POST request to the server
    fetch('http://localhost:5001/trips', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify(tripData),
    })
    .then(response => {
        if (!response.ok) {
            //throw new Error(`HTTP error! status: ${response.status}`);
        }
        return response.json();
    })
    .then(data => {
        // Handle success
        console.log('Trip published successfully:', data);
        alert('Trip published successfully!');
    })
    .catch(error => {
        // Handle errors
        console.error('Error publishing trip:', error);
        alert(`Error publishing trip: ${error.message}`);
    });

    document.getElementById('publishTripForm').addEventListener('submit', publishTripHandler);
}



function manageTrips() {
    // Code to handle managing trips
}





// PASSENGER MENU
function browseTrips() {
    // Code to handle browsing trips
    // Make an AJAX call to get the available trips
    fetch('http://localhost:5001/available-trips')
        .then(response => response.json())
        .then(trips => {
            // Assuming main-content is the id of the main content area where you want to load the trips
            const mainContent = document.getElementById('main-content');
            mainContent.innerHTML = '<h2>Available Trips:</h2>';
            
            // Create a list of trips
            trips.forEach(trip => {
                mainContent.innerHTML += `
                    <div class="trip">
                        <p>Pick Up Location: ${trip.pick_up_location}</p>
                        <p>Destination Address: ${trip.destination_address}</p>
                        <p>Start Traveling Time: ${new Date(trip.travel_start_time).toLocaleString()}</p>
                        <button onclick="enroll(${trip.ID})">Enroll</button>
                    </div>
                `;
            });
        })
        .catch(error => {
            console.error('Error fetching available trips:', error);
        });
}

function enroll(tripId) {
    // Make an AJAX call to enroll in a trip
    const passengerId = sessionStorage.getItem('userId'); 

    fetch('http://localhost:5001/enroll', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/json',
        },
        body: JSON.stringify({ passenger_id: passengerId, trip_id: tripId }),
    })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            alert('Successfully enrolled in the trip!');
            // Update the UI accordingly
        })
        .catch(error => {
            console.error('Error enrolling in trip:', error);
        });
}

function viewEnrolledTrips() {
    // Code to handle viewing enrolled trips
    const passengerId = sessionStorage.getItem('userId');
    fetch(`http://localhost:5001/enrolled-trips?passenger_id=${passengerId}`) // Assuming this is the correct endpoint
    .then(response => response.json())
    .then(trips => {
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = '<h2>Your Enrolled Trips:</h2>';
        
        // Create a list of enrolled trips
        trips.forEach(trip => {
            mainContent.innerHTML += `
                <div class="trip">
                    <p>Trip ID: ${trip.id}</p>
                    <p>Pick Up Location: ${trip.pick_up_location}</p>
                    <p>Destination Address: ${trip.destination_address}</p>
                    <p>Travel Start Time: ${new Date(trip.travel_start_time).toLocaleString()}</p>
                    <button onclick="cancelEnrollment(${trip.id})">Cancel Enrollment</button>
                </div>
            `;
        });
    })
    .catch(error => {
        console.error('Error fetching enrolled trips:', error);
    });
}




function logout() {
    // Code to handle logout
    sessionStorage.clear(); // Clear session storage
    window.location.href = 'login.html'; // Redirect to login page
}

