package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/urfave/negroni"
)

// Generated by https://quicktype.io

type MyCompaniesResponse struct {
	Total  int     `json:"_total"`
	Values []Value `json:"values"`
}

type Value struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

// Generated by https://quicktype.io

type MyCompaniesUpdatesResponse struct {
	Count  int64                             `json:"_count"`
	Start  int64                             `json:"_start"`
	Total  int64                             `json:"_total"`
	Values []MyCompaniesUpdatesResponseValue `json:"values"`
}

type MyCompaniesUpdatesResponseValue struct {
	IsCommentable  bool           `json:"isCommentable"`
	IsLikable      bool           `json:"isLikable"`
	IsLiked        bool           `json:"isLiked"`
	NumLikes       int64          `json:"numLikes"`
	Timestamp      int64          `json:"timestamp"`
	UpdateComments UpdateComments `json:"updateComments"`
	UpdateContent  UpdateContent  `json:"updateContent"`
	UpdateKey      string         `json:"updateKey"`
	UpdateType     UpdateType     `json:"updateType"`
	Likes          *Likes         `json:"likes,omitempty"`
}

type Likes struct {
	Total  int64        `json:"_total"`
	Values []LikesValue `json:"values"`
}

type LikesValue struct {
	Person Person `json:"person"`
}

type Person struct {
	FirstName string `json:"firstName"`
	ID        string `json:"id"`
	LastName  string `json:"lastName"`
}

type UpdateComments struct {
	Total int64 `json:"_total"`
}

type UpdateContent struct {
	Company             Company             `json:"company"`
	CompanyStatusUpdate CompanyStatusUpdate `json:"companyStatusUpdate"`
}

type Company struct {
	ID   int64       `json:"id"`
	Name CompanyName `json:"name"`
}

type CompanyStatusUpdate struct {
	Share Share `json:"share"`
}

type Share struct {
	Comment    string     `json:"comment"`
	Content    Content    `json:"content"`
	ID         string     `json:"id"`
	Source     Source     `json:"source"`
	Timestamp  int64      `json:"timestamp"`
	Visibility Visibility `json:"visibility"`
}

type Content struct {
	Description       string `json:"description"`
	EyebrowURL        string `json:"eyebrowUrl"`
	ShortenedURL      string `json:"shortenedUrl"`
	SubmittedImageURL string `json:"submittedImageUrl"`
	SubmittedURL      string `json:"submittedUrl"`
	ThumbnailURL      string `json:"thumbnailUrl"`
	Title             string `json:"title"`
}

type Source struct {
	Application            Application `json:"application"`
	ServiceProvider        Application `json:"serviceProvider"`
	ServiceProviderShareID string      `json:"serviceProviderShareId"`
}

type Application struct {
	Name ApplicationName `json:"name"`
}

type Visibility struct {
	Code Code `json:"code"`
}

type CompanyName string

const (
	Websays CompanyName = "Websays"
)

type ApplicationName string

const (
	Linkedin ApplicationName = "LINKEDIN"
	Postcron ApplicationName = "Postcron"
)

type Code string

const (
	Anyone Code = "anyone"
)

type UpdateType string

const (
	Cmpy UpdateType = "CMPY"
)

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/linkedin/getmycompanies", GetMyCompanies).Methods("GET")

	ngroni := negroni.Classic()

	ngroni.UseHandler(handlers.CORS(
		handlers.AllowedHeaders([]string{
			"Origin", "Content-Type", "Accept", "Authorization", "Accept-Charset", "Bearer"}),
		handlers.AllowedOrigins([]string{"*"}),
		handlers.AllowedMethods([]string{"GET", "HEAD", "POST", "PUT", "OPTIONS"}),
		handlers.AllowCredentials(),
		handlers.MaxAge(1209600))(router))

	log.Fatal(http.ListenAndServe(":8081", ngroni))
}

func GetMyCompanies(w http.ResponseWriter, r *http.Request) {
	fmt.Println("GetMyCompanies")
	accessToken := r.Header.Get("Bearer")
	fmt.Println(accessToken)
	client := &http.Client{
		CheckRedirect: nil,
	}

	if accessToken != "" {
		req, _ := http.NewRequest("GET", "https://api.linkedin.com/v1/companies?format=json&is-company-admin=true", nil)
		req.Header.Add("oauth_token", accessToken)
		resp, err := client.Do(req)
		if err != nil {
			fmt.Println("GetMyCompanies ERROR")
			fmt.Println(err)
		} else {
			fmt.Println("GetMyCompanies RESP")
			defer resp.Body.Close()
			if resp.StatusCode == http.StatusOK {
				bodyBytes, _ := ioutil.ReadAll(resp.Body)
				bodyString := string(bodyBytes)
				fmt.Println(bodyString)
				// json.NewEncoder(w).Encode(bodyString)

				myCResponse := MyCompaniesResponse{}
				json.Unmarshal(bodyBytes, &myCResponse)

				GetCompaniesUpdates(myCResponse.Values[0].ID, accessToken)
			} else {
				json.NewEncoder(w).Encode(resp)
			}

		}
	} else {
		fmt.Println("boom")
	}

	client = nil
}

