// Mapping of domains to queries
const domainToQueries = {
  "E-commerce": ["E-commerce Query 1 Used car", "E-commerce Query 2 Gold Price", "E-commerce Query 3 Silver Price"],
  RealEstate: ["Real Estate Query 1 Philadelphia", "Real Estate Query 2 New York", "Real Estate Query 3 California"],
  JobMarket: ["Job Market Query 1", "Job Market Query 2", "Job Market Query 3"],
  // Add other domains and queries as necessary
};

// Cache the selectors for domain and query dropdowns
const $domainSelect = $('#domainSelect');
const $querySelect = $('#querySelect');

$(document).ready(function() {
  // Function to update the query select based on the selected domain
  function updateQuerySelect() {
    const selectedDomain = $domainSelect.val();
    const queries = domainToQueries[selectedDomain] || [];

    $querySelect.empty();
    $.each(queries, function (index, query) {
      $querySelect.append($('<option>', {
        value: query,
        text: query
      }));
    });
  }

  // Initialize the query selection based on the initial domain selection
  updateQuerySelect();

  // Update the query dropdown when the domain selection changes
  $domainSelect.change(updateQuerySelect);

  // When the 'Learn More' button is clicked, show the login modal
  $('#learnMoreBtn').click(function () {
    $('#loginModal').modal('show');
  });

  // Handle login form submission
  $('#loginForm').submit(function (event) {
    event.preventDefault();
    var username = $('#username').val();
    var password = $('#password').val();

    // AJAX call to the Go backend for login
    $.ajax({
      url: '/login',
      type: 'POST',
      contentType: 'application/json',
      data: JSON.stringify({username: username, password: password}),
      success: function (response) {
        onLoginSuccess(response.token);
      },
      error: function (xhr, status, error) {
        console.error('Login failed:', error);
        alert('Login failed: ' + error);
      }
    });
  });

  // Attach event listeners for settings buttons
  $('.btn-settings').click(function () {
    var action = $(this).text().trim();
    var module = $(this).closest('.tab-pane').attr('id');
    performAction(module, action);
  });

  // Logout functionality
  $('#logoutBtn').click(function () {
    logout();
  });

  $('#aiPredictionForm').submit(function(event) {
    event.preventDefault();

    var domain = $('#domainSelect').val();
    var queryType = $('#querySelect').val();

    $.ajax({
      url: '/api/predictions?domain=' + encodeURIComponent(domain) + '&queryType=' + encodeURIComponent(queryType),
      type: 'GET',
      success: function(response) {
        // Check if the response contains valid prediction info
        if (response && response.prediction_info) {
          var predictionInfo = response.prediction_info;
          // Check if the prediction info is a file path
          if (/\.(jpg|png|gif)$/.test(predictionInfo)) {
            // Display the image
            var imagePath = 'static/Assets/MachineLearning/' + predictionInfo;
            $('#predictionResult').html('<img src="' + imagePath + '" alt="Prediction Result" style="max-width: 100%; height: auto;">');
          } else {
            // Display text prediction
            $('#predictionResult').text(predictionInfo);
          }
        } else {
          // Handle missing or invalid prediction info
          $('#predictionResult').text("No prediction data available for the selected query.");
        }
      },
      error: function(xhr, status, error) {
        // Display a user-friendly error message
        console.error('Error fetching prediction:', error);
        $('#predictionResult').text("An error occurred while fetching the prediction. Please try again later.");
      }
    });
  });

// Function to handle settings actions
    function performAction(module, action) {
      $.ajax({
        url: '/settings/' + module + '/' + action.toLowerCase(),
        type: 'POST',
        success: function (response) {
          console.log(response.message);
        },
        error: function (xhr, status, error) {
          console.error(module + ' ' + action + ' failed:', error);
        }
      });
    }

// Function to run on login success
    function onLoginSuccess(token) {
      $('#loginModal').modal('hide');
      localStorage.setItem('token', token);
    }

// Function to logout
    function logout() {
      localStorage.removeItem('token');
      window.location.href = '/';
    }
  });

