package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/go-resty/resty/v2"
)

func main() {
	tokenPtr := flag.String("token", "", "The API token to use")
	hostPtr := flag.String("host", "exed.canvas.harvard.edu", "The Canvas host to connect to")
	filePtr := flag.String("file", "", "File to read - must contain a list of email addresses, one per line")
	accountIdPtr := flag.Int("account_id", 139, "The Canvas account ID to use")
	courseIdPtr := flag.Int("course_id", 0, "The Canvas course ID to use")
	flag.Parse()

	// make sure we have a token
	if *tokenPtr == "" {
		log.Fatal("Must provide a Canvas API token")
	}
	client := resty.New()
	client.SetAuthToken(*tokenPtr)
	client.SetHostURL("https://" + *hostPtr)

	log.Println("*****************************")
	log.Println("* Zero-L Enrollment Cleaner *")
	log.Println("*****************************")
	// Get the Canvas account
	account, err := getAccount(*accountIdPtr, client)
	log.Printf("Using account: %s\n", account.Name)

	// Get the Canvas course
	course, err := getCourse(*courseIdPtr, client)
	log.Printf("Using course: %s\n", course.Name)

	file, err := os.Open(*filePtr)

	if err != nil {
		log.Fatal(err)
	}

	// Find a user for each email address in the file
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		email := scanner.Text()
		user, err := findUser(email, client)
		if err != nil {
			log.Printf("Email not found: %s\n", email)
			continue
		}
		log.Printf("Email Found: %s Name: %s ID: %d \n", email, user.Name, user.ID)

		// Delete user's enrollments in the course
		err = deleteEnrollments(course.ID, user.ID, client)
		if err != nil {
			log.Printf("Error deleting enrollments: %s\n", err)
			continue
		}
	}
}

func findUser(email string, client *resty.Client) (*User, error) {
	url := "/api/v1/accounts/1/users?include[]=email&search_term=" + email
	resp, err := client.R().Get(url)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%s: %s", email, resp.String())
	}
	var users []User
	err = json.Unmarshal(resp.Body(), &users)
	if err != nil {
		return nil, err
	}
	if len(users) == 0 {
		return nil, fmt.Errorf("%s: not found", email)
	}

	for _, user := range users {
		if user.Email == email {
			return &user, nil
		}
	}

	return nil, fmt.Errorf("%s: not found", email)

}

// User represents a Canvas user
type User struct {
	ID    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

// Get a Canvas account
func getAccount(id int, client *resty.Client) (*Account, error) {
	resp, err := client.R().Get("/api/v1/accounts/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%d: %s", id, resp.String())
	}
	var account Account
	err = json.Unmarshal(resp.Body(), &account)
	if err != nil {
		return nil, err
	}
	return &account, nil
}

type Account struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Get a Canvas course
func getCourse(id int, client *resty.Client) (*Course, error) {
	resp, err := client.R().Get("/api/v1/courses/" + strconv.Itoa(id))
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != 200 {
		return nil, fmt.Errorf("%d: %s", id, resp.String())
	}
	var course Course
	err = json.Unmarshal(resp.Body(), &course)
	if err != nil {
		return nil, err
	}
	return &course, nil
}

type Course struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Delete a user's enrollments in a course
func deleteEnrollments(courseId int, userId int, client *resty.Client) error {
	url := "/api/v1/courses/" + strconv.Itoa(courseId) + "/enrollments?user_id=" + strconv.Itoa(userId)
	resp, err := client.R().Get(url)
	if err != nil {
		return err
	}
	if resp.StatusCode() != 200 {
		return fmt.Errorf("%d: %s", userId, resp.String())
	}
	var enrollments []Enrollment
	err = json.Unmarshal(resp.Body(), &enrollments)
	if err != nil {
		return err
	}
	if len(enrollments) == 0 {
		log.Printf("No enrollments found for user %d in course %d\n", userId, courseId)
		return nil
	}

	for _, enrollment := range enrollments {
		log.Printf("Deleting enrollment %d for user %d in course %d\n", enrollment.ID, userId, courseId)
		delete_url := "/api/v1/courses/" + strconv.Itoa(courseId) + "/enrollments/" + strconv.Itoa(enrollment.ID)
		resp, err := client.R().Delete(delete_url)
		if err != nil {
			return err
		}
		if resp.StatusCode() != 200 {
			return fmt.Errorf("course=%d enrollment=%d: %s", courseId, enrollment.ID, resp.String())
		}
	}
	return nil
}

type Enrollment struct {
	ID   int    `json:"id"`
	Type string `json:"type"`
	Role string `json:"role"`
}
