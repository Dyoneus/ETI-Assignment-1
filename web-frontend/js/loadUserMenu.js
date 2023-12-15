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
            manageTrips();
        })
        .catch(error => {
            // Handle errors
            console.error('Error publishing trip:', error);
            alert(`Error publishing trip: ${error.message}`);
        });

        document.getElementById('publishTripForm').addEventListener('submit', publishTripHandler);
    }

    function manageTrips() {
        // Fetch trips from the server and display them
        const user = JSON.parse(sessionStorage.getItem('user'));
        const userID = user.userID;

        fetch(`http://localhost:5001/trips?carOwnerID=${userID}`)
            .then(response => {
                if (!response.ok) {
                    //throw new Error(`HTTP error! status: ${response.status}`);
                }
                return response.json();
            })
            .then(trips => {
                //console.log('Raw trips data:', trips); // Log the raw data
                const mainContent = document.getElementById('main-content');
                mainContent.innerHTML = '<h2>Your Trips:</h2>';
            
                trips.forEach(trip => {
                    //console.log('Trip ID:', trip.ID); // Log the trip ID here for debugging

                    // For each trip, create HTML to display it
                    mainContent.innerHTML += `
                        <div class="trip">
                            <p>Pick Up Location: ${trip.pick_up_location}</p>
                            <p>Destination Address: ${trip.destination_address}</p>
                            <p>Travel Start Time: ${new Date(trip.travel_start_time).toLocaleString()}</p>
                            <button onclick="editTrip(${trip.ID})">Edit</button>
                            <button onclick="deleteTrip(${trip.ID})">Delete</button>
                        </div>
                    `;
                });
            })
            .catch(error => {
                console.error('Error fetching trips:', error);
            });
    }

    // This function is called when the "Edit" button is clicked.
    function editTrip(tripId) {
        if (!tripId) {
            console.error('Trip ID is undefined or not a number.');
            return;
        }

        // Fetch the trip data from the server.
        fetch(`http://localhost:5001/trips/${tripId}`)
        .then(response => response.json())
        .then(trip => {
            // Populate the form with the trip data.
            showEditForm(trip);
        })
        .catch(error => {
            console.error('Error fetching trip details:', error);
            alert(`Error fetching trip details: ${error.message}`);
        });
    }

    // This function displays the edit form populated with the trip data.
    function showEditForm(trip) {
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = `
            <h2>Edit Trip</h2>
            <form id="editTripForm">
                <input type="hidden" id="edit_trip_id" value="${trip.ID}" readonly>
                <label for="edit_pick_up_location">Pick-up Location:</label>
                <input type="text" id="edit_pick_up_location" value="${trip.pick_up_location}" required>
                <label for="edit_alt_pick_up_location">Alternative Pick-up Location:</label>
                <input type="text" id="edit_alt_pick_up_location" value="${trip.alternative_pick_up}">
                <label for="edit_destination_address">Destination Address:</label>
                <input type="text" id="edit_destination_address" value="${trip.destination_address}" required>
                <label for="edit_travel_start_time">Travel Start Time:</label>
                <input type="datetime-local" id="edit_travel_start_time" value="${formatDateToInput(trip.travel_start_time)}" required>
                <label for="edit_number_of_seats">Available Seats:</label>
                <input type="text" id="edit_number_of_seats" value="${trip.available_seats}" readonly>
                <button type="submit">Save Changes</button>
            </form>
        `;

        // Add submit event listener to the form.
        document.getElementById('editTripForm').addEventListener('submit', function(event) {
            event.preventDefault();
            submitUpdatedTrip(trip.ID); // Pass the trip.ID for updating.
        });
    }

    function submitUpdatedTrip(tripId) {
        //console.log('Trip ID:', tripId); // Log the trip ID here for debugging
        if (!tripId) {
            console.error('Trip ID is undefined or not a number.');
            return;
        }
        const user = JSON.parse(sessionStorage.getItem('user'));
        const carOwnerId = user.userID; 
        const carOwnerName = user.first_name + ' ' + user.last_name;
    
        
        // Collect trip data from form inputs
        var pickUpLocation = document.getElementById('edit_pick_up_location').value;
        var alternativePickUp = document.getElementById('edit_alt_pick_up_location').value || ""; // Optional field
        var destinationAddress = document.getElementById('edit_destination_address').value;
        var travelStartTime = document.getElementById('edit_travel_start_time').value; // Needs to be in ISO format
        var availableSeats = document.getElementById('edit_number_of_seats').value;

        // Convert to Date object and check if it's at least 30 minutes ahead
        var startTime = new Date(travelStartTime);
        var currentTime = new Date();
        var thirtyMinutesLater = new Date(currentTime.getTime() + 30 * 60000);

        if (startTime < thirtyMinutesLater) {
            alert('Trip start time must be at least 30 minutes in the future.');
            return; // Stop execution if the time is not valid
        }

        // Collect data from form fields
        const updatedTripData = {
            id: tripId,
            car_owner_id: carOwnerId,
            car_owner_name: carOwnerName,
            pick_up_location: pickUpLocation,
            alternative_pick_up: alternativePickUp,
            destination_address: destinationAddress,
            travel_start_time: new Date(travelStartTime).toISOString(), // Convert to ISO format
            available_seats: parseInt(availableSeats, 10) // Parse to int
        };
        
        // Check if the date is valid
        if (!updatedTripData.travel_start_time) {
            console.error('Travel start time is invalid.');
            alert('Please enter a valid travel start time.');
            return;
        }
    
        // Make a PATCH request to the server
        fetch(`http://localhost:5001/trips/${tripId}`, {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify(updatedTripData),
        })
        .then(response => {
            if (!response.ok) {
                // If the HTTP status code is not successful, throw an error
                return response.text().then(text => Promise.reject(new Error(text)));
            }
            return response.text(); // Get the response as text
        })
        .then(text => {
            // The server responds with a text message. No need to parse as JSON.
            console.log('Trip updated successfully:', text);
            alert('Trip updated successfully!');
            manageTrips(); // Reload the list of trips
        })
        .catch(error => {
            console.error('Error updating trip:', error);
            alert(`Error updating trip: ${error}`);
        });
    }

    function formatDateToInput(dateTime) {
        // Create a new date object using the provided dateTime
        const date = new Date(dateTime);
    
        // Convert it to an ISO string and then slice it to match the datetime-local input format
        return date.toISOString().slice(0, 16);
    }

    function deleteTrip(tripId) {
        if (!tripId) {
            console.error('Trip ID is undefined or not a number.');
            return;
        }

        fetch(`http://localhost:5001/trips/${tripId}`, {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            }
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(() => {
            alert('Trip deleted successfully!');
            manageTrips(); // Reload the list of trips
        })
        .catch(error => {
            //console.error('Error deleting trip:', error);
            alert('Trip deleted successfully!');
            manageTrips(); // Reload the list of trips
        });
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
                            <p>Driver: ${trip.car_owner_name}</p>
                            <p>Pick-up Location: ${trip.pick_up_location}</p>
                            <p>Alt. Pick-up Location: ${trip.alternative_pick_up}</p>
                            <p>Destination Address: ${trip.destination_address}</p>
                            <p>Start Traveling Time: ${new Date(trip.travel_start_time).toLocaleString()}</p>
                            <p>Available Seat(s): ${trip.available_seats}</p>
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
        const user = JSON.parse(sessionStorage.getItem('user'));
        const passengerId = user.userID;
    
        fetch('http://localhost:5001/enroll', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ trip_id: tripId, passenger_id: passengerId }),
        })
        .then(response => {
            if (!response.ok) {
                // If the HTTP status code is not successful, get the response as text and throw an error
                return response.text().then(text => Promise.reject(text));
            }
            // If the response is ok, parse it as JSON
            return response.json();
        })
        .then(data => {
            alert('Successfully enrolled in the trip!');
            // Update the UI accordingly
            browseTrips();
        })
        .catch(error => {
            // Check if the error is a string (which means it's the text response from the server)
            if (typeof error === 'string') {
                if (error.startsWith("User is al")) {
                    alert('You have already enrolled in this trip!');
                } else {
                    alert('Successfully enrolled in the trip!');
                    viewEnrolledTrips();
                }
            } else if (error instanceof Error) {
                alert('Successfully enrolled in the trip!');
                viewEnrolledTrips();
            } else {
                // Handle other cases or unknown errors
                console.error('Unknown error enrolling in trip:', error);
                alert('An unknown error occurred while enrolling in the trip.');
            }
            viewEnrolledTrips();
        });
    }

    function viewEnrolledTrips() {
        const user = JSON.parse(sessionStorage.getItem('user'));
        const passengerId = user ? user.userID : null;
        if (!passengerId) {
            console.error('No user ID found in session storage');
            return;
        }
    
        fetch(`http://localhost:5001/enrolled-trips?passenger_id=${passengerId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! Status: ${response.status}`);
            }
            return response.json();
        })
        .then(trips => {
            //console.log('Enrolled trips:', trips); // Log the trips to console for debugging
            const mainContent = document.getElementById('main-content');
            mainContent.innerHTML = '<h2>Your Enrolled Trips:</h2>';
    
            if (trips.length === 0) {
                mainContent.innerHTML += '<p>You have no enrolled trips.</p>';
            } else {
                // Create a list of enrolled trips
                trips.forEach(trip => {
                    mainContent.innerHTML += `
                        <div class="trip">
                            <p>Driver: ${trip.car_owner_name}</p>
                            <p>Pick Up Location: ${trip.pick_up_location}</p>
                            <p>Alt. Pick-up Location: ${trip.alternative_pick_up}</p>
                            <p>Destination Address: ${trip.destination_address}</p>
                            <p>Travel Start Time: ${new Date(trip.travel_start_time).toLocaleString()}</p>
                            <button onclick="cancelEnrollment(${trip.ID})">Cancel Enrollment</button>
                        </div>
                    `;
                });
            }
        })
        .catch(error => {
            console.error('Error fetching enrolled trips:', error);
        });
    }   

    function cancelEnrollment(tripId) {
        const user = JSON.parse(sessionStorage.getItem('user'));
        const passengerId = user.userID;
        
        fetch('http://localhost:5001/cancel-enrollment', {
            method: 'DELETE',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ trip_id: tripId, passenger_id: passengerId }),
        })
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(data => {
            alert('Successfully canceled enrollment in the trip. You will not be able to enroll in this trip again.');
            // Update the UI here
            // For example, refresh the list of enrolled trips or update the seat count
            viewEnrolledTrips(); // This would re-fetch the enrolled trips
        })
        .catch(error => {
            //console.error('Error canceling enrollment:', error);
            //alert(`Error canceling enrollment: ${error.message}`);
            alert('Successfully canceled enrollment in the trip. You will not be able to enroll in this trip again.');
            viewEnrolledTrips(); // This would re-fetch the enrolled trips
        });
    }



    function logout() {
        // Code to handle logout
        sessionStorage.clear(); // Clear session storage
        window.location.href = 'login.html'; // Redirect to login page
    }

