package ranger

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"
	"time"

	"gopkg.in/yaml.v3"
)

const (
	getUsersEndpoint           = `/service/xusers/users`
	getServicePoliciesEndpoint = `/service/public/v2/api/service/%s/policy`
	getGroupsEndpoint          = `/service/xusers/groups`
)

type User struct {
	ID           int64    `json:"id,omitempty"`
	Name         string   `json:"name,omitempty"`
	FirstName    string   `json:"firstName,omitempty"`
	LastName     string   `json:"lastName,omitempty"`
	EmailAddress string   `json:"emailAddress,omitempty"`
	UserRoleList []string `json:"userRoleList,omitempty"`
	Password     string   `json:"password,omitempty"`
	SyncSource   string   `json:"syncSource,omitempty"`
	GroupIdList  []int64  `json:"groupIdList,omitempty"`
}

type Group struct {
	ID          int64  `json:"id,omitempty"`
	Name        string `json:"name,omitempty"`
	Description string `json:"description,omitempty"`
	SyncSource  string `json:"syncSource,omitempty"`
}

type getResponse struct {
	PageSize   int      `json:"pageSize"`
	StartIndex int      `json:"startIndex"`
	ResultSize int      `json:"resultSize"`
	VXUsers    []*User  `json:"vXUsers,omitempty"`
	VXGroups   []*Group `json:"vXGroups,omitempty"`
}

//go:generate go run github.com/vektra/mockery/v2@v2.53.4 --name=Client --output=./mocks --outpkg=mocks

type ClientWrapper struct {
	Client Client
}

func (aw *ClientWrapper) UnmarshalYAML(value *yaml.Node) error {
	var cl client
	if err := value.Decode(&cl); err != nil {
		return err
	}
	aw.Client = &cl
	cl.client = &http.Client{}
	return nil
}

func (cw *ClientWrapper) GetUsers() (map[string]*User, error) {
	return cw.Client.GetUsers()
}

func (cw *ClientWrapper) GetGroups() (map[string]*Group, error) {
	return cw.Client.GetGroups()
}

func (cw *ClientWrapper) GetPolicies(serviceName string) ([]*Policy, error) {
	return cw.Client.GetPolicies(serviceName)
}

type Client interface {
	GetUsers() (map[string]*User, error)
	GetGroups() (map[string]*Group, error)
	GetPolicies(serviceName string) ([]*Policy, error)
}

type client struct {
	URL      string `yaml:"url" json:"url" omitempty"`
	Username string `yaml:"username" json:"username" omitempty"`
	Password string `yaml:"password" json:"password" omitempty"`
	client   *http.Client
}

func NewClient(url, username, password string) Client {
	return &client{
		URL:      url,
		Username: username,
		Password: password,
		client:   &http.Client{},
	}
}
func (c *client) GetUsers() (map[string]*User, error) {

	responses, err := c.executeBatchRequest(http.MethodGet, getUsersEndpoint)
	if err != nil {
		return nil, err
	}

	usersMap := make(map[string]*User)

	// Process all response batches into the map
	for _, resp := range responses {
		// Use type assertion to get the correct response type
		for _, user := range resp.VXUsers {
			usersMap[user.Name] = user
		}

	}

	log.Printf("Number of Ranger Users pulled: %d\n", len(usersMap))
	return usersMap, nil

}

func (c *client) GetGroups() (map[string]*Group, error) {

	responses, err := c.executeBatchRequest(http.MethodGet, getGroupsEndpoint)
	if err != nil {
		return nil, err
	}

	groupsMap := make(map[string]*Group)

	for _, resp := range responses {
		for _, group := range resp.VXGroups {
			groupsMap[group.Name] = group
		}

	}

	log.Printf("Number of Ranger Groups pulled: %d\n", len(groupsMap))
	return groupsMap, nil

}

func (c *client) GetPolicies(serviceName string) ([]*Policy, error) {
	var policies []*Policy
	err := c.executeRequest(http.MethodGet, fmt.Sprintf(getServicePoliciesEndpoint, serviceName), &policies, nil)
	return policies, err
}

func (c *client) createRequest(method, endpoint string, reqBody interface{}) (*http.Request, error) {

	// Ensure client exists
	if c.client == nil {
		c.client = &http.Client{}
	}

	var jsonBody []byte
	var err error

	// Marshal body if POST request
	if reqBody != nil {
		jsonBody, err = json.Marshal(reqBody)
		if err != nil {
			return nil, err
		}
	}

	// Create Request
	req, err := http.NewRequest(method, fmt.Sprintf("%s%s", c.URL, endpoint), bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, err
	}

	// Add auth headers
	req.SetBasicAuth(c.Username, c.Password)
	req.Header.Set("Content-Type", "application/json")

	return req, nil

}

func (r *client) executeRequest(method string, endpoint string, v interface{}, reqBody interface{}) error {

	req, err := r.createRequest(method, endpoint, reqBody)
	if err != nil {
		return err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		bodyString := string(bodyBytes)

		if strings.Contains(bodyString, "INVALID_INPUT_DATA") {
			return nil
		}

		return fmt.Errorf("request to %s failed with status %s\n%s", req.URL.String(), resp.Status, bodyString)
	}

	if v != nil {
		return json.NewDecoder(resp.Body).Decode(v)
	}

	return nil

}

// executeBatchRequest performs paginated API requests and returns all aggregated results
func (r *client) executeBatchRequest(method string, endpoint string) ([]getResponse, error) {

	results := make([]getResponse, 500)
	pageSize := 200
	startIndex := 0

	for {

		batchEndpoint := fmt.Sprintf("%s?pageSize=%d&startIndex=%d", endpoint, pageSize, startIndex)

		// Marshall into generic get
		getResponse := &getResponse{}
		if err := r.executeRequest(method, batchEndpoint, getResponse, nil); err != nil {
			return nil, err
		}

		// Add this batch's response to our results
		results = append(results, *getResponse)

		fmt.Printf("%v - Pulled batch with %d items...\n", time.Now().Format("2006-01-02 15:04:05"), getResponse.ResultSize)

		if getResponse.ResultSize < int(pageSize) {
			break
		}

		startIndex += pageSize

	}

	return results, nil

}
