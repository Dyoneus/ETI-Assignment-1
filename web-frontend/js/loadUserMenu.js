//  /web-frontend/js/loadUserMenu.js

function loadUserMenu() {
    // Fetch the user type from session storage or an API call
    let userType = sessionStorage.getItem('userType'); // e.g., 'passenger' or 'car_owner'
    console.log('User type:', userType); // Debugging line
    
    let menu = document.getElementById('user-menu');

    // Clear the menu
    menu.innerHTML = ''; 

    // Check if userType is null and handle the case
    if (!userType) {
        console.log('No user type found in sessionStorage. Redirecting to login.');
        window.location.href = 'login.html'; // Redirect to login page if userType is not set
        return;
    }

    // Based on user type, add relevant menu items
    if (userType === 'car_owner') {
        // Add car owner specific menu items
        menu.innerHTML += '<button onclick="publishTrip()">Publish a Trip</button>';
        // ... other menu items
    } else if (userType === 'passenger') {
        // Add passenger specific menu items
        menu.innerHTML += '<button onclick="browseTrips()">Browse Trips</button>';
        // ... other menu items
    }
}
// MENU OPTION FOR BOTH USER TYPES

// CAR OWNER MENU
function publishTrip() {
    // Code to handle publishing a trip
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
}


function logout() {
    // Code to handle logout
    sessionStorage.clear(); // Clear session storage
    window.location.href = 'login.html'; // Redirect to login page
}

