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

  // Handle AI Prediction Form submission
  $('#aiPredictionForm').submit(function(event) {
    event.preventDefault();

    console.log("Form submission handler called");

    var query = $('#querySelect').val();
    var domain = $('#domainSelect').val();

    console.log("Selected Query: " + query + ", Domain: " + domain);

    var results = {
      'query1': {
        'finance': 'Result for Query 1 and Finance',
        'healthcare': 'Result for Query 1 and Healthcare',
        'technology': 'Result for Query 1 and Technology'
      },
      'query2': {
        'finance': 'Result for Query 2 and Finance',
        'healthcare': 'Result for Query 2 and Healthcare',
        'technology': 'Result for Query 2 and Technology'
      },
      'query3': {
        'finance': 'Result for Query 3 and Finance',
        'healthcare': 'Result for Query 3 and Healthcare',
        'technology': 'Result for Query 3 and Technology'
      }
      // Add more combinations as needed
    };

    var resultText = results[query][domain] || 'No result found';
    $('#predictionResult').text(resultText);

    console.log("Result Text: " + resultText);
  });
});

// Function to handle settings actions
function performAction(module, action) {
  console.log(module + ' action:', action);
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
}

// Function to logout
function logout() {
  localStorage.removeItem('token');
  window.location.href = '/';
}