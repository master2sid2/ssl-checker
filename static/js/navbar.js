document.addEventListener('DOMContentLoaded', function() {
    var backToHomeLink = document.getElementById('back-to-home');
    var currentUrl = window.location.pathname;

    if (currentUrl === '/home') {
        if (backToHomeLink) {
            backToHomeLink.style.display = 'none';
        }
    }
});