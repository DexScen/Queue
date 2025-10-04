document.addEventListener('DOMContentLoaded', async function () {
    const username = localStorage.getItem('username');

    if (username) {
        const greeting = document.querySelector('.text-title');
        const profileLink = document.querySelector('.nav-a[href="/profile/index.html"]');

        if (greeting) {
            greeting.textContent = `${username}, вот все ваши очереди:`;
        }

        if (profileLink) {
            profileLink.textContent = username;
        }

        await loadAndRenderQueues(username);
        startAutoRefresh(username);  
    }
});

async function loadAndRenderQueues(username) {
    try {
            const response = await fetch(`http://localhost:8080/queue/${username}`);

            if (response.ok) {
                const userData = await response.json();
                renderQueues(userData);
            }
        } catch (error) {
            console.error('Ошибка при загрузке очередей:', error);
        }
}

function startAutoRefresh(username) {
    setInterval(async () => {
        console.log("автоматическое обновление данных... (кд - 5 сек)");
        await loadAndRenderQueues(username);
    }, 5000); // 5000 мс 
}

function renderQueues(queues) {
    const container = document.querySelector('[data-container]');
    container.innerHTML = '';

    if (!queues || !Array.isArray(queues)) {
        console.log('Нет данных об очередях или данные некорректны');
        container.innerHTML = '<p>У вас нет активных очередей</p>';
        return;
    }

    if (queues.length === 0) {
        console.log('Массив очередей пуст');
        container.innerHTML = '<p>У вас нет активных очередей</p>';
        return;
    }

    queues.forEach(stand => {
        const card = document.createElement('div');
        card.className = 'card-profile back-color-white';
        card.setAttribute('data-stand-id', stand.id);

        const waitTimeMinutes = Math.ceil(stand.current_people * stand.duration_seconds / 60);

        //  <button type="button" class="btn-standart">обновить</button>
        card.innerHTML = `
            <h3 class="queue-name">${stand.name}</h3>
            <p class="queue-info">
                еще ${waitTimeMinutes} мин<br>
                очередь: ${stand.current_people} чел.
            </p>
            <div class="queue-actions">
                <button type="button" class="btn-delete" data-stand-id="${stand.id}">удалить</button>
            </div>
        `;

        container.appendChild(card);
    });

    document.querySelectorAll('.btn-delete').forEach(button => {
        button.addEventListener('click', async function () {
            const standId = this.getAttribute('data-stand-id');
            const username = localStorage.getItem('username');
            const response = await fetch(`http://localhost:8080/auth/${username}`, {
                method: 'GET',
                headers: {
                    'Content-Type': 'application/json'
                }
            });

            const data = await response.json();

            const userID = data.id;
            console.log(userID);

            try {
                const response = await fetch('http://localhost:8080/remove', {
                    method: 'DELETE',
                    headers: {
                        'Content-Type': 'application/json'
                    },
                    body: JSON.stringify({
                        user_id: userID,
                        game_id: parseInt(standId)
                    })
                });

                if (response.ok) {
                    const cardToRemove = document.querySelector(`[data-stand-id="${standId}"]`);
                    if (cardToRemove) {
                        cardToRemove.remove();
                    }
                    await loadAndRenderQueues(username);

                    console.log(`Очередь ${standId} успешно удалена`);
                } else {
                    console.error('Ошибка при удалении:', response.status);
                    alert('Не удалось удалить запись из очереди');
                }
            } catch (error) {
                console.error('Ошибка:', error);
                alert('Произошла ошибка при удалении');
            }
        });
    });
}