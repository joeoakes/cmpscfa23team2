$(document).ready(function() {
  // When the 'Learn More' button is clicked, show the login modal
  $('#learnMoreBtn').click(function() {
    $('#loginModal').modal('show');
  });

  // Handle login form submission
  $('#loginForm').submit(function(event) {
    event.preventDefault();
    var username = $('#username').val();
    var password = $('#password').val();

    // AJAX call to the Go backend for login
    $.ajax({
      url: '/login',
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify({ username: username, password: password }),
      success: function(response) {
        // Assuming the response contains a token
        onLoginSuccess(response.token);
      },
      error: function(xhr, status, error) {
        console.error('Login failed:', error);
        alert('Login failed: ' + error);
      }
    });
  });

  // Attach event listeners for settings buttons
  $('.btn-settings').click(function() {
    var action = $(this).text().trim();
    var module = $(this).closest('.tab-pane').attr('id');
    performAction(module, action);
  });

  // Logout functionality
  $('#logoutBtn').click(function() {
    logout();
  });
});

// Function to handle settings actions
function performAction(module, action) {
  console.log(module + ' action:', action);
  // Implement AJAX calls or other functionality based on the action
  // Example AJAX call:
  $.ajax({
    url: '/settings/' + module + '/' + action.toLowerCase(),
    type: 'POST',
    success: function(response) {
      console.log(response.message);
    },
    error: function(xhr, status, error) {
      console.error(module + ' ' + action + ' failed:', error);
    }
  });
}

// Function to run on login success
function onLoginSuccess(token) {
  console.log('Login successful, token:', token);
  $('#loginModal').modal('hide');
  localStorage.setItem('token', token);
  // Redirect to dashboard or update UI
}

// Function to logout
function logout() {
  localStorage.removeItem('token');
  // Redirect to the home page or update UI
  window.location.href = '/';
}
