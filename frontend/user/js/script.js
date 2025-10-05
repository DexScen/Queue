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

        await requestNotificationPermission();
        
        await loadAndRenderQueues(username);
        startAutoRefresh(username);  
    }
});

async function requestNotificationPermission() {
    if ('Notification' in window && Notification.permission === 'default') {
        try {
            const permission = await Notification.requestPermission();
            console.log('Разрешение на уведомления:', permission);
        } catch (error) {
            console.error('Ошибка при запросе разрешения:', error);
        }
    }
}

function sendMobileNotification(standName) {
    if (!('Notification' in window)) {
        console.log('Браузер не поддерживает уведомления');
        return;
    }

    if (Notification.permission === 'granted') {
        const notification = new Notification('Скоро ваша очередь!', {
            body: `На стенде "${standName}" перед вами остался 1 человек!`,
            // icon: '/icon.png', 
            // tag: 'queue-notification', 
            requireInteraction: true 
        });

        if ('vibrate' in navigator) {
            navigator.vibrate([200, 100, 200]); 
        }

        notification.onclick = function() {
            window.focus();
            notification.close();
        };
    }
}

function checkForNotification(queues) {
    if (!queues || !Array.isArray(queues)) return;

    queues.forEach(stand => {
        if (stand.current_people === 1) {
            const notificationKey = `notified_${stand.id}`;
            if (!localStorage.getItem(notificationKey)) {
                sendMobileNotification(stand.name);
                localStorage.setItem(notificationKey, 'true');
            }
        } else {
            const notificationKey = `notified_${stand.id}`;
            if (localStorage.getItem(notificationKey)) {
                localStorage.removeItem(notificationKey);
            }
        }
    });
}

async function loadAndRenderQueues(username) {
    try {
        const response = await fetch(`http://localhost:8080/queue/${username}`);

        if (response.ok) {
            const userData = await response.json();
            renderQueues(userData);
            checkForNotification(userData);
        }
    } catch (error) {
        console.error('Ошибка при загрузке очередей:', error);
    }
}

function startAutoRefresh(username) {
    setInterval(async () => {
        console.log("автоматическое обновление данных... (кд - 10 сек)");
        await loadAndRenderQueues(username);
    }, 5000); // 5000 мс = 5 секунд
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

        const waitTimeMinutes = Math.ceil((stand.position-1) * stand.duration_seconds / 60);

        card.innerHTML = `
            <h3 class="queue-name">${stand.name}</h3>
            <p class="queue-info">
                еще ${waitTimeMinutes} мин<br>
                позиция: ${stand.position-1}
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
                    cardToRemove.classList.add('smoke-vanish'); // добавляем анимацию
                    setTimeout(() => {
                        cardToRemove.remove(); // удаляем после анимации
                        loadAndRenderQueues(username); // перерисовываем
                    }, 800); // 0.8 сек — совпадает с CSS
                }

                // очищаем флаг уведомления
                localStorage.removeItem(`notified_${standId}`);

                console.log(`Очередь ${standId} успешно удалена`);
                }
                else {
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

function checkForNotification(queues) {
    if (!queues || !Array.isArray(queues)) return;

    // Сначала проверим все существующие уведомления в localStorage
    const allKeys = Object.keys(localStorage);
    const notificationKeys = allKeys.filter(key => key.startsWith('notified_'));
    
    // Очищаем уведомления для стендов, которых больше нет в очередях пользователя
    notificationKeys.forEach(key => {
        const standId = key.replace('notified_', '');
        const standExists = queues.some(stand => stand.id.toString() === standId);
        
        if (!standExists) {
            localStorage.removeItem(key);
            console.log(`Удалено уведомление для несуществующего стенда: ${standId}`);
        }
    });

    // Теперь обрабатываем текущие очереди
    queues.forEach(stand => {
        const notificationKey = `notified_${stand.id}`;
        
        if (stand.current_people === 1) {
            // Если перед пользователем 1 человек и уведомление еще не показывалось
            if (!localStorage.getItem(notificationKey)) {
                sendMobileNotification(stand.name);
                localStorage.setItem(notificationKey, 'true');
                console.log(`Показано уведомление для стенда: ${stand.name}`);
            } else {
                console.log(`Уведомление для стенда ${stand.name} уже было показано`);
            }
        } else {
            // Если состояние изменилось (не 1 человек), очищаем флаг уведомления
            if (localStorage.getItem(notificationKey)) {
                localStorage.removeItem(notificationKey);
                console.log(`Очищен флаг уведомления для стенда: ${stand.name}`);
            }
        }
    });
}

window.addEventListener('beforeunload', function() {
    // Очищаем все флаги уведомлений
    const allKeys = Object.keys(localStorage);
    const notificationKeys = allKeys.filter(key => key.startsWith('notified_'));
    
    notificationKeys.forEach(key => {
        localStorage.removeItem(key);
    });
    
    console.log('Очищены флаги уведомлений при выходе со страницы');
});