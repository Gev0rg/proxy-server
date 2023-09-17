package storage

import (
	"context"
	"encoding/json"
	"io/ioutil"
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
	clientOptions := options.Client().ApplyURI("mongodb://localhost:27017")

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

	collection := s.client.Database("test").Collection("requests")

	_, err := collection.InsertOne(context.TODO(), parseReq)
	if err != nil {
		log.Fatal(err)
	}
}

func (s *Storage) SaveResponse(res *http.Response) {
	var parseResp models.Response

	parseResp.Code = res.StatusCode

	parseResp.Message = res.Status

	parseResp.Headers = make(map[string][]string)
	for key, value := range res.Header {
		parseResp.Headers[key] = value
	}

	if res.Body != nil {
		bodyBytes, _ := ioutil.ReadAll(res.Body)
		parseResp.Body = string(bodyBytes)
	}

	collection := s.client.Database("test").Collection("responses")

	_, err := collection.InsertOne(context.TODO(), parseResp)
	if err != nil {
		log.Fatal(err)
	}
}
