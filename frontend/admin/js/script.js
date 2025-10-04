document.addEventListener("DOMContentLoaded", async () => {
  const container = document.querySelector("main .container");
  const gameId = 1;
  const headerInfo = document.querySelector(".div-flex-jestify-between p");
  const GAME_DURATION = 600; // секунд

  async function loadPlayers() {
    try {
      const response = await fetch(`http://localhost:8080/players/${gameId}`);
      if (!response.ok) throw new Error("Ошибка загрузки игроков");

      const data = await response.json();
      container.innerHTML = "";

      const playersCount = data ? data.length : 0;
      const waitMinutes = Math.round((playersCount * GAME_DURATION) / 60);
      headerInfo.innerHTML = `${playersCount} чел <br> ~ ${waitMinutes} мин`;

      if (playersCount === 0) {
        const emptyMsg = document.createElement("p");
        emptyMsg.textContent = "Очередь пуста";
        emptyMsg.style.textAlign = "center";
        emptyMsg.style.color = "#666";
        container.appendChild(emptyMsg);
        return;
      }

      data.forEach((player, index) => {
        const card = document.createElement("div");
        card.classList.add("card-profile", "back-color-white");

        const actionBtn = index === 0
          ? `<button class="btn-standart" data-user-id="${player.id}">закончил</button>`
          : `<button style="visibility: hidden" class="btn-standart" data-user-id="${player.id}">закончил</button>`;

        card.innerHTML = `
          <p>${index + 1}</p>
          <h3>${player.login}</h3>
          ${actionBtn}
          <div>
            <button class="btn-delete" data-user-id="${player.id}">удалить</button>
          </div>
        `;

        container.appendChild(card);

        // Обработчики кнопок сразу при создании карточки
        const deleteBtn = card.querySelector(".btn-delete");
        deleteBtn.addEventListener("click", () => removePlayer(player.id, card));

        if (index === 0) {
          const finishBtn = card.querySelector(".btn-standart");
          finishBtn.addEventListener("click", () => removePlayer(player.id, card));
        }
      });

    } catch (error) {
      console.error(error);
    }
  }

  async function removePlayer(userId, card) {
    try {
      const resp = await fetch("http://localhost:8080/remove", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ user_id: Number(userId), game_id: gameId })
      });

      if (resp.ok) {
        // Анимация удаления
        card.classList.add('smoke-vanish'); // добавь CSS-анимацию
        setTimeout(() => {
          card.remove();
          updateUserCount();
        }, 800); // совпадает с длительностью анимации
      } else {
        const errText = await resp.text();
        alert("Ошибка при удалении: " + errText);
      }
    } catch (error) {
      console.error(error);
      alert("Ошибка соединения с сервером");
    }
  }

  function updateUserCount() {
    const cards = document.querySelectorAll('.card-profile');
    const waitMinutes = Math.round((cards.length * GAME_DURATION) / 60);
    headerInfo.innerHTML = `${cards.length} чел <br> ~ ${waitMinutes} мин`;
  }

  loadPlayers();
  setInterval(loadPlayers, 1000); // автообновление каждые 1 секунд
});
