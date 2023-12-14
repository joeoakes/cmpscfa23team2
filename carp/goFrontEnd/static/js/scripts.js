$(document).ready(function() {
  const $domainSelect = $('#domainSelect');
  const $querySelect = $('#querySelect');
  const domainToQueries = {
    "E-commerce (Price Prediction)": ["E-commerce Query 1 Used car", "E-commerce Query 2 Gold Price", "E-commerce Query 3 Silver Price"],
    "Gas Prices (Industry Trend Analysis)": ["Gas Prices Prediction for the year 2023", "Gas Prices Query 2", "Gas Prices Query 3"],
    "Job Market (Industry Trend Analysis)": ["Top 3 Tech Jobs", "Top 3 Law Jobs", "Top 3 Business Jobs"]
  };

  $domainSelect.change(function() {
    const selectedDomain = $domainSelect.val();
    const queries = domainToQueries[selectedDomain] || [];
    $querySelect.empty();
    $.each(queries, function(index, query) {
      $querySelect.append($('<option>', { value: query, text: query }));
    });
  }).change();

  $('#aiPredictionForm').submit(function(event) {
    event.preventDefault();
    const domain = $domainSelect.val();
    const queryType = $querySelect.val();

    $.ajax({
      url: `/api/predictions?domain=${encodeURIComponent(domain)}&queryType=${encodeURIComponent(queryType)}`,
      type: 'GET',
      success: function(response) {
        $('#predictionResult').empty();
        displayPredictionResults(domain, queryType, response);
      },
      error: function(xhr, status, error) {
        $('#predictionResult').text(`Error: ${error}`);
      }
    });
  });

  function displayPredictionResults(domain, queryType, response) {
    // Display skills for "Job Market (Industry Trend Analysis)"
    if (domain === "Job Market (Industry Trend Analysis)") {
      const skillsMapping = {
        "Top 3 Tech Jobs": "Skills: Software, Java, React, C++, JavaScript, DevOps, Cloud, AWS, Backend",
        "Top 3 Law Jobs": "Skills: Law, Litigation, Legal, Contract, Compliance",
        "Top 3 Business Jobs": "Skills: Management, Finance, Marketing, Sales, Microsoft Office"
      };
      $('#keywordsOutput').html(skillsMapping[queryType] || "");

      if (response.job_listings) {
        $('#predictionResult').append(formatJobListings(response.job_listings));
      } else {
        $('#predictionResult').append($('<p>').text('No job listings found.'));
      }
    } else {
      // Display prediction text for other domains
      if (response.prediction_info) {
        $('#predictionResult').append($('<p>').text(response.prediction_info));
      }
      $('#keywordsOutput').empty();
    }

    // Handle image path for predictions that include a visual component
    if (response.image_path) {
      displayImage(response.image_path);
    }
  }

  function displayImage(imagePath) {
    var image = $('<img>', {
      src: imagePath,
      alt: 'Prediction Result',
      style: 'max-width: 100%; height: auto;'
    }).on('error', function() {
      $('#predictionResult').append($('<p>').text("Error loading prediction image."));
    });
    $('#predictionResult').append(image);
  }

  function formatJobListings(jobListings) {
    let formattedListings = '';
    jobListings.forEach(listing => {
      const descriptionListItems = listing.description.split('\n').map(line => `<li>${line}</li>`).join('');
      formattedListings += `
        <div class="job-listing">
          <h3 class="job-title text-center">${listing.title}</h3>
          <p class="text-center"><a href="${listing.url}" target="_blank" class="job-link">View Listing</a></p>
          <p class="text-center"><strong>Company:</strong> ${listing.company}</p>
          <p class="text-center"><strong>Location:</strong> ${listing.location}</p>
          <p class="text-center"><strong>Salary:</strong> ${listing.salary}</p>
          <p class="description-label"><strong>Description:</strong></p>
          <ul class="description-list">${descriptionListItems}</ul>
        </div>`;
    });
    return formattedListings;
  }

// Additional functionalities (Settings, Login, Logout, etc.) can be added here as needed
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

});