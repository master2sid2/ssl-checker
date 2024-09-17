    $(document).ready(function() {
        $("#addUserForm").on("submit", function(e) {
            e.preventDefault();

            $.ajax({
                type: "POST",
                url: "/admin/add",
                data: $(this).serialize(),
                success: function(response) {
                    if (response.error) {
                        $("#error-message").text(response.error);
                    } else {
                        window.location.href = "/admin";
                    }
                },
                error: function() {
                    $("#error-message").text("An error occurred. Please try again.");
                }
            });
        });
    });
