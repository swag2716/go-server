package main

import (
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"time"
)

const PhotoApi = "https://api.pexels.com/v1"
const VideoApi = "https://api.pexels.com/videos"

type Client struct {
	Token          string
	hc             http.Client
	RemainingTimes int32
}

func NewClient(token string) *Client {
	c := http.Client{}
	return &Client{Token: token, hc: c}
}

type SearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:"total_results"`
	NextPage     string  `json:"next_page"`
	Photos       []Photo `json:"photos"`
}

type Photo struct {
	Id              int32       `json:"id"`
	Width           int32       `json:"width"`
	Height          int32       `json:"height"`
	Url             string      `json:"url"`
	Photographer    string      `json:"photographer"`
	PhotographerUrl string      `json:"photographer_url"`
	Src             PhotoSource `json:"src"`
}

type curatedResult struct {
	Page     int32   `json:"page"`
	PerPage  int32   `json:"per_page"`
	NextPage int32   `json:"next_page"`
	Photos   []Photo `json:"photos"`
}

type PhotoSource struct {
	Original  string `json:"original"`
	Large     string `json:"large"`
	Large2x   string `json:"large2x"`
	Medium    string `json:"medium"`
	Small     string `json:"small"`
	Portrait  string `json:"portrait"`
	Square    string `json:"square"`
	Landscape string `json:"landscape"`
	Tiny      string `json:"tiny"`
}

func (c *Client) SearchPhotos(query string, perPage int32, page int32) (*SearchResult, error) {
	url := fmt.Sprintf(PhotoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	response, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var result SearchResult
	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) requestDoWithAuth(method, url string) (*http.Response, error) {
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return nil, err
	}

	req.Header.Add("Authorization", c.Token)
	response, err := c.hc.Do(req)

	if err != nil {
		return response, err
	}

	times, err := strconv.Atoi(response.Header.Get("X-Ratelimit-Remaining"))

	if err != nil {
		return response, nil
	} else {
		c.RemainingTimes = int32(times)
	}
	return response, nil
}

func (c *Client) GetPhoto(id int32) (*Photo, error) {
	url := fmt.Sprintf(PhotoApi+"/photos/%d", id)
	response, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var result Photo

	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) curatedPhotos(perPage, page int) (*curatedResult, error) {
	url := fmt.Sprintf(PhotoApi+"/curated?per_page%d&page%d", perPage, page)
	response, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var result curatedResult

	err = json.Unmarshal(data, &result)
	return &result, err
}

func (c *Client) GetRandomPhoto() (*Photo, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.curatedPhotos(1, randNum)

	if err == nil && len(result.Photos) == 1 {
		return &result.Photos[0], nil
	}

	return nil, err
}

type Video struct {
	Id            int32           `json:"id"`
	Width         int32           `json:"width"`
	Height        int32           `json:"height"`
	Url           string          `json:"url"`
	Image         string          `json:"image"`
	FullRes       interface{}     `json:"full_res"`
	Duration      float64         `json:"duration"`
	VideoFiles    []VideoFiles    `json:"video_files"`
	VideoPictures []VideoPictures `json:"video_pictures"`
}

type VideoPictures struct {
	Id      int32  `json:"id"`
	Picture string `json:"picutre"`
	Nr      int32  `json:"nr"`
}
type VideoSearchResult struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:""total_results`
	NextPage     string  `json:"next_page"`
	Videos       []Video `json:"videos"`
}
type PopularVideos struct {
	Page         int32   `json:"page"`
	PerPage      int32   `json:"per_page"`
	TotalResults int32   `json:""total_results`
	Url          string  `json:"url"`
	Videos       []Video `json:"videos"`
}
type VideoFiles struct {
	Id       int32  `json:"id"`
	Quality  string `json:"quality"`
	FileType string `json:"file_type"`
	Width    int32  `json:"width"`
	Height   int32  `json:"height"`
	Link     string `json:"link"`
}

func (c *Client) SearchVideo(query string, perPage int, page int) (*VideoSearchResult, error) {
	url := fmt.Sprintf(VideoApi+"/search?query=%s&per_page=%d&page=%d", query, perPage, page)
	response, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var result VideoSearchResult

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) PopularVideo(perPage int, page int) (*PopularVideos, error) {
	url := fmt.Sprintf(VideoApi+"/popular?per_page=%d&page=%d", perPage, page)
	response, err := c.requestDoWithAuth("GET", url)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	data, err := io.ReadAll(response.Body)

	if err != nil {
		return nil, err
	}

	var result PopularVideos

	err = json.Unmarshal(data, &result)

	return &result, err
}

func (c *Client) GetRemainingRequestInThisMonth() int32 {
	return c.RemainingTimes
}

func (c *Client) GetRandomVideo() (*Video, error) {
	rand.Seed(time.Now().Unix())
	randNum := rand.Intn(1001)
	result, err := c.PopularVideo(1, randNum)

	if err == nil && len(result.Videos) == 1 {
		return &result.Videos[0], nil
	}

	return nil, err
}

func main() {
	os.Setenv("PexelsToken", "vyLgRS3icz3zekH1UdM6PaBBwQ8GgYuPzPlW2UZ6X49vCOnoPVVtGdxx")
	TOKEN := os.Getenv("PexelsToken")
	c := NewClient(TOKEN)

	result, err := c.SearchPhotos("waves", 15, 1)

	if err != nil {
		fmt.Errorf("Search error : %v", err)
	}

	if result.Page == 0 {
		fmt.Errorf("Search result wrong")
	}

	fmt.Println(result)
}
