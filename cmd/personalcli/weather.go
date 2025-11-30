package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url" // Import the net/url package
	"os"
	"strings" // Import the strings package

	"github.com/spf13/cobra"
)

// GeocodingResponse defines the structure for the Geocoding API response.
type GeocodingResponse struct {
	Lat float64 `json:"lat"`
	Lon float64 `json:"lon"`
}

// WeatherResponse defines the structure for the Current Weather API response.
type WeatherResponse struct {
	Name string `json:"name"`
	Main struct {
		Temp      float64 `json:"temp"`
		FeelsLike float64 `json:"feels_like"`
		Humidity  int     `json:"humidity"`
	} `json:"main"`
	Weather []struct {
		Main        string `json:"main"`
		Description string `json:"description"`
	} `json:"weather"`
	Wind struct {
		Speed float64 `json:"speed"`
	} `json:"wind"`
}

func init() {
	rootCmd.AddCommand(weatherCmd)
	weatherCmd.Flags().StringP("zip", "z", "", "Zip code for weather lookup")
	weatherCmd.Flags().StringP("location", "l", "", "City name or 'City, State' for weather lookup")
	weatherCmd.Flags().StringP("api-key", "k", "", "OpenWeatherMap API key")
}

// ZipGeocodingResponse defines the structure for the Geocoding API response for zip codes.
type ZipGeocodingResponse struct {
	Zip  string `json:"zip"`
	Name string `json:"name"`
	Lat  float64 `json:"lat"`
	Lon  float64 `json:"lon"`
}

