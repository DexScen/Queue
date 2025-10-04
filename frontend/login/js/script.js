const form = document.querySelector("[data-login-dialog]");
const secondaryBtn = document.querySelector("[data-secondary-btn]");

secondaryBtn.addEventListener('click', () => {
    console.log("Кнопка 'Зарегистрироваться' нажата");
    window.location.href = '/registration/';
});

form.addEventListener('submit', async (e) => {
    e.preventDefault();

    const login = document.querySelector("[data-login]").value;
    const password = document.querySelector("[data-password]").value;

    const data = {
        login: login,
        password: password
    };

    try {
        const response = await fetch('http://localhost:8080/auth/login', {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify(data)
        });

        const answer = await response.json();

        if (answer.role === 'user') {
            console.log("успешный вход");
            localStorage.setItem('username', login);
            window.location.href = '/user/';
        } else {
            if (answer.role === 'user not found') {
                alert("Пользователь не найден");
            } else if (answer.role === 'wrong password') {
                alert("Неверный пароль");
            } else {
                alert("Ошибка входа: неизвестная ошибка");
            }
        }
    } catch (error) {
        console.error("Ошибка:", error);
        alert("Ошибка соединения с сервером");
    }
});