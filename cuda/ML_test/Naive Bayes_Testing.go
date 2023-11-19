package ML_test

//
//func main1() {
//	// Create a Naive Bayesian Classifier
//	classifier := bayesian.NewClassifier(bayesian.MultinomialTf)
//
//	// Train the classifier with job descriptions and their corresponding skill demands
//	trainClassifier(classifier)
//
//	// Sample job description to classify
//	jobDescription := "We are looking for a data scientist with expertise in machine learning and Python."
//
//	// Classify the job description
//	probabilities, bestClass, certain := classifier.Classify(tokenizeText(jobDescription))
//
//	// Display the classification results
//	fmt.Println("Job Description:", jobDescription)
//	fmt.Println("Best Class:", bestClass)
//	fmt.Println("Probabilities:", probabilities)
//	fmt.Println("Is certain:", certain)
//}
//
//// Train the classifier with sample data
//func trainClassifier(classifier bayesian.Classifier) {
//	// Sample job descriptions and their associated classes
//	jobData := []bayesian.Document{
//		bayesian.NewDocument("Data Scientist", "data scientist", "machine learning", "Python"),
//		bayesian.NewDocument("Software Engineer", "software engineer", "Java", "web development"),
//		bayesian.NewDocument("Database Administrator", "database administrator", "SQL", "database management"),
//	}
//
//	// Learn from the training data
//	classifier.Learn(jobData...)
//}
//
//// Tokenize text into words
//func tokenizeText(text string) []string {
//	// This is a simple tokenizer, you may want to implement a more advanced one
//	// based on your requirements.
//	// You can also preprocess the text to remove stop words and perform stemming.
//	return bayesian.Tokenize(text)
//}
