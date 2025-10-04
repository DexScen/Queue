document.addEventListener('DOMContentLoaded', async function () {
    const username = localStorage.getItem('username');

    try {
        const standsData = await loadAllStands();
        
        if (standsData) {
            renderAllStands(standsData);
        } else {
            console.error('не удалось загрузить данные стендов');
        }

    } catch (error) {
        console.error('Ошибка при загрузке страницы:', error);
    }
});

async function loadAllStands() {
    try {
        const response = await fetch(`http://localhost:8080/games`);
        
        if (response.ok) {
            const standsData = await response.json();
            return standsData;
        } else {
            console.error('Ошибка:', response.status);
            return null;
        }
    } catch (error) {
        console.error('Ошибка при загрузке данных стендов:', error);
        return null;
    }
}

function renderAllStands(stands) {
    const container = document.querySelector('.div-card-stand');
    container.innerHTML = '';

    if (!stands || !Array.isArray(stands) || stands.length === 0) {
        container.innerHTML = '<p>Нет доступных стендов</p>';
        return;
    }

    stands.forEach(stand => {
        const waitTimeMinutes = Math.ceil(stand.current_people * stand.duration_seconds / 60);
        
        const standCard = document.createElement('article');
        standCard.className = 'card-stand';
        standCard.innerHTML = `
            <section>
                <h2 class="text-title">${stand.name}</h2>
                <p class="text-description">${stand.description}</p>
            </section>
            <section class="div-container">
                <div class="div-info">
                    <p>Очередь: ${stand.current_people}/${stand.max_slots}</p>
                    <p>Время ~ ${waitTimeMinutes} мин</p>
                </div>
                <button class="btn-standart" data-stand-id="${stand.id}">
                    Записаться
                </button>
            </section>
            <section class="div-img">
                <img class="img-stand" src="../assets/photo1.jpg" alt="${stand.name}">
            </section>
        `;

        container.appendChild(standCard);
    });

    document.querySelectorAll('.btn-standart').forEach(button => {
        button.addEventListener('click', async function() {
            const standId = this.getAttribute('data-stand-id');
            await signUpForStand(standId);
        });
    });
}

async function signUpForStand(standId) {
    const username = localStorage.getItem('username');

    if (!username) {
        alert('Пожалуйста, войдите в систему');
        return;
    }

    try {
        const authResponse = await fetch(`http://localhost:8080/auth/${username}`);
        
        if (!authResponse.ok) {
            throw new Error('Ошибка при получении user_id');
        }
        
        const userData = await authResponse.json();
        const userId = userData.id;

        const queuesResponse = await fetch(`http://localhost:8080/queue/${username}`);
        if (queuesResponse.ok) {
            const userQueues = await queuesResponse.json();
            
            if (userQueues && Array.isArray(userQueues)) {
                const alreadyRegistered = userQueues.some(queue => queue.id === parseInt(standId));
                if (alreadyRegistered) {
                    alert('Вы уже записаны в эту очередь!');
                    return;
                }
            }
        }

        const signupResponse = await fetch(`http://localhost:8080/add`, {
            method: 'POST',
            headers: {
                'Content-Type': 'application/json'
            },
            body: JSON.stringify({
                user_id: userId,
                game_id: parseInt(standId)
            })
        });

        console.log('Статус ответа:', signupResponse.status);
        
        if (signupResponse.ok) {
            alert('Вы успешно записались в очередь!');
            const standsData = await loadAllStands();
            if (standsData) {
                renderAllStands(standsData);
            }
        } else {
            const errorText = await signupResponse.text();
            console.error('Ошибка сервера:', errorText);
            alert(`Не удалось записаться в очередь. Ошибка: ${signupResponse.status}`);
        }

    } catch (error) {
        console.error('Ошибка при записи:', error);
        alert('Произошла ошибка при подключении к серверу');
    }
}