var weatherCmd = &cobra.Command{
	Use:   "weather",
	Short: "Get the current weather for a specific location",
	Long:  `Fetches and displays the current weather for a given city, 'City, State' or zip code.`,
	Run: func(cmd *cobra.Command, args []string) {
		apiKey, _ := cmd.Flags().GetString("api-key")
		if apiKey == "" {
			apiKey = os.Getenv("WEATHER_API_KEY")
		}

		if apiKey == "" {
			fmt.Println("Error: OpenWeatherMap API key not provided.")
			fmt.Println("Please provide it using the --api-key flag or by setting the WEATHER_API_KEY environment variable.")
			fmt.Println("You can get a free API key from https://openweathermap.org/")
			os.Exit(1)
		}

		zipCode, _ := cmd.Flags().GetString("zip")
		locationName, _ := cmd.Flags().GetString("location")

		var lat, lon float64
		var displayLocation string
		var err error

		if zipCode != "" {
			displayLocation = zipCode
			fmt.Printf("Fetching weather for zip code %s...\n", zipCode)
			// Geocoding for zip code
			geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/zip?zip=%s,US&appid=%s", zipCode, apiKey)
			resp, err := http.Get(geoURL)
			if err != nil {
				log.Fatalf("Error fetching geocoding data for zip code: %v", err)
			}
			defer resp.Body.Close()

			body, err := io.ReadAll(resp.Body)
			if err != nil {
				log.Fatalf("Error reading geocoding response body for zip code: %v", err)
			}

			var zipGeoResponse ZipGeocodingResponse
			if err := json.Unmarshal(body, &zipGeoResponse); err != nil {
				log.Fatalf("Error decoding geocoding response for zip code: %v\nRaw response: %s", err, string(body))
			}

			if zipGeoResponse.Zip == "" { // Check if zip was found
				fmt.Println("Could not find location for zip code:", zipCode)
				os.Exit(1)
			}

			lat = zipGeoResponse.Lat
			lon = zipGeoResponse.Lon
			displayLocation = zipGeoResponse.Name // Use city name from zip code response
		} else if locationName != "" {
			displayLocation = locationName
			fmt.Printf("Fetching weather for %s...\n", locationName)

			var geoResponse []GeocodingResponse
			foundLocation := false

			// Attempt 1: Try with ",US" appended if a comma is present (for "City, State" formats)
			if strings.Contains(locationName, ",") {
				fmt.Printf("Trying '%s,US' for better geocoding...\n", locationName)
				geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s,US&limit=1&appid=%s", url.QueryEscape(locationName), apiKey)
				resp, err := http.Get(geoURL)
				if err != nil {
					log.Printf("Warning: Error fetching geocoding data for '%s,US': %v", locationName, err)
				} else {
					defer resp.Body.Close()
					body, err := io.ReadAll(resp.Body)
					if err != nil {
						log.Printf("Warning: Error reading geocoding response body for '%s,US': %v", locationName, err)
					} else {
						if err := json.Unmarshal(body, &geoResponse); err != nil {
							log.Printf("Warning: Error decoding geocoding response for '%s,US': %v\nRaw response: %s", locationName, err, string(body))
						} else if len(geoResponse) > 0 {
							foundLocation = true
						}
					}
				}
			}

			// If not found in Attempt 1 or no comma was present, try with original locationName
			if !foundLocation {
				fmt.Printf("Trying original location name '%s'...\n", locationName)
				geoURL := fmt.Sprintf("http://api.openweathermap.org/geo/1.0/direct?q=%s&limit=1&appid=%s", url.QueryEscape(locationName), apiKey)
				resp, err := http.Get(geoURL)
				if err != nil {
					log.Fatalf("Error fetching geocoding data for location: %v", err)
				}
				defer resp.Body.Close()

				body, err := io.ReadAll(resp.Body)
				if err != nil {
					log.Fatalf("Error reading geocoding response body for location: %v", err)
				}

				if err := json.Unmarshal(body, &geoResponse); err != nil {
					log.Fatalf("Error decoding geocoding response for location: %v\nRaw response: %s", err, string(body))
				}

				if len(geoResponse) > 0 {
					foundLocation = true
				}
			}

			if !foundLocation {
				fmt.Println("Could not find location:", locationName)
				os.Exit(1)
			}

			lat = geoResponse[0].Lat
			lon = geoResponse[0].Lon
		} else {
			fmt.Println("Error: Please provide a zip code (-z) or a location (-l).")
			os.Exit(1)
		}

		// 2. Weather: Get current weather using coordinates (imperial for Fahrenheit)
		weatherURLImperial := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=imperial", lat, lon, apiKey)
		respImperial, err := http.Get(weatherURLImperial)
		if err != nil {
			log.Fatalf("Error fetching imperial weather data: %v", err)
		}
		defer respImperial.Body.Close()

		var weatherResponseImperial WeatherResponse
		if err := json.NewDecoder(respImperial.Body).Decode(&weatherResponseImperial); err != nil {
			log.Fatalf("Error decoding imperial weather response: %v", err)
		}

		// 3. Weather: Get current weather using coordinates (metric for Celsius)
		weatherURLMetric := fmt.Sprintf("http://api.openweathermap.org/data/2.5/weather?lat=%f&lon=%f&appid=%s&units=metric", lat, lon, apiKey)
		respMetric, err := http.Get(weatherURLMetric)
		if err != nil {
			log.Fatalf("Error fetching metric weather data: %v", err)
		}
		defer respMetric.Body.Close()

		var weatherResponseMetric WeatherResponse
		if err := json.NewDecoder(respMetric.Body).Decode(&weatherResponseMetric); err != nil {
			log.Fatalf("Error decoding metric weather response: %v", err)
		}


		// 4. Display Weather
		fmt.Printf("\nWeather in %s:\n", displayLocation)
		fmt.Printf("Temperature: %.1f°F (Feels like %.1f°F) | %.1f°C (Feels like %.1f°C)\n",
			weatherResponseImperial.Main.Temp, weatherResponseImperial.Main.FeelsLike,
			weatherResponseMetric.Main.Temp, weatherResponseMetric.Main.FeelsLike)
		if len(weatherResponseImperial.Weather) > 0 { // Weather conditions should be the same for both units
			fmt.Printf("Condition: %s (%s)\n", weatherResponseImperial.Weather[0].Main, weatherResponseImperial.Weather[0].Description)
		}
		fmt.Printf("Humidity: %d%%\n", weatherResponseImperial.Main.Humidity)
		fmt.Printf("Wind Speed: %.1f mph | %.1f m/s\n", weatherResponseImperial.Wind.Speed, weatherResponseMetric.Wind.Speed)
	},
}
