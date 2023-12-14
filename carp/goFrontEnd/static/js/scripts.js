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


  $('#aiPredictionForm').submit(function(event) {
    event.preventDefault();

    var domain = $('#domainSelect').val();
    var queryType = $('#querySelect').val();

    $.ajax({
      url: '/api/predictions?domain=' + encodeURIComponent(domain) + '&queryType=' + encodeURIComponent(queryType),
      type: 'GET',
      success: function(response) {
        console.log("Response received:", response);
        $('#predictionResult').empty();

        // Display skills if available and if the domain is 'Job Market (Industry Trend Analysis)'
        if (domain === 'Job Market (Industry Trend Analysis)' && response.skills) {
          $('#predictionResult').append($('<h2>').text("Skills: " + response.skills));
        }

        // Display specific job if available and if the domain is 'Job Market (Industry Trend Analysis)'
        if (domain === 'Job Market (Industry Trend Analysis)' && response.specificJob) {
          const specificJob = formatSpecificJob(response.specificJob);
          $('#predictionResult').append(specificJob);
        }

        // Display general job listings if the domain is 'Job Market (Industry Trend Analysis)'
        if (domain === 'Job Market (Industry Trend Analysis)' && response.job_listings && response.job_listings.length > 0) {
          const jobListings = formatJobListing(response.job_listings);
          $('#predictionResult').append(jobListings);
        } else if (domain !== 'Job Market (Industry Trend Analysis)') {
          // Handle other predictions
          if (response.prediction_info) {
            var predictionText = $('<p>').text(response.prediction_info);
            $('#predictionResult').append(predictionText);
          }
        } else {
          $('#predictionResult').append($('<p>').text("No job listings found for the selected query."));
        }

        // Handle image path for predictions that include a visual component
        if (response.image_path) {
          var imagePath = response.image_path;
          var image = $('<img>', {
            src: imagePath,
            alt: 'Prediction Result',
            style: 'max-width: 100%; height: auto;'
          }).on('error', function() {
            $('#predictionResult').append($('<p>').text("Error loading prediction image."));
          });
          $('#predictionResult').append(image);
        }
      },
      error: function(xhr, status, error) {
        $('#predictionResult').text("An error occurred while fetching the prediction.");
      }
    });
  });


  function formatSpecificJob(job) {
    let jobHTML = `
    <div class="specific-job">
      <h3 class="job-title">${job.title}</h3>
      <p><strong>Company:</strong> ${job.company}</p>
      <p><strong>Location:</strong> ${job.location}</p>
      <p><strong>Salary:</strong> ${job.salary}</p>
      <p><strong>Description:</strong> ${job.description.replace(/\n/g, '<br>')}</p>
      <a href="${job.url}" target="_blank" class="btn btn-primary">View Listing</a>
    </div>`;
    return jobHTML;
  }

  function formatJobListing(jobListings) {
    let formattedListings = '';
    jobListings.forEach(listing => {
      formattedListings += `
    <div class="job-listing">
      <h3 class="job-title">${listing.title}</h3>
      <a href="${listing.url}" target="_blank" class="job-link">View Listing</a>
      <p class="job-info">Company: ${listing.company}</p>
      <p class="job-info">Location: ${listing.location}</p>
      <p class="job-info">Salary: ${listing.salary}</p>
      <p class="job-description">${listing.description.replace(/\n/g, '<br>')}</p>
    </div>`;
    });
    return formattedListings;
  }
});

// // Attach event listeners for settings buttons
// $('.btn-settings').click(function () {
//   var action = $(this).text().trim();
//   var module = $(this).closest('.tab-pane').attr('id');
//   performAction(module, action);
// });
//
// // Logout functionality
// $('#logoutBtn').click(function () {
//   logout();
// });
// // Function to handle settings actions
// function performAction(module, action) {
//   $.ajax({
//     url: '/settings/' + module + '/' + action.toLowerCase(),
//     type: 'POST',
//     success: function (response) {
//       console.log(response.message);
//     },
//     error: function (xhr, status, error) {
//       console.error(module + ' ' + action + ' failed:', error);
//     }
//   });
// }
//
// // Function to run on login success
// function onLoginSuccess(token) {
//   $('#loginModal').modal('hide');
//   localStorage.setItem('token', token);
// }
//
// // Function to logout
// function logout() {
//   localStorage.removeItem('token');
//   window.location.href = '/';
// }
