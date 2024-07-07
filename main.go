package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"time"
)

type Weather struct {
	Current struct {
		TempC    float64 `json:"temperature_2m"`
		Humidity float64 `json:"relative_humidity_2m"`
	} `json:"current"`
	Hourly struct {
		Time          []string  `json:"time"`
		Temperature2m []float64 `json:"temperature_2m"`
		// Agrega otros campos relevantes aquí, como la humedad horaria si es necesario
	} `json:"hourly"`
}

func main() {
	res, err := http.Get("https://api.open-meteo.com/v1/forecast?latitude=-32.9468&longitude=-60.6393&current=temperature_2m,relative_humidity_2m,is_day,precipitation,rain&hourly=temperature_2m")

	if err != nil {
		panic(err)
	}
	defer res.Body.Close()

	if res.StatusCode != 200 {
		panic("Weather API no disponible")
	}

	body, err := io.ReadAll(res.Body)

	if err != nil {
		panic(err)
	}

	var weather Weather
	err = json.Unmarshal(body, &weather)
	if err != nil {
		panic(err)
	}

	// Verificar si se proporcionó un argumento en la línea de comandos
	if len(os.Args) > 1 {

		if os.Args[1] == "1" {
			// Mostrar los horarios de los siguientes tres días
			fmt.Println("Pronóstico para los siguientes tres días:")
			now := time.Now()
			for i, time := range weather.Hourly.Time {
				t, err := timeParse(time)
				if err != nil {
					continue
				}
				if t.After(now) && t.Before(now.Add(72*time.Hour)) {
					fmt.Printf("⏰ %s: 🌡️ %.1f°C\n", t.Format("2006-01-02 15:04"), weather.Hourly.Temperature2m[i])
				}
			}
		} else if os.Args[1] == "2" {
			// Mostrar el pronóstico horario completo
			fmt.Println("Pronóstico horario completo:")
			for i, time := range weather.Hourly.Time {
				fmt.Printf("⏰ %s: 🌡️ %.1f°C\n", time, weather.Hourly.Temperature2m[i])
			}
		} else {
			fmt.Println("Argumento no válido. Debes usar '1' o '2'.")
		}

	} else {

		fmt.Printf("Actual weather at Rosario, Santa Fe, 🇦🇷: 🌡️ %.1f°C , 💧 %.1f%% de humedad.\n", weather.Current.TempC, weather.Current.Humidity)

	}
}

func timeParse(s string) (time.Time, error) {
	return time.Parse("2006-01-02T15:04", s)
}
