// Mapping of domains to queries
const domainToQueries = {
  "E-commerce (Price Prediction)": ["E-commerce Query 1 Used car", "E-commerce Query 2 Gold Price", "E-commerce Query 3 Silver Price"],
  "Gas Prices (Industry Trend Analysis)": ["Gas Prices Prediction for the year 2023", "Gas Prices Query 2", "Gas Prices Query 3"],
  "Job Market (Industry Trend Analysis)": ["Top 3 Tech Jobs with most demand skills", "Top 3 Law Jobs with most demand skills", "Top 3 Bus Jobs with most demand skills"],
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
        console.log("Response received:", response);
        // Clear previous results
        $('#predictionResult').empty();

        // If the domain is 'Job Market (Industry Trend Analysis)', format the job listings
        if (domain === 'Job Market (Industry Trend Analysis)' && response.prediction_info) {
          const jobListings = formatJobListing(response.prediction_info);
          $('#predictionResult').html(jobListings); // Use .html() to set the HTML content
        } else {
          // Handle other predictions
          if (response.prediction_info) {
            var predictionText = $('<p>').text(response.prediction_info);
            $('#predictionResult').append(predictionText);
          }
        }

        // Handle image path for predictions that include a visual component
        if (response.image_path) {
          var imagePath = response.image_path;
          console.log("Image path:", imagePath);
          var image = $('<img>', {
            src: imagePath,
            alt: 'Prediction Result',
            style: 'max-width: 100%; height: auto;'
          });

          // Add error handling for the image
          image.on('error', function() {
            console.error('Error loading image:', imagePath);
            $('#predictionResult').append($('<p>').text("Error loading prediction image."));
          });

          $('#predictionResult').append(image);
        } else {
          // Display message if no prediction info is available
          $('#predictionResult').text("No prediction data available for the selected query.");
        }
      },
      error: function(xhr, status, error) {
        console.error('Error fetching prediction:', error);
        $('#predictionResult').text("An error occurred while fetching the prediction. Please try again later.");
      }
    });
  });


  function formatJobListing(predictionInfo) {
    let formattedListings = '';

    // Handle 'Most Demand Skills' separately
    if (predictionInfo.includes('Most Demand Skills:')) {
      formattedListings += `<h2>${predictionInfo.substring(0, predictionInfo.indexOf('Top Jobs for'))}</h2>`;
    }

    // Split the predictionInfo by job listings
    const jobListings = predictionInfo.split('Job Title:').slice(1); // Skip the first empty split

    jobListings.forEach(listing => {
      const parts = listing.split(', ').map(part => part.trim());

      // Extract the job title
      const jobTitle = parts[0].split('\n')[0];

      // Construct the job listing HTML
      let listingHTML = `<div class="job-listing"><h3 class="job-title">${jobTitle}</h3>`;

      let isSalaryInfo = false;
      let isDescription = false;

      // Process other parts of the listing
      parts.forEach(part => {
        if (part.startsWith('URL:')) {
          const url = part.replace('URL:', '').trim().split(' ')[0];
          listingHTML += `<a href="${url}" target="_blank" class="job-link">View Listing</a>`;
        } else if (part.startsWith('Company:') || part.startsWith('Location:')) {
          listingHTML += `<p class="job-info">${part}</p>`;
        } else if (part.startsWith('Salary:')) {
          listingHTML += `<p class="job-info">${part}</p>`;
          isSalaryInfo = true;
        } else if (isSalaryInfo && !isDescription) {
          // First line after salary info is the start of the description
          isDescription = true;
          listingHTML += `<p class="job-description">${part.replace(/\n/g, '<br>')}</p>`;
        } else if (isDescription) {
          // Continue appending description lines
          listingHTML += `${part.replace(/\n/g, '<br>')}</p>`;
        }
      });

      listingHTML += '</div>';
      formattedListings += listingHTML;
    });

    return formattedListings;
  }





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

