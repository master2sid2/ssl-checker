document.addEventListener("DOMContentLoaded", function() {
    const profileIcon = document.querySelector(".profile-icon");
    const dropdownMenu = document.querySelector(".dropdown");

    profileIcon.addEventListener("click", function(event) {
        dropdownMenu.style.display = dropdownMenu.style.display === "none" ? "block" : "none";
    });

    document.addEventListener("click", function(event) {
        if (!profileIcon.contains(event.target) && !dropdownMenu.contains(event.target)) {
            dropdownMenu.style.display = "none";
        }
    });
});