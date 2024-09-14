document.addEventListener("DOMContentLoaded", function() {
    const profileIcon = document.querySelector(".profile-icon");
    const dropdownMenu = document.querySelector(".dropdown");

    profileIcon.addEventListener("click", function(event) {
        // Переключаем видимость меню
        dropdownMenu.style.display = dropdownMenu.style.display === "none" ? "block" : "none";
    });

    // Закрытие меню при клике вне его
    document.addEventListener("click", function(event) {
        if (!profileIcon.contains(event.target) && !dropdownMenu.contains(event.target)) {
            dropdownMenu.style.display = "none";
        }
    });
});