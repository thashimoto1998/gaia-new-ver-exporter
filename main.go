package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var get_url = "https://api.github.com/repos/cosmos/gaia/releases/latest"
var setted_latest_release_url string
var setted_latest_ver string

var gaiaHasNewVerGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "gaia_has_new_ver",
		Help: "1 if gaia(CosmosHub) blockchain has new ver, 0 if no",
	},
	[]string{"exporter_service"},
)

var errorOccuredGauge = prometheus.NewGaugeVec(
	prometheus.GaugeOpts{
		Name: "gaia_new_ver_exportor_error_occured",
		Help: "1 if gaia new ver exportor error occured, 0 if no",
	},
	[]string{"exporter_service"},
)

type ErrorResponse struct {
	Message string `json:"message"`
}

func checkNewVer(w http.ResponseWriter, req *http.Request) {
	current_latest_release_url, current_latest_ver, err := getLatestRelease()
	if err != nil {
		fmt.Errorf("fail to get latest release: %v", err)
		gaiaHasNewVerGauge.With(prometheus.Labels{"exporter_service": "gaia_new_ver_exporter"}).Set(0)
		errorOccuredGauge.With(prometheus.Labels{"exporter_service": "gaia_new_ver_exporter"}).Set(1)

		errorResponse := ErrorResponse{Message: "Internal Server Error"}
		jsonResponse, _ := json.Marshal(errorResponse)
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		w.Write(jsonResponse)
		return
	}

	hasNewVer := float64(0)
	if strings.Compare(setted_latest_release_url, current_latest_release_url) != 0 || strings.Compare(setted_latest_ver, current_latest_ver) != 0 {
		hasNewVer = 1
	}

	gaiaHasNewVerGauge.With(prometheus.Labels{"exporter_service": "gaia_new_ver_exporter"}).Set(hasNewVer)
	errorOccuredGauge.With(prometheus.Labels{"exporter_service": "gaia_new_ver_exporter"}).Set(0)
	fmt.Fprint(w, "Success to check latest ver of gaia")
}

func main() {
	var err error
	setted_latest_release_url, setted_latest_ver, err = getLatestRelease()
	if err != nil {
		log.Panicf("fail to get latest rease: error %v", err)
	}
	fmt.Printf("Setted current latest release url: %s\n", setted_latest_release_url)
	fmt.Printf("Setted current latest ver: %s\n", setted_latest_ver)

	prometheus.MustRegister(gaiaHasNewVerGauge)
	prometheus.MustRegister(errorOccuredGauge)

	http.HandleFunc("/", checkNewVer)
	http.Handle("/metrics", promhttp.Handler())
	fmt.Println("gaia new ver exporter server start")
	http.ListenAndServe(":8080", nil)
}

func getLatestRelease() (string, string, error) {
	response, err := http.Get(get_url)
	if err != nil {
		return "", "", err
	}
	defer response.Body.Close()

	body, err := io.ReadAll(response.Body)
	if err != nil {
		return "", "", err
	}

	var latestRelease GitHubRelease
	if err := json.Unmarshal(body, &latestRelease); err != nil {
		return "", "", err
	}

	if len(latestRelease.URL) == 0 || len(latestRelease.Name) == 0 {
		return "", "", errors.New("fail to get name, ver from gaia latest release page")
	}

	return latestRelease.URL, latestRelease.Name, nil
}
