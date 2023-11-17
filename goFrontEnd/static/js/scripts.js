$(document).ready(function () {
  $("#learnMoreBtn").click(function () {
    $("#loginModal").modal("show");
  });

  $("#loginForm").on("submit", function (event) {
    event.preventDefault();
    $.ajax({
      url: "/login",
      method: "POST",
      contentType: "application/json",
      data: JSON.stringify({
        username: $("#username").val(),
        password: $("#password").val(),
      }),
      success: function (response) {
        console.log("Login successful:", response);
        $("#loginModal").modal("hide");
        // You can store the token in localStorage and redirect to a dashboard etc.
        // localStorage.setItem('token', response.token);
        // window.location.href = '/dashboard'; // The dashboard URL if exists
      },
      error: function (xhr, status, error) {
        console.error("Login failed:", error);
      },
    });
  });
});

$(document).ready(function () {
  $("#learnMoreBtn").click(function () {
    $("#loginModal").modal("show");
  });

  $("#loginForm").on("submit", function (event) {
    event.preventDefault();
    if ($("#username").val() === "test" && $("#password").val() === "test") {
      // Dummy admin login success
      onLoginSuccess({ token: "dummy-token" });
    } else {
      alert("Incorrect username or password.");
    }
  });
});

function onLoginSuccess(response) {
  console.log("Login successful:", response);
  $("#loginModal").modal("hide");
  $("#adminPage").show();
  localStorage.setItem("token", response.token);
}

function logout() {
  localStorage.removeItem("token");
  $("#adminPage").hide();
  window.location.href = "/";
}

// Add a successful login callback
// function onLoginSuccess(response) {
//   console.log("Login successful:", response);

//   $("#loginModal").modal("hide");
//   // Display the admin page section

//   $("#adminPage").show();

//   // Assuming you have a token in the response
//   localStorage.setItem("token", response.token);

//   // Redirect to the admin page if it's a separate HTML file
//   window.location.href = "/admin.html"; // The admin page URL if exists
// }
