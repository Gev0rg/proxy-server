package storage

import (
	"context"
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/Gev0rg/proxy-server/models"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type Storage struct {
	client *mongo.Client
}

func (s *Storage) Connect() {
	clientOptions := options.Client().ApplyURI("mongodb://root:root@mongo:27017")

	client, err := mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}

	s.client = client
}

func (s *Storage) GetClient() *mongo.Client {
	return s.client
}

func (s *Storage) SaveRequest(req *http.Request) {
	var parseReq models.Request

	parseReq.HttpMethod = req.Method

	if req.URL.Scheme == "" {
		parseReq.Path = "https://" + req.URL.Host + req.URL.Path
	} else {
		parseReq.Path = req.URL.Scheme + "://" + req.URL.Host + req.URL.Path
	}

	parseReq.Headers = make(map[string][]string)
	for key, value := range req.Header {
		parseReq.Headers[key] = value
	}

	parseReq.Cookie = make(map[string]string)
	for _, cookie := range req.Cookies() {
		parseReq.Cookie[cookie.Name] = cookie.Value
	}

	parseReq.GetParams = make(map[string][]string)
	if req.Method == "GET" {
		for key, value := range req.URL.Query() {
			parseReq.GetParams[key] = value
		}
	}

	byt := []byte{}
	read, _ := req.Body.Read(byt)

	var data map[string]string

	if read != 0 {
		if err := json.Unmarshal(byt, &data); err != nil {
			log.Fatal(err)
			return
		}
	}

	parseReq.PostParams = make(map[string]string)
	for key, value := range data {
		parseReq.PostParams[key] = value
	}

	collection := s.client.Database("admin").Collection("requests")

	_, err := collection.InsertOne(context.TODO(), parseReq)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Storage) SaveResponse(res *http.Response) {
	parsedResp := models.Response{
		Code:    res.StatusCode,
		Message: res.Status,
		Headers: make(map[string][]string),
	}
	
	for key, value := range res.Header {
		parsedResp.Headers[key] = value
	}

	if res.Request.Method == http.MethodPost && res.Body != nil {
		body, err := io.ReadAll(res.Body)
		if err != nil {
			panic(err)
		}

		parsedResp.Body = string(body)
	}

	collection := s.client.Database("admin").Collection("responses")

	_, err := collection.InsertOne(context.TODO(), parsedResp)
	if err != nil {
		log.Fatal(err)
	}
}
