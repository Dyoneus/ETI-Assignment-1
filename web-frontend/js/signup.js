//  /web-frontend/js/signup.js

document.getElementById('signupForm').addEventListener('submit', function(event) {
    event.preventDefault();

    var firstName = document.getElementById('firstName').value;
    var lastName = document.getElementById('lastName').value;
    var mobile = document.getElementById('mobile').value;
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;

    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:5000/users', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
        if (xhr.status === 201) {
            alert('Signup successful!');
            window.location.href = 'login.html'; // Redirect to login page after successful signup
        } else {
            alert('Signup failed: ' + xhr.responseText);
        }
    };
    xhr.send(JSON.stringify({
        first_name: firstName,
        last_name: lastName,
        mobile: mobile,
        email: email,
        password: password
    }));
});