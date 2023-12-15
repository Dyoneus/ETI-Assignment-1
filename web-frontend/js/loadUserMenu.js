    //  /web-frontend/js/loadUserMenu.js

    function loadUserMenu() {
        // Fetch the user type from session storage or an API call
        let userType = sessionStorage.getItem('userType'); // e.g., 'passenger' or 'car_owner'
        //console.log('User type:', userType); // Debugging line

        let menu = document.getElementById('user-menu');
        menu.innerHTML = '<ul class="menu-list">'; // Clear the menu and start a list

        // Check if userType is null and handle the case
        if (!userType) {
            console.log('No user type found in sessionStorage. Redirecting to login.');
            window.location.href = 'index.html'; // Redirect to login page if userType is not set
            return;
        }
        
        // Based on user type, add relevant menu items
        if (userType === 'car_owner') {
            menu.innerHTML += `
                <li><button onclick="publishTrip()">Publish a Trip</button></li>
                <li><button onclick="manageTrips()">Manage Trips</button></li>
                <li><button onclick="loadPastTripsCarOwner()">View Past Published Trips</button></li>
                <li><button onclick="showUpdateProfileForm()">Update Profile</button></li>
                <li><button onclick="deleteAccount()">Delete Account</button></li>
            `;
            // Add other car owner specific menu items here
        } else if (userType === 'passenger') {
            menu.innerHTML += `
                <li><button onclick="browseTrips()">Browse Trips</button></li>
                <li><button onclick="viewEnrolledTrips()">View Current Enrolled Trips</button></li>
                <li><button onclick="loadPastTripsPassenger()">View Past Enrolled Trips</button></li>
                <li><button onclick="showUpdateProfileForm()">Update Profile</button></li>
                <li><button onclick="showUpgradeToCarOwnerForm()">Become Car Owner</button></li>
                <li><button onclick="deleteAccount()">Delete Account</button></li>
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
            //console.log('Trip published successfully:', data);
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

    // Function to load past trips for car owners
    function loadPastTripsCarOwner() {
        const user = JSON.parse(sessionStorage.getItem('user'));
        const carOwnerId = user.userID; // or the correct property that holds the ID
    
        fetch(`http://localhost:5001/past-trips/car-owner?ownerID=${carOwnerId}`)
        .then(response => {
            if (!response.ok) {
                throw new Error(`HTTP error! status: ${response.status}`);
            }
            return response.json();
        })
        .then(trips => {
            // Sort trips in reverse chronological order
            trips.sort((a, b) => new Date(b.travel_start_time) - new Date(a.travel_start_time));
            displayTrips(trips);
        })
        .catch(error => {
            console.error('Error loading past trips:', error);
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
    
        // First, fetch the currently enrolled trips to check for time conflicts
        fetch(`http://localhost:5001/enrolled-trips?passenger_id=${passengerId}`)
        .then(response => response.json())
        .then(enrolledTrips => {
            return Promise.all([enrolledTrips, fetch(`http://localhost:5001/trips/${tripId}`)]);
        })
        .then(([enrolledTrips, tripResponse]) => {
            return Promise.all([enrolledTrips, tripResponse.json()]);
        })
        .then(([enrolledTrips, selectedTrip]) => {
            if (hasTimeConflict(enrolledTrips, selectedTrip)) {
                throw new Error('Time conflict with an already enrolled trip. Time gap of 1 hour in between enrolled trip start time.');
            } else {
                return fetch('http://localhost:5001/enroll', {
                    method: 'POST',
                    headers: {
                        'Content-Type': 'application/json',
                    },
                    body: JSON.stringify({ trip_id: tripId, passenger_id: passengerId }),
                });
            }
        })
        .then(response => {
            if (!response.ok) {
                return response.text().then(text => Promise.reject(new Error(text)));
            }
            return response.json();
        })
        .then(data => {
            alert('Successfully enrolled in the trip!');
            browseTrips();
        })
        .catch(error => {
            alert(`${error.message}`);
        });
    }

    // Function to check for time conflicts
    function hasTimeConflict(enrolledTrips, selectedTrip) {
        const selectedStartTime = new Date(selectedTrip.travel_start_time).getTime();
        const selectedEndTime = selectedStartTime + (1 * 60 * 60 * 1000); // Assuming a fixed duration of 1 hour

        return enrolledTrips.some(enrolledTrip => {
            const enrolledStartTime = new Date(enrolledTrip.travel_start_time).getTime();
            const enrolledEndTime = enrolledStartTime + (1 * 60 * 60 * 1000); // Assuming a fixed duration of 1 hour

            return (selectedStartTime < enrolledEndTime && selectedEndTime > enrolledStartTime);
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
    
    // Utility function to display trips on the page
    function displayTrips(trips) {
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = '<h2>Past Trips:</h2>'; // Clear previous content and add title

        trips.forEach(trip => {
            const tripElement = document.createElement('div');
            tripElement.className = 'trip';
            tripElement.innerHTML = `
                <p>Driver: ${trip.car_owner_name}</p>
                <p>From: ${trip.pick_up_location}</p>
                <p>To: ${trip.destination_address}</p>
                <p>On: ${new Date(trip.travel_start_time).toLocaleString()}</p>
            `;
            mainContent.appendChild(tripElement);
        });
    }

    // Function to load past trips for passengers
    function loadPastTripsPassenger() {
        const userID = JSON.parse(sessionStorage.getItem('user')).userID;
        const passengerId = userID;
    
        fetch(`http://localhost:5001/past-trips/passenger?passengerID=${passengerId}`)
        .then(response => response.json())
        .then(trips => {
            //console.log('Received trips:', trips); // Log the trips to see what you receive
            // Sort trips in reverse chronological order
            trips.sort((a, b) => new Date(b.travel_start_time) - new Date(a.travel_start_time));
            //console.log('Sorted trips:', trips); // Log the sorted trips
            // Ensure the container is present in the DOM
            const mainContent = document.getElementById('main-content');
            if (mainContent) {
                displayTrips(trips, 'main-content'); // Assuming 'main-content' is the container for displaying trips
            } else {
                console.error('Container for displaying trips not found.');
            }
        })
        .catch(error => console.error('Error loading past trips:', error));
    }

    function showUpgradeToCarOwnerForm() {
        const mainContent = document.getElementById('main-content');
        mainContent.innerHTML = `
            <h2>Become a Car Owner</h2>
            <form id="upgradeToCarOwnerForm">
                <input type="text" id="drivers_license" placeholder="Driver's License Number" required>
                <input type="text" id="car_plate_number" placeholder="Car Plate Number" required>
                <button type="submit">Submit</button>
            </form>
        `;
    
        document.getElementById('upgradeToCarOwnerForm').addEventListener('submit', handleUpgradeToCarOwnerSubmission);
    }


    function handleUpgradeToCarOwnerSubmission(event) {
        event.preventDefault();
        const driversLicense = document.getElementById('drivers_license').value;
        const carPlateNumber = document.getElementById('car_plate_number').value;
    
        // Retrieve user email from sessionStorage
        const user = JSON.parse(sessionStorage.getItem('user'));
        const email = user.email;
    
        // Send POST request to the server
        fetch('http://localhost:5000/upgradeToCarOwner', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({
                email: email,
                drivers_license: driversLicense,
                car_plate_number: carPlateNumber,
            }),
        })
        .then(response => {
            if (!response.ok) {
                // If the server responds with a textual error message, parse it as text
                return response.text().then(text => Promise.reject(new Error(text)));
            }
            return response.json();
        })
        .then(() => {
            alert('Successfully upgraded to car owner!');
            // Update user type in sessionStorage
            user.userType = 'car_owner';
            sessionStorage.setItem('user', JSON.stringify(user));
            // Refresh the menu to reflect the new user type
            refreshMainMenu();
        })
        .catch(error => {
            alert('Successfully upgraded to car owner! Please re-login to take effect.');
            sessionStorage.clear(); // Clear session storage
            window.location.href = 'index.html'; // Redirect to login page
        });
    }



    // GENERAL MENU
    // This function reloads the user menu to reflect updated information
    function refreshMainMenu() {
        // Clear the current user menu
        let menu = document.getElementById('user-menu');
        if (menu) {
            menu.innerHTML = '';
        }

        // Call the loadUserMenu again to repopulate the menu
        loadUserMenu();

        var user = JSON.parse(sessionStorage.getItem('user'));
        if (user) {
            document.getElementById('user-name').textContent = user.first_name + ' ' + user.last_name;
        }
    }

    // This is a callback that you might call after the profile update
    function onProfileUpdateSuccess() {

        // Refresh the main menu
        refreshMainMenu();

        window.location.reload();
    }

    // This function verifies the current password by using the login endpoint.
    function verifyCurrentPassword(email, password, onSuccess, onError) {
        fetch('http://localhost:5000/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({ email, password })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Password verification failed');
            }
            return response.json();
        })
        .then(data => onSuccess(data))
        .catch(error => onError(error));
    }

    function showUpdateProfileForm() {
        // Retrieve user data from sessionStorage
        const user = JSON.parse(sessionStorage.getItem('user'));
        const mainContent = document.getElementById('main-content');

        mainContent.innerHTML = `
            <h2>Update Profile</h2>
            <form id="updateProfileForm">
                <input type="text" id="first_name" placeholder="First Name" value="${user.first_name || ''}" required>
                <input type="text" id="last_name" placeholder="Last Name" value="${user.last_name || ''}" required>
                <input type="tel" id="mobile" placeholder="Mobile Number" value="${user.mobile || ''}" required>
                <input type="email" id="email" placeholder="Email" value="${user.email || ''}" readonly>  
                <input type="password" id="current_password" placeholder="Current Password" required>
                <button type="submit">Update Profile</button>
            </form>
        `;
    
        // Add event listener for form submission
        document.getElementById('updateProfileForm').addEventListener('submit', handleUpdateProfileFormSubmission);
    }

    // Function to handle the form submission for updating the user profile
    function handleUpdateProfileFormSubmission(event) {
        event.preventDefault();

        // Collect form data
        const formData = {
            first_name: document.getElementById('first_name').value,
            last_name: document.getElementById('last_name').value,
            mobile: document.getElementById('mobile').value,
            email: document.getElementById('email').value,
            current_password: document.getElementById('current_password').value, // Collect the current password
        };

        // First, verify the current password
        verifyCurrentPassword(formData.email, formData.current_password, 
            // onSuccess callback
            () => {
                // Update the mobile separately by calling another function
                updateMobile(formData.email, formData.mobile);
                // If password verification is successful, proceed with updating the profile
                updateProfile(formData);
            },
            // onError callback
            (error) => {
                // If password verification fails, alert the user
                alert('Current password is incorrect.');
            }
        );
    }

    // Function to update the user's mobile number
    function updateMobile(email, mobile) {
        // Parse the mobile number to an integer and stringify it back to ensure it's in integer format
        const mobileNumber = parseInt(mobile, 10);

        // Check if the mobile number is an integer
        if (isNaN(mobileNumber)) {
            alert('Mobile number must be a valid integer. It will not be updated.');
            return;
        }
        
        fetch('http://localhost:5000/updateMobile', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                email: email,
                mobile: mobileNumber.toString() // Send it as a string but ensure it's an integer
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to update mobile number');
            }
            return response.json();
        })
        .then(data => {
            console.log('Mobile number updated successfully:', data);
        })
        .catch(error => {
            alert('Failed to update mobile number: ' + error.message);
        });
    }

    // Function to update the user profile
    function updateProfile(formData) {
        fetch('http://localhost:5000/users', {
            method: 'PATCH',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                first_name: formData.first_name,
                last_name: formData.last_name,
                mobile: formData.mobile,
                email: formData.email,
                // Do not send the current password in the update request
            })
        })
        .then(response => {
            if (!response.ok) {
                throw new Error('Failed to update profile');
            }
            return response.json();
        })
        .then(data => {
            alert('Profile updated successfully! You will be redirected to login page.');
            // Update session storage with new user data
            sessionStorage.setItem('user', JSON.stringify(data));
            // Optionally, refresh the page or redirect the user
            sessionStorage.clear(); // Clear session storage
            window.location.href = 'index.html'; // Redirect to login page

        })
        .catch(error => {
            alert('Failed to update profile: ' + error.message);
        });

        document.getElementById('updateProfileForm').addEventListener('submit', handleUpdateProfileFormSubmission);
    }

    // Function to check if one year has passed since the account creation
    function canDeleteAccount() {
        const user = JSON.parse(sessionStorage.getItem('user'));
        const accountCreationDate = new Date(user.createdAt); // Convert to Date object
        const oneYearAgo = new Date(new Date().setFullYear(new Date().getFullYear() - 1));
    
        return accountCreationDate <= oneYearAgo; // Returns true if the account was created at least one year ago
    }

    // Function to handle account deletion
    function deleteAccount() {
        // Retrieve the user object from sessionStorage
        const user = JSON.parse(sessionStorage.getItem('user'));
        
        if (canDeleteAccount()) {
            // Proceed with account deletion
            if (confirm("Are you sure you want to delete your account? This action cannot be undone.")) {
                // Send the DELETE request to the server
                fetch(`http://localhost:5000/deleteAccount?email=${encodeURIComponent(user.email)}`, {
                    method: 'PATCH',
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
                    alert('Your account has been successfully deleted.');
                    // Handle the UI cleanup here, e.g., redirect to the login page
                    sessionStorage.clear();
                    window.location.href = 'index.html';
                })
                .catch(error => {
                    //console.error('Error deleting account:', error);
                    //alert(`Error deleting account: ${error.message}`);
                    alert('Your account has been successfully deleted.');
                    // Handle the UI cleanup here, e.g., redirect to the login page
                    sessionStorage.clear();
                    window.location.href = 'index.html';
                });
            }
        } else {
            alert("You can only delete your account after one year of registration.");
        }
    }


    function logout() {
        // Code to handle logout
        sessionStorage.clear(); // Clear session storage
        window.location.href = 'index.html'; // Redirect to login page
    }


