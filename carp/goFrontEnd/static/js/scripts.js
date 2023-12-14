$(document).ready(function() {
  const $domainSelect = $('#domainSelect');
  const $querySelect = $('#querySelect');
  const domainToQueries = {
    "Job Market": ["Top 3 Tech Jobs", "Top 3 Law Jobs", "Top 3 Business Jobs"],
    "Gas Prices": ["Gas Prices Prediction 2023", "Gas Prices Prediction 2024", "Gas prices target prediction for years similar to 2023 prediction"],
    "Airfare Prices": ["Airfare Prices Prediction 2024", "Airfare Prices Prediction 2025", "Airfare Prices Prediction 2030"],
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
    // Clear previous results for both predictions and keywords
    $('#predictionResult').empty();
    $('#keywordsOutput').empty(); // Clear the keywords for every new query

    // Display skills for "Job Market (Industry Trend Analysis)"
    if (domain === "Job Market") {
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
      // Since we're not in the "Job Market" domain, ensure any previous skills are not shown
      $('#keywordsOutput').empty();

      if (queryType === "Gas prices target prediction for years similar to 2023 prediction") {
        // Specific logic to handle tabular data for Gas Prices target prediction
        if (response.prediction_info) {
          $('#predictionResult').append(createTableFromPrediction(response.prediction_info));
        }
      } else {
        // Other domains
        if (response.prediction_info) {
          $('#predictionResult').append($('<p>').text(response.prediction_info));
        }
      }

      // Handle image path for predictions that include a visual component
      if (response.image_path) {
        displayImage(response.image_path);
      }
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
  function createTableFromPrediction(predictionInfo) {
    // Split the input into lines and then process
    let lines = predictionInfo.trim().split('\n');
    let tableContainer = $('<div class="table-container">');
    let table = $('<table class="table table-striped">');
    let tbody = $('<tbody>');

    // Data rows - start from the second line, exclude the first and last lines
    lines.slice(1, -1).forEach((line, index) => {
      if (line.trim().length === 0) return; // Skip empty lines
      let rowData = line.split(/\s+/).filter(cell => cell.trim().length > 0);
      let row = $('<tr>');
      rowData.forEach(cellText => {
        row.append($('<td>').text(cellText));
      });
      tbody.append(row);
    });
    table.append(tbody);
    tableContainer.append(table);

    // Summary line - kept outside the table for distinct styling
    let summary = $('<div class="summary-text">').text(lines[lines.length - 1]);
    tableContainer.append(summary);

    // Return the container with the table and summary
    return tableContainer;
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