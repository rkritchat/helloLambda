package user

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"net/http"
	"os"
)

type Service interface {
	GetUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	CreateUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	JustReturnErr(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error)
	MethodNotAllowed() (*events.APIGatewayProxyResponse, error)
}

type service struct {
	Users map[string]*CustomerDetail
}

func NewService() Service {
	return &service{
		Users: make(map[string]*CustomerDetail),
	}
}

type CustomerDetail struct {
	Email     string `json:"email"`
	Firstname string `json:"firstname"`
	Lastname  string `json:"lastname"`
	Age       int    `json:"age"`
	Password  string `json:"password"`
}

func (s *service) GetUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("GET: %v", req)
	email := req.QueryStringParameters["email"]
	fmt.Printf("request email is %v", email)
	if len(email) == 0 {
		return toJson(http.StatusBadRequest, Response{ErrMsg: errors.New("email is required").Error()})
	}

	custInfo := s.Users[email]
	if custInfo == nil {
		return toJson(http.StatusBadRequest, Response{ErrMsg: errors.New("email is not found").Error()})
	}

	return toJson(http.StatusOK, custInfo)
}

type Response struct {
	Message string `json:"message,omitempty"`
	ErrMsg  string `json:"errMsg,omitempty"`
}

func (s *service) CreateUser(req events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Printf("POST: %v", req)
	request, err := validateCustDetail(req)
	if err != nil {
		return toJson(http.StatusBadRequest, Response{ErrMsg: err.Error()})
	}

	isExist := s.Users[request.Email] != nil
	if isExist {
		return toJson(http.StatusBadRequest, Response{ErrMsg: errors.New(fmt.Sprintf("email %v is already exists", request.Email)).Error()})
	}

	//setup password by reading from env
	request.Password = os.Getenv("ENV_AUTO_GEN_PWD")

	//store to memory
	s.Users[request.Email] = request
	return toJson(http.StatusOK, Response{Message: "Create user successfully"})
}

func validateCustDetail(req events.APIGatewayProxyRequest) (*CustomerDetail, error) {
	var request CustomerDetail
	err := json.Unmarshal([]byte(req.Body), &request)
	if err != nil {
		return nil, errors.New("request is invalid json format")
	}

	if len(request.Email) == 0 {
		return nil, errors.New("email is required")
	}

	if len(request.Firstname) == 0 {
		return nil, errors.New("firstname is required")
	}

	if len(request.Lastname) == 0 {
		return nil, errors.New("lastname is required")
	}

	if request.Age == 0 {
		return nil, errors.New("age is required")
	}

	return &request, nil
}

func toJson(httpStatus int, b interface{}) (*events.APIGatewayProxyResponse, error) {
	r, _ := json.Marshal(&b)
	return &events.APIGatewayProxyResponse{
		Headers:    map[string]string{"content-type": "application/json"},
		Body:       string(r),
		StatusCode: httpStatus,
	}, nil
}

func (s service) JustReturnErr(_ events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	fmt.Println("coming with method PUT, then just return err")
	return nil, errors.New("do nothing just Error")
}
func (s *service) MethodNotAllowed() (*events.APIGatewayProxyResponse, error) {
	return toJson(http.StatusBadRequest, errors.New("request http method is not allowed"))
}
