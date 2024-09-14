    $(document).ready(function() {
        $("#addUserForm").on("submit", function(e) {
            e.preventDefault(); // Предотвращаем отправку формы по умолчанию

            $.ajax({
                type: "POST",
                url: "/admin/add", // URL для отправки данных
                data: $(this).serialize(), // Собираем данные формы
                success: function(response) {
                    if (response.error) {
                        $("#error-message").text(response.error);
                    } else {
                        window.location.href = "/admin"; // Перенаправляем на /admin при успехе
                    }
                },
                error: function() {
                    $("#error-message").text("An error occurred. Please try again.");
                }
            });
        });
    });
