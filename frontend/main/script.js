document.addEventListener('DOMContentLoaded', function () {
    const form = document.getElementById('expressionForm');
    const logoutButton = document.getElementById('logoutButton');
    const userNameSpan = document.getElementById('userName');
    
    const token = document.cookie.split('; ').find(row => row.startsWith('auth_token=')).split('=')[1];
    if (token) {
        const decodedToken = jwt_decode(token);
        
        if (!localStorage.getItem('username')) {
            fetch('/auth/username/', {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json',
                    'Authorization': `Bearer ${token}`
                }
            })
            .then(response => response.json())
            .then(data => {
                localStorage.setItem('username', data.username);
                userNameSpan.textContent = data.username;
            })
            .catch(error => console.error('Ошибка при получении имени пользователя:', error));
        } else {
            userNameSpan.textContent = localStorage.getItem('username');
        }
    }

    form.addEventListener('submit', function (event) {
        event.preventDefault();
        const expression = form.elements.expression.value;
        fetch('http://localhost:8080/expression/', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json',
            },
            body: JSON.stringify({ expression: expression }),
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка сети');
                }
                return response.json();
            })
            .then(data => {
                console.log('Успех:', data);
                form.elements.expression.value = '';
                loadExpressions();
            })
            .catch((error) => {
                console.error('Ошибка:', error);
                form.elements.expression.value = '';
            });
    });

    function loadExpressions() {
        fetch('http://localhost:8080/expression/', {
            method: 'GET',
            headers: {
                'Content-Type': 'application/json',
            },
        })
            .then(response => {
                if (!response.ok) {
                    throw new Error('Ошибка сети');
                }
                return response.json();
            })
            .then(data => {
                console.log('Получено:', data);
                const expressionsDiv = document.getElementById('savedExpressions');
                expressionsDiv.innerHTML = '';
                data.forEach(expression => {
                    const p = document.createElement('p');
                    p.textContent = `ID: ${expression.id}, Выражение: ${expression.expression}, Ответ: ${expression.answer}, Статус: ${expression.status}`;
                    expressionsDiv.appendChild(p);
                });
            })
            .catch(error => {
                console.error('Ошибка при получении данных:', error);
            });
    }

    setInterval(loadExpressions, 5000);

    function jwt_decode(token) {
        const base64Url = token.split('.')[1];
        const base64 = base64Url.replace(/-/g, '+').replace(/_/g, '/');
        const jsonPayload = decodeURIComponent(atob(base64).split('').map(function (c) {
            return '%' + ('00' + c.charCodeAt(0).toString(16)).slice(-2);
        }).join(''));
        return JSON.parse(jsonPayload);
    }

    logoutButton.addEventListener('click', function () {
        document.cookie = "auth_token=; expires=Thu, 01 Jan 1970 00:00:00 UTC; path=/;";
        localStorage.removeItem('username');
        window.location.href = "http://localhost:8080/auth";
    });

    loadExpressions();
});