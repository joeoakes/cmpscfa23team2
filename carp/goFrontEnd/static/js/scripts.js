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

  // Handle AI Prediction Form submission
  $('#aiPredictionForm').submit(function (event) {
    event.preventDefault();

    var domain = $('#domainSelect').val();
    var query = $('#querySelect').val();

    // Define the base path for images
    var basePath = 'static/Assets/MachineLearning/';


    // Define a variable for the image path
    var imagePath, imageDescription;


    // Check domain and query combinations
    if (domain === "E-commerce") {
      switch (query) {
        case "E-commerce Query 1 Used car":
          imagePath = basePath + 'KNN/e-commerceQuery1Pic.png';
          imageDescription = 'Description for E-commerce Query 1...';
          break;
        case "E-commerce Query 2 Gold Price":
          imagePath = basePath + 'KNN/e-commerceQuery2Pic.png';
          imageDescription = 'Description for E-commerce Query 2...';
          break;
        case "E-commerce Query 3 Silver Price":
          imagePath = basePath + 'KNN/e-commerceQuery3Pic.png';
          imageDescription = 'Description for E-commerce Query 3...';
          break;
      }
    } else if (domain === "RealEstate") {
      switch (query) {
        case "Real Estate Query 1 Philadelphia":
          imagePath = basePath + 'LinearRegression/realEstateQuery1Pic.png';
          imageDescription = 'Description for Real Estate Query 1...';
          break;
        case "Real Estate Query 2 New York":
          imagePath = basePath + 'LinearRegression/realEstateQuery2Pic.png';
          imageDescription = 'Description for Real Estate Query 2...';
          break;
        case "Real Estate Query 3 California":
          imagePath = basePath + 'LinearRegression/realEstateQuery3Pic.png';
          imageDescription = 'Description for Real Estate Query 3...';
          break;
      }
    }

    // Display the image if a path is set, otherwise show the simulation result
    if (imagePath) {
      $('#predictionResult').html(`
            <div id="image-description">${imageDescription}</div>
            <img src="${imagePath}" alt="Result" style="width: 100%; height: auto; object-fit: contain;">
        `);
    } else {
      // Simulation of the prediction result for other selections or Job Market
      console.log("Selected Query: " + query + ", Domain: " + domain);
      $('#predictionResult').text("Simulation of the prediction result.");
    }
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