let token = "";
transactionList.innerHTML = "";

if (localStorage.getItem("username") != null){  
    document.getElementById("username").placeholder = localStorage.getItem("username");
    loadTransactions();
}

document.getElementById("login").addEventListener("click", () => {
    const username = document.getElementById("username").value;
    if (!username) {
        document.getElementById("response").innerText = "Введите имя!";
        return;
    }
    localStorage.setItem("username", username); 
    document.getElementById("username").placeholder = username;

    fetch("/api/auth", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ username: username })
    })
    
    .then(res => res.json())
    .then(data => {
        if (data.token) {
            token = data.token;
            localStorage.setItem("token", token);  // Сохраняем токен в localStorage
            document.getElementById("response").innerText = "Успешный вход! Токен сохранен.";
            loadTransactions();
        } else {
            document.getElementById("response").innerText = "Ошибка авторизации.";
        }
    console.log("Отправляем JSON:", JSON.stringify({ username: username }));

    });
});


document.getElementById("sendCoins").addEventListener("click", () => {
    const recipient = document.getElementById("recipient").value;
    const amount = document.getElementById("amount").value;

    if (!recipient || !amount) {
        document.getElementById("response").innerText = "Введите все данные!";
        return;
    }

    token = localStorage.getItem("token");
    if (!token) {
        alert("Вы не авторизованы!");
        return;
    }

    fetch("/api/sendCoin", {
        method: "POST",
        headers: { 
            "Content-Type": "application/json",
            "Authorization": `Bearer ${token}`  // Отправляем токен в заголовке Authorization
        },
        body: JSON.stringify({ toUser: recipient, amount: parseInt(amount) })
    })
    .then(res => res.json())
    .then(data => {
        if (data.error) {
            alert("Ошибка: " + data.error);
        } else {
            document.getElementById("response").innerHTML = 
            `Вы отправили ${data.recipient} ${data.amount} сoin. Ваш текущий баланс ${data.sender_new_balance} сoin`;
            loadTransactions();
        }
    })
});


document.getElementById("getInfo").addEventListener("click", () => {
    token = localStorage.getItem("token");
    if (!token) {
        alert("Вы не авторизованы!");
        return;
    }

    fetch("/api/info", {
        method: "GET",
        headers: { 
            "Authorization": `Bearer ${token}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert("Ошибка: " + data.error);
        } else {
            document.getElementById("response").innerHTML = `
            <p><b>Имя:</b> ${data.username}</p>
            <p><b>Баланс:</b> ${data.balance} монет</p>
            <p><b>Инвентарь:</b> ${data.inventory.map(item => `${item.ItemType} (${item.Quantity})`).join(", ")}</p>
        `;
        }
    })
    .catch(error => console.error("Ошибка загрузки данных:", error));
});

document.getElementById("buyItem").addEventListener("click", () => {
    token = localStorage.getItem("token");
    let item = document.getElementById("itemSelect").value;

    if (!token) {
        alert("Вы не авторизованы!");
        return;
    }

    fetch(`/api/buy/${item}`, {
        method: "GET",
        headers: { 
            "Authorization": `Bearer ${token}`
        }
    })
    .then(response => response.json())
    .then(data => {
        if (data.error) {
            alert("Ошибка: " + data.error);
        } else {
            alert(`Куплен товар: ${item}\nОстаток монет: ${data.balance}`);
            document.getElementById("getInfo").click();
        }
    })
    .catch(error => console.error("Ошибка:", error));
});


document.addEventListener("DOMContentLoaded", () => {
    fetch("/api/items")
        .then(response => response.json())
        .then(items => {
            let select = document.getElementById("itemSelect");
            select.innerHTML = "";
            items.forEach(item => {
                let option = document.createElement("option");
                option.value = item.name;
                option.textContent = `${item.name}: ${item.price} coin`;
                select.appendChild(option);
            });
        })
        .catch(error => console.error("Ошибка загрузки товаров:", error));
});

function loadTransactions() {
    fetch("/api/transactions", {
        headers: {
            "Authorization": "Bearer " + localStorage.getItem("token")
        }
    })
    .then(response => {
        if (!response.ok) {
            throw new Error("Ошибка загрузки");
        }
        return response.json();
    })
    .then(data => {
        const transactionList = document.getElementById("transactionList");
        transactionList.innerHTML = "";
        const currentUser = localStorage.getItem("username");

        data.forEach(tx => {
            const li = document.createElement("li");
            li.textContent = `${tx.send} -> ${tx.rec}: ${tx.amount} coin`;

            if (tx.rec === currentUser) {
                li.style.color = "green"; // Входящий перевод (получили монеты)
            } else if (tx.send === currentUser) {
                li.style.color = "red"; // Исходящий перевод (отправили монеты)
            }

            transactionList.appendChild(li);
        });
    })
    .catch(error => console.error("Ошибка загрузки транзакций:", error));
}
