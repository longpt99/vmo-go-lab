package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"lab01/weather/weatherpb"
	"log"
	"net"
	"net/http"

	"google.golang.org/grpc"
)

type server struct {
	weatherpb.WeatherServiceServer
}

func main() {
	lis, err := net.Listen("tcp", "0.0.0.0:8888")

	if err != nil {
		log.Fatalf("err while create listen %v\n", err)
	}

	s := grpc.NewServer()

	weatherpb.RegisterWeatherServiceServer(s, &server{})

	fmt.Println("Server connecting...")

	err = s.Serve(lis)

	if err != nil {
		log.Fatalf("err while serve %v", err)
	}
}

type WeatherResponse struct {
	Coord struct {
		Lon float64 `json:"lon"`
		Lat float64 `json:"lat"`
	} `json:"coord"`
	Weather []struct {
		ID          int    `json:"id"`
		Main        string `json:"main"`
		Description string `json:"description"`
		Icon        string `json:"icon"`
	} `json:"weather"`
	Base string `json:"base"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		TempMin   float64 `json:"temp_min"`
		TempMax   float64 `json:"temp_max"`
		Pressure  int     `json:"pressure"`
		Humidity  int     `json:"humidity"`
		SeaLevel  int     `json:"sea_level"`
		GrndLevel int     `json:"grnd_level"`
	} `json:"main"`
	Visibility int `json:"visibility"`
	Wind       struct {
		Speed float64 `json:"speed"`
		Deg   int     `json:"deg"`
		Gust  float64 `json:"gust"`
	} `json:"wind"`
	Clouds struct {
		All int `json:"all"`
	} `json:"clouds"`
	Dt  int `json:"dt"`
	Sys struct {
		Sunrise int `json:"sunrise"`
		Sunset  int `json:"sunset"`
	} `json:"sys"`
	Timezone int    `json:"timezone"`
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Cod      int    `json:"cod"`
}

func (*server) GetTemp(ctx context.Context, req *weatherpb.GetTempRequest) (*weatherpb.GetTempResponse, error) {
	log.Println("Temp is coming..")
	response, err := http.Get("https://api.openweathermap.org/data/2.5/weather?lat=7.857940&lon=51.956055&units=metric&appid=81804a170a781fff7807ddc3eb8fb016")

	if err != nil {
		log.Fatalf("Get temp err %v", err)
	}

	defer response.Body.Close()
	body, err := io.ReadAll(response.Body)

	if err != nil {
		log.Fatalf("Get body err %v", err)
	}

	var weather WeatherResponse
	json.Unmarshal(body, &weather)

	resp := &weatherpb.GetTempResponse{
		Temp: weather.Main.Temp,
	}

	return resp, nil
}