func GetCompaniesUpdates(companiesId int, accessToken string) {

	client := &http.Client{
		CheckRedirect: nil,
	}

	url := fmt.Sprintf("https://api.linkedin.com/v2/organizations/%s", strconv.Itoa(companiesId))
	fmt.Println(url)
	// url := fmt.Sprintf("https://api.linkedin.com/v1/companies/%s/updates?format=json", strconv.Itoa(companiesId))

	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("oauth_token", accessToken)
	fmt.Println(accessToken)
	resp, err := client.Do(req)

	defer resp.Body.Close()
	if err != nil {
		fmt.Println("GetCompaniesUpdates ERROR")
		fmt.Println(err)
	} else {
		fmt.Println("GetCompaniesUpdates RESP")

		if resp.StatusCode == http.StatusOK {
			bodyBytes, _ := ioutil.ReadAll(resp.Body)
			myCUResponse := MyCompaniesUpdatesResponse{}
			json.Unmarshal(bodyBytes, &myCUResponse)
		} else {
			fmt.Println("GetCompaniesUpdates ERROR")
			fmt.Println(resp)
		}
	}

}

// func HandleAuthCodeResponse(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("handleAuthCodeResponse")
// 	fmt.Println("code")
// }

// func GetLnAccessToken(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("GetLnAccesToken")
// 	keys, _ := r.URL.Query()["code"]

// 	authCode := keys[0]
// 	fmt.Println(authCode)
// 	redirectURI := "http://localhost:8081/linkedin/getaccesstokenresponse"

// 	// u, _ := url.ParseRequestURI(redirectURI)

// 	clientID := "86gkduq4srneks"
// 	secretID := "5ZLK7n8nid9nZDvT"
// 	client := &http.Client{
// 		CheckRedirect: nil,
// 	}

// 	if authCode != "" {
// 		var Url *url.URL
// 		Url, err := url.Parse("https://www.linkedin.com")
// 		// url := fmt.Sprintf("https://www.linkedin.com/oauth/v2/accessToken?&client_id=%s&client_secret=%s", accessToken, redirURL, clientID, secretID)
// 		// tokenURL := "https://www.linkedin.com/oauth/v2/accessToken"
// 		// u, _ := url.Parse(tokenURL)
// 		Url.Path += "/oauth/v2/accessToken"
// 		params := url.Values{}
// 		fmt.Println("------------------------")
// 		fmt.Println("------------------------")
// 		params.Add("grant_type", "authorization_code")
// 		params.Add("code", authCode)
// 		params.Add("redirect_uri", redirectURI)
// 		params.Add("client_id", clientID)
// 		params.Add("client_secret", secretID)
// 		params.Add("state", "987654321")
// 		Url.RawQuery = params.Encode()
// 		req, _ := http.NewRequest("POST", Url.String(), nil)

// 		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
// 		req.Header.Add("Host", "www.linkedin.com")

// 		fmt.Println("++++++++++++++++++++++++++++++++")
// 		fmt.Println(req.URL)
// 		fmt.Println("++++++++++++++++++++++++++++++++")
// 		resp, err := client.Do(req)
// 		if err != nil {
// 			fmt.Println(" ********************** ERROR ********************** ")
// 			fmt.Println(err)
// 		} else {
// 			fmt.Println(" ********************** RESP ********************** ")
// 			fmt.Println(resp)
// 			// json.NewEncoder(w).Encode(resp)
// 		}
// 	}
// }

// func HandleGetLnAccessTokenResponse(w http.ResponseWriter, r *http.Request) {
// 	fmt.Println("HandleGetLnAccessTokenResponse")
// 	decoder := json.NewDecoder(r.Body)
// 	var t string
// 	err := decoder.Decode(&t)
// 	if err != nil {
// 		panic(err)
// 	}
// 	log.Println(t)
// 	json.NewEncoder(w).Encode("Todo bien")
// }

// volumes:
// - /home/fabio/Documentos/proyectos/dashboard-api/application/api-0.1/classes/:/var/www/html/application/api-0.1/classes/
// - /home/fabio/Documentos/proyectos/dashboard-api/application/api-1.0/classes/:/var/www/html/application/api-1.0/classes/
// - /home/fabio/Documentos/proyectos/dashboard-api/application/api-0.1/config/:/var/www/html/application/api-0.1/config/
// - /home/fabio/Documentos/proyectos/dashboard-api/application/api-1.0/config/:/var/www/html/application/api-0.1/config/
// - /home/fabio/Documentos/proyectos/dashboard-api/modules/:/var/www/html/modules/
// adminer:
