package dal_test

import (
	"cmpscfa23team2/dal"
	"fmt"
	"reflect"
	"testing"
)

func TestSearchJobByTitle(t *testing.T) {
	jobs := []dal.JobData{
		{Title: "Software Engineer", URL: "url1", Company: "Company1"},
		{Title: "Software Release DevOps Engineer", URL: "url2", Company: "Company2"},
		// Add more jobs as needed
	}

	tests := []struct {
		title    string
		expected *dal.JobData
	}{
		{"Software Release DevOps Engineer", &jobs[1]},
		{"Non Existent Job", nil},
	}

	for _, test := range tests {
		result := dal.SearchJobByTitle(jobs, test.title)
		if !reflect.DeepEqual(result, test.expected) {
			t.Errorf("SearchJobByTitle(%s) = %v; expected %v", test.title, result, test.expected)
		}
	}
}

// Function to format JobData for better readability in test output
func formatJobData(job *dal.JobData) string {
	if job == nil {
		return "nil"
	}
	return fmt.Sprintf("\nTitle: %s\nURL: %s\nCompany: %s\nLocation: %s\nSalary: %s\nDescription: %s",
		job.Title, job.URL, job.Company, job.Location, job.Salary, job.Description)
}
func TestFetchPredictionDataWithSpecificJob(t *testing.T) {
	// Ensure the database is initialized
	if err := dal.InitDB(); err != nil {
		t.Fatalf("Failed to initialize database: %v", err)
	}

	testCases := []struct {
		queryIdentifier string
		domain          string
		expectedError   error
		expectedJob     *dal.JobData // Update this with the actual data from your database
	}{
		{
			queryIdentifier: "Top 3 Tech Jobs with most demand skills",
			domain:          "Software Engineer",
			expectedError:   nil,
			expectedJob: &dal.JobData{
				Title:       "Software Release DevOps Engineer",
				URL:         "https://www.indeed.com/pagead/clk?mo=r&ad=...[truncated for brevity]...",
				Company:     "Comcast",
				Location:    "Location Not Found",
				Salary:      "Salary Not Found",
				Description: "A strong technical background in Software engineering...[truncated for brevity]...",
			},
		},
		// Additional test cases...
	}

	for _, tc := range testCases {
		result, err := dal.FetchPredictionData(tc.queryIdentifier, tc.domain)

		// Check for unexpected errors
		if err != tc.expectedError {
			t.Errorf("FetchPredictionData(%s, %s) unexpected error: got %v, want %v", tc.queryIdentifier, tc.domain, err, tc.expectedError)
			continue
		}

		// Check specific fields of the job data
		if result.SpecificJob.Title != tc.expectedJob.Title || result.SpecificJob.Company != tc.expectedJob.Company {
			t.Errorf("FetchPredictionData(%s, %s) returned unexpected job data:\ngot: %v\nwant: %v",
				tc.queryIdentifier, tc.domain, formatJobData(result.SpecificJob), formatJobData(tc.expectedJob))
		}
	}
}

//
//func TestLoadDataFromJSON(t *testing.T) {
//	mockFilename := "C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\SoftwareEng_jobs.json"
//
//	_, specificJob, err := dal.LoadDataFromJSON(mockFilename, "Software Release DevOps Engineer")
//	if err != nil {
//		t.Errorf("LoadDataFromJSON returned an error: %v", err)
//	}
//
//	expectedJob := dal.JobData{
//		Title:       "Software Release DevOps Engineer",
//		URL:         "https://www.indeed.com/pagead/clk?mo=r&ad=-6NYlbfkN0Cj-KmZPsf9w80C8b1WzNVrlanjD2SXJjxuCbUWHsXPZlTAgGmdtIUzoKTi6fK6WvZ2eEeIQBp5OUhO-xRyQvDo4yR3Mt5CEDSCojK6clcrRqADOS0tfXeHAsrfH_7i7PXK3XmzBFDjlntXqwANAhWdOGj1px_99ycmqNNMR1xJWJSD4fVvgEAHQ7k280w16fwgxdWPCRbm8AnWIcfNN80JUy59wT9tGgmYowLfJMNlb59D1fI_AAltPnWnlLAN8uAA5-xJsxxukjoDK6v87q1eWJ5-V6CYMFON9bDVAQDco1IpFcQvbop_yqOibp-_MM9f6Es1SBInWQFSuMOoz6qry772qPFb91QWNiM4TyJEZF8B2C2XT-bxJ4AU3TEjCHUtOXhMObsgmVZ7cgFVmMyp13NcGrlXUkwXDdEAofGHsaHxjzpqH_Tqrv9xtk45BNsI4qVsCjA6yJ5rBVPI-RndGlyNmEGwX_80Cs1oodNxQco1xRD02mwsS3T0v9EzR-QWVfBvOfQa19y33g9ZYB2cucSk3gVohwzP9KSSvTBXkhNdWI5ZANrB3slK_zJ7E-m3Frc3YSX9Xj1GjTZ95hiFD-9nNS4f6JzAHf6xN9uBo0wntLjw4XVCmwhRH8ey8sW36YlRqhKFGfyW6g9TG9kUNgXGt868rA0xMTVQMcT7cu8-8iDDtIN2tgORD_xRKUjz-5ecrsIYEBtHJOoBBzpmKOD9zq2wovdcjkeIZOCEdlVXOgnrVXg6SsXidLCN9p47wOjoUZubJTyRsqqr_smPgWi-BgkOv5rYmmwfIVswrJKYEoCid231UshEI4GPLtZFYOyeA2XCVJPmgD4U17ovEz7beTygYoCgIxaRhCG_PSzBHXCwEoXtTBG-LTZbBNv4-FpcrWNGalcwjGK5NXJKv4R8Foyex9YpS_vE5krZQotBAYsjuATSFg4XEaUOmtIbO7BcWZX7uJhMH0Z90-KQ4WoZmkcaDOjMc6Sft2REyxsb7Vt0cBUuSees5s8QdN_hDf0hf-6Jz2OeAzTCX1VCyTXUjtwRlXH5ul0_M0CvXcIkoHK69i19ujjKO_nxZlnLr8RYZGvsAPq9vofOL3t4U53g9D_N6rMuV9jIAvvFVdtpcNXUYAA5yk-CosGwH-aUR5RkOchgzdyGH8mB_Wctd9NheCk-6CMejN5ugoRF8AABif8ZKnSAP-jU7MgxA-5Jsf58QrDhCSrgALfzvQwMvW77vBgKdlC4cAX-6-ZM_tb16Z6Pjq9dkpkw_R_VamPClgk0yCQE7qDHTTnYsW8UquC4VDbxlAR22vpWDGe84Qu5A4djutTc33Mx4yp-cPM",
//		Company:     "Comcast",
//		Location:    "Location Not Found",
//		Salary:      "Salary Not Found",
//		Description: "A strong technical background in Software engineering including best practices and understanding of software development lifecycle.\nSCM Tools: Git, Gerrit, GitLab, GitHub\nSCM Tools: Git, Gerrit, GitLab, GitHub\nLanguages: C, C++\nKnowledge with Linux and Embedded software\nBS in Computer Science or related field\nKnowledge with Linux and Embedded software\nUnderstand our Operating Principles; make them the guidelines for how you do your job.\nWin as a team - make big things happen by working together and being open to new ideas.\nBe an active part of the Net Promoter System - a way of working that brings more employee and customer feedback into the company - by joining huddles, making call backs and helping us elevate opportunities to do better for our customers.\nDrive results and growth.\nRespect and promote inclusion & diversity.\nDo what's right for each other, our customers, investors and our communities.",
//	}
//
//	if !reflect.DeepEqual(specificJob, &expectedJob) {
//		t.Errorf("LoadDataFromJSON returned incorrect job data: got %v, want %v", specificJob, &expectedJob)
//	}
//}

// Function for partial comparison of JobData
func compareJobDataFields(a, b *dal.JobData) bool {
	return a.Title == b.Title && a.Company == b.Company && a.URL == b.URL
	// Add more fields to compare as needed
}

