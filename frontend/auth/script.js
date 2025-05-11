document.addEventListener('DOMContentLoaded', function() {
    const authForm = document.getElementById('authForm');
    authForm.addEventListener('submit', function(event) {
        if (event.submitter.id === 'signupButton') {
            event.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            fetch('http://localhost:8080/auth/signup/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                }),
            })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Network response was not ok');
                }
                return response.json();
            })
            .then(data => {
                console.log('Signup successful:', data);
                authForm.elements.username.value = '';
                authForm.elements.password.value = '';
            })
            .catch(error => console.error('Error:', error));
        } else if (event.submitter.id === 'loginButton') {
            event.preventDefault();
            const username = document.getElementById('username').value;
            const password = document.getElementById('password').value;
            fetch('http://localhost:8080/auth/login/', {
                method: 'POST',
                headers: {
                    'Content-Type': 'application/json',
                },
                body: JSON.stringify({
                    username: username,
                    password: password
                }),
            })
            .then(response => {
                if (!response.ok) {
                    console.error('Response status:', response.status)
                    return response.text().then(text => {
                        throw new Error('Network response was not ok: ' + text);
                    });
                }
                const contentType = response.headers.get("content-type");
                if (!contentType || !contentType.includes("application/json")) {
                    window.location.href = "http://localhost:8080/";
                    return;
                }
                return response.json();
            })
            .then(data => {
                console.log('Login successful, token stored.');
                authForm.elements.username.value = '';
                authForm.elements.password.value = '';
            })
            .catch(error => { 
                console.error('Error:', error);
                alert('Login failed: ' + error.message);
            });
        }
    });
});