const form = document.querySelector("[data-reg-dialog]");
const secondaryBtn = document.querySelector("[data-secondary-btn]");

secondaryBtn.addEventListener('click', () => {
    console.log("Кнопка 'Войти' нажата"); 
    window.location.href = '../login/index.html';
});

form.addEventListener('submit', async(e) => {
    e.preventDefault();

    const login = document.querySelector("[data-login]").value;
    const password = document.querySelector("[data-password]").value;
    const secondPassword = document.querySelector("[data-second-password]").value;

    if (password !== secondPassword) {
        alert("Ошибка: пароли не совпадают.");
        return; 
    }

    const data = {
        login: login,
        password: password
    };

    
    try {
        const response = await fetch('http://localhost:8080/auth/register', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        // user exists - уже зарегался, user - не зарегался еще
        const answer = await response.json();

        if (response.ok) {
            console.log("успешная проверка");
            if (answer.role === 'user') {
                console.log("регистрация прошла успешно, пользователя еще нет в системе.");
                localStorage.setItem('username', login);
                window.location.href = '../user/index.html';
            } else if (answer.role === 'user exists') {
                console.log("пользователь уже существует.");
                alert("Пользователь с таким логином уже существует.");
            }
        } else {
            console.log("ошибка регистрации.");
            alert("Ошибка регистрации: неизвестная ошибка.");
        } 
    }
    catch (error) {
        console.error("Ошибка:", error);
    }
});