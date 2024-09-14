document.getElementById('logout-link').addEventListener('click', function(event) {
    event.preventDefault(); // Предотвращаем переход по ссылке
    
    // Создаем и отправляем POST-запрос на сервер
    fetch('/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: '', // Тело запроса можно оставить пустым для выхода
    }).then(function(response) {
        if (response.ok) {
            window.location.href = '/login'; // Перенаправляем пользователя после выхода
        }
    }).catch(function(error) {
        console.error('Ошибка при выходе:', error);
    });
});