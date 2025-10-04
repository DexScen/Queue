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

        try {
            const response = await fetch(`http://localhost:8080/queue/${username}`);

            if (response.ok) {
                const userData = await response.json();
                renderQueues(userData);
            }
        } catch (error) {
            console.error('Ошибка:', error);
        }
    }
});

function renderQueues(queues) {
    const container = document.querySelector('[data-container]');
    container.innerHTML = '';

    queues.forEach(stand => {
        const card = document.createElement('div');
        card.className = 'card-profile back-color-white';
        card.setAttribute('data-stand-id', stand.id);

        const waitTimeMinutes = Math.ceil(stand.current_people * stand.duration_seconds / 60);

        card.innerHTML = `
            <h3 class="queue-name">${stand.name}</h3>
            <p class="queue-info">
                еще ${waitTimeMinutes} мин<br>
                очередь: ${stand.current_people} чел.
            </p>
            <div class="queue-actions">
                <button type="button" class="btn-standart">обновить</button>
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

            // Сначала преобразуем ответ в JSON
            const data = await response.json();

            // Теперь можно получить id
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