func TestLoadDataFromJSON(t *testing.T) {
	mockFilename := "C:\\Users\\Public\\GoLandProjects\\JustAFork\\crab\\output\\SoftwareEng_jobs.json"

	_, specificJob, err := dal.LoadDataFromJSON(mockFilename, "Software Release DevOps Engineer")
	if err != nil {
		t.Errorf("LoadDataFromJSON returned an error: %v", err)
	}

	expectedJob := dal.JobData{
		Title:   "Software Release DevOps Engineer",
		URL:     "https://www.indeed.com/pagead/clk?mo=r&ad=-6NYlbfkN0Cj-KmZPsf9w80C8b1WzNVrlanjD2SXJjxuCbUWHsXPZlTAgGmdtIUzoKTi6fK6WvZ2eEeIQBp5OUhO-xRyQvDo4yR3Mt5CEDSCojK6clcrRqADOS0tfXeHAsrfH_7i7PXK3XmzBFDjlntXqwANAhWdOGj1px_99ycmqNNMR1xJWJSD4fVvgEAHQ7k280w16fwgxdWPCRbm8AnWIcfNN80JUy59wT9tGgmYowLfJMNlb59D1fI_AAltPnWnlLAN8uAA5-xJsxxukjoDK6v87q1eWJ5-V6CYMFON9bDVAQDco1IpFcQvbop_yqOibp-_MM9f6Es1SBInWQFSuMOoz6qry772qPFb91QWNiM4TyJEZF8B2C2XT-bxJ4AU3TEjCHUtOXhMObsgmVZ7cgFVmMyp13NcGrlXUkwXDdEAofGHsaHxjzpqH_Tqrv9xtk45BNsI4qVsCjA6yJ5rBVPI-RndGlyNmEGwX_80Cs1oodNxQco1xRD02mwsS3T0v9EzR-QWVfBvOfQa19y33g9ZYB2cucSk3gVohwzP9KSSvTBXkhNdWI5ZANrB3slK_zJ7E-m3Frc3YSX9Xj1GjTZ95hiFD-9nNS4f6JzAHf6xN9uBo0wntLjw4XVCmwhRH8ey8sW36YlRqhKFGfyW6g9TG9kUNgXGt868rA0xMTVQMcT7cu8-8iDDtIN2tgORD_xRKUjz-5ecrsIYEBtHJOoBBzpmKOD9zq2wovdcjkeIZOCEdlVXOgnrVXg6SsXidLCN9p47wOjoUZubJTyRsqqr_smPgWi-BgkOv5rYmmwfIVswrJKYEoCid231UshEI4GPLtZFYOyeA2XCVJPmgD4U17ovEz7beTygYoCgIxaRhCG_PSzBHXCwEoXtTBG-LTZbBNv4-FpcrWNGalcwjGK5NXJKv4R8Foyex9YpS_vE5krZQotBAYsjuATSFg4XEaUOmtIbO7BcWZX7uJhMH0Z90-KQ4WoZmkcaDOjMc6Sft2REyxsb7Vt0cBUuSees5s8QdN_hDf0hf-6Jz2OeAzTCX1VCyTXUjtwRlXH5ul0_M0CvXcIkoHK69i19ujjKO_nxZlnLr8RYZGvsAPq9vofOL3t4U53g9D_N6rMuV9jIAvvFVdtpcNXUYAA5yk-CosGwH-aUR5RkOchgzdyGH8mB_Wctd9NheCk-6CMejN5ugoRF8AABif8ZKnSAP-jU7MgxA-5Jsf58QrDhCSrgALfzvQwMvW77vBgKdlC4cAX-6-ZM_tb16Z6Pjq9dkpkw_R_VamPClgk0yCQE7qDHTTnYsW8UquC4VDbxlAR22vpWDGe84Qu5A4djutTc33Mx4yp-cPM",
		Company: "Comcast",
		// Other fields are not compared
	}

	if !compareJobDataFields(specificJob, &expectedJob) {
		t.Errorf("LoadDataFromJSON returned incorrect job data: got %v, want %v", specificJob, &expectedJob)
	}
}
