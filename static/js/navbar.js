document.addEventListener('DOMContentLoaded', function() {
    var backToHomeLink = document.getElementById('back-to-home');
    var currentUrl = window.location.pathname;

    // Скрыть ссылку, если текущий URL - это страница /home
    if (currentUrl === '/home') {
        if (backToHomeLink) {
            backToHomeLink.style.display = 'none';
        }
    }
});