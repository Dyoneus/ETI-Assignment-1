//  /web-frontend/js/auth.js

// User login
document.getElementById('loginForm').addEventListener('submit', function(event) {
    event.preventDefault();
    var email = document.getElementById('email').value;
    var password = document.getElementById('password').value;
    var xhr = new XMLHttpRequest();
    xhr.open('POST', 'http://localhost:5000/login', true);
    xhr.setRequestHeader('Content-Type', 'application/json');
    xhr.onload = function() {
        if (xhr.status === 200) {
            var userInfo = JSON.parse(xhr.responseText);
            sessionStorage.setItem('user', JSON.stringify(userInfo));
            sessionStorage.setItem('userType', userInfo.userType); // Assuming userType is part of the userInfo object
            window.location.href = 'mainMenu.html'; // Redirect to main menu page
        } else {
            alert('Login failed: ' + xhr.responseText);
        }
    };
    xhr.send(JSON.stringify({email: email, password: password}));
});