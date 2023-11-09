# Prediction Engine API Interaction Flow
Updated on Nov 09.

## E-commerce (Price Prediction) 
## Ben, Sarah, Emily
- **User Interaction**: Request for a price prediction of a product.
- **ML Algorithms**: Linear Regression
- **Flow**:
  - Web UI → CUDA: Formats and sends the request.
  - CUDA → CRAB: Instructs data retrieval.
  - CRAB → DAL/SQL: Fetches historical price data and user reviews.
  - DAL/SQL → CRAB: Returns the requested data.
  - CUDA: Processes the data to predict the price.
  - User: Receives the predicted price.

## Real Estate (Market Value Prediction) 
## Ben, Sarah, Emily
- **User Interaction**: Market value estimation request for a property.
- **ML Algorithms**: K-Nearest Neighbors.
- **Flow**:
  - Web UI → CUDA: Forwards the request.
  - CUDA → CRAB: Asks for real estate listings.
  - CRAB → DAL/SQL: Fetches property data.
  - DAL/SQL → CRAB: Returns the data.
  - CUDA: Predicts property value.
  - User: Obtains estimated market value.

## Job Market (Industry Trend Analysis)
## Hansi, Eni, Mat, Shiv, Matthew
- **User Interaction**: Inquiry about the demand for specific job skills.
- **ML Algorithms**: Naive Bayes Classifier.
- **Flow**:
  - Web UI → CUDA: Sends the inquiry.
  - CUDA → CRAB: Requests job listing data.
  - CRAB → DAL/SQL: Retrieves job data.
  - DAL/SQL → CRAB: Sends back the data.
  - CUDA: Analyzes for demand trends.
  - User: Gets the trend analysis.

Each flow is designed to be a seamless interaction from the user's request to the final output, leveraging machine learning algorithms for accurate and efficient predictions.
