package main

import (
	"context"
	"lab01/weather/weatherpb"
	"log"

	"google.golang.org/grpc"
)

func main() {
	cc, err := grpc.Dial("localhost:8888", grpc.WithInsecure())

	if err != nil {
		log.Fatalf("err while dial %v\n", err)
	}

	defer cc.Close()

	client := weatherpb.NewWeatherServiceClient(cc)
	getTemp(client)
}

func getTemp(w weatherpb.WeatherServiceClient) {
	log.Println("Get temp calling...")
	resp, err := w.GetTemp(context.Background(), &weatherpb.GetTempRequest{
		Lat: 7.857940,
		Lon: 51.956055,
	})

	if err != nil {
		log.Fatalf("call err %v\n", err)
	}

	log.Printf("Temperate %v\n", resp.Temp)
}
