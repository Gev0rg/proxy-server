package api

import (
	"bufio"
	"context"
	"encoding/json"
	"github.com/Gev0rg/proxy-server/models"
	"github.com/Gev0rg/proxy-server/storage"
	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
)

type Handlers struct {
	Storage *storage.Storage
}

type Not404Response struct {
	Path       string `json:"path"`
	StatusCode int    `json:"status_code"`
}

func Respond(w http.ResponseWriter, code int, data interface{}) {
	w.WriteHeader(code)
	if data != nil {
		err := json.NewEncoder(w).Encode(data)
		if err != nil {
			log.Println(err)
			return
		}
	}
}

func (h *Handlers) GetRequests(w http.ResponseWriter, r *http.Request) {
	collection := h.Storage.GetClient().Database("admin").Collection("requests")

	var results []*models.ReplyRequest

	findOptions := options.Find()

	cur, err := collection.Find(context.TODO(), bson.D{{}}, findOptions)
	if err != nil {
		log.Fatal(err)
	}

	for cur.Next(context.TODO()) {
		var elem models.ReplyRequest
		err := cur.Decode(&elem)
		if err != nil {
			log.Fatal(err)
		}
		results = append(results, &elem)
	}

	Respond(w, 200, results)
}

func (h *Handlers) GetRequestByID(w http.ResponseWriter, r *http.Request) {
	collection := h.Storage.GetClient().Database("admin").Collection("requests")

	var result models.ReplyRequest

	params := mux.Vars(r)
	ID, _ := params["id"]

	objectId, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{{"_id", objectId}}

	_ = collection.FindOne(context.TODO(), filter).Decode(&result)

	Respond(w, 200, result)
}

func (h *Handlers) RepeatRequest(w http.ResponseWriter, r *http.Request) {

	collection := h.Storage.GetClient().Database("admin").Collection("requests")

	var result models.ReplyRequest

	params := mux.Vars(r)
	ID, _ := params["id"]

	objectId, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{{"_id", objectId}}

	_ = collection.FindOne(context.TODO(), filter).Decode(&result)

	var repeatReq *http.Request

	var url string
	url = result.Path
	if result.HttpMethod == "GET" {
		if len(result.GetParams) != 0 {
			url += "?"
		}
		for key, value := range result.GetParams {
			url += key + "=" + strings.Join(value, ",")
		}
	} else {
		for key, value := range result.PostParams {
			repeatReq.PostForm.Set(key, value)
		}
	}

	repeatReq, _ = http.NewRequest(result.HttpMethod, url, strings.NewReader(result.Body))

	for key, value := range result.Headers {
		repeatReq.Header.Set(key, strings.Join(value, ""))
	}

	for key, value := range result.Cookie {
		repeatReq.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}

	client := http.Client{}

	resp, _ := client.Do(repeatReq)

	w.WriteHeader(resp.StatusCode)

	copyHeader(w.Header(), resp.Header)

	io.Copy(w, resp.Body)
}

func (h *Handlers) DirSearch(w http.ResponseWriter, r *http.Request) {

	file, err := os.OpenFile("dirsearch", os.O_RDONLY, 0666)
	if err != nil {
		log.Fatal(err)
	}
	fileScanner := bufio.NewScanner(file)

	var responseAnswer []Not404Response

	collection := h.Storage.GetClient().Database("admin").Collection("requests")

	var result models.ReplyRequest

	params := mux.Vars(r)
	ID, _ := params["id"]

	objectId, _ := primitive.ObjectIDFromHex(ID)

	filter := bson.D{{"_id", objectId}}

	_ = collection.FindOne(context.TODO(), filter).Decode(&result)

	var repeatReq *http.Request

	var url string
	url = result.Path
	if result.HttpMethod == "GET" {
		if len(result.GetParams) != 0 {
			url += "?"
		}
		for key, value := range result.GetParams {
			url += key + "=" + strings.Join(value, ",")
		}
	} else {
		for key, value := range result.PostParams {
			repeatReq.PostForm.Set(key, value)
		}
	}

	repeatReq, _ = http.NewRequest(result.HttpMethod, url, strings.NewReader(result.Body))

	for key, value := range result.Headers {
		repeatReq.Header.Set(key, strings.Join(value, ""))
	}

	for key, value := range result.Cookie {
		repeatReq.AddCookie(&http.Cookie{
			Name:  key,
			Value: value,
		})
	}

	for fileScanner.Scan() {
		repeatReq.URL.Path += fileScanner.Text()
		pathLen := len(fileScanner.Text())

		client := http.Client{}

		resp, err := client.Do(repeatReq)
		if err != nil {
			log.Println(err)
			repeatReq.URL.Path = repeatReq.URL.Path[:len(repeatReq.URL.Path)-pathLen]
			continue
		}

		log.Println(resp.StatusCode, fileScanner.Text())
		if resp.StatusCode != 404 {
			responseAnswer = append(responseAnswer, Not404Response{
				StatusCode: resp.StatusCode,
				Path:       repeatReq.URL.Path,
			})
		}
		repeatReq.URL.Path = repeatReq.URL.Path[:len(repeatReq.URL.Path)-pathLen]
	}

	file.Close()
	Respond(w, 200, responseAnswer)
}

func copyHeader(dst, src http.Header) {
	for k, vv := range src {
		for _, v := range vv {
			dst.Add(k, v)
		}
	}
}
