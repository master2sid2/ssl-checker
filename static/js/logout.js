document.getElementById('logout-link').addEventListener('click', function(event) {
    event.preventDefault();
    
    fetch('/logout', {
        method: 'POST',
        headers: {
            'Content-Type': 'application/x-www-form-urlencoded',
        },
        body: '',
    }).then(function(response) {
        if (response.ok) {
            window.location.href = '/login';
        }
    }).catch(function(error) {
        console.error('Ошибка при выходе:', error);
    });
});