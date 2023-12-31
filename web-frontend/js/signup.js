//  /web-frontend/js/signup.js
// Function to validate mobile number (only integers allowed)
function isValidMobile(mobile) {
    // Regular expression pattern for integers only
    var pattern = /^\d+$/;
    return pattern.test(mobile);
}

document.getElementById('signupForm').addEventListener('submit', function(event) {
    event.preventDefault();

    var firstName = document.getElementById('firstName').value;
    var lastName = document.getElementById('lastName').value;
    var mobile = document.getElementById('mobile').value;
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;
    var confirmPassword = document.getElementById('confirmPassword').value;

    // Check if mobile number is valid
    if (!isValidMobile(mobile)) {
        alert("Mobile number must contain digits only.");
        return; // Stop execution if mobile number is not valid
    }

    // Check if both passwords match
    if (password === confirmPassword) {
        console.log('Passwords match. Proceed with sign up.');
        var xhr = new XMLHttpRequest();
        xhr.open('POST', 'http://localhost:5000/users', true);
        xhr.setRequestHeader('Content-Type', 'application/json');
        xhr.onload = function() {
            if (xhr.status === 201) {
                alert('Sign Up successful!');
                window.location.href = 'index.html'; // Redirect to login page after successful signup
            } else {
                alert('Sign Up failed: ' + xhr.responseText);
            }
        };
        xhr.send(JSON.stringify({
            first_name: firstName,
            last_name: lastName,
            mobile: mobile,
            email: email,
            password: password
        }));
    } else {
        // If passwords do not match, prevent form submission and alert the user
        alert("Passwords do not match. Please try again.");
    }
});