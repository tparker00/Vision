package main

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strings"

	"github.com/tparker00/Vision/events"
	v "github.com/tparker00/Vision/types"

	"github.com/labstack/echo"
)

func getEvents(c echo.Context) error {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	config := getConfig()

	vmList, err := events.GetEvents(ctx, config)
	if err != nil {
		fmt.Printf("Error: %s", err)
		return c.JSON(http.StatusInternalServerError, err)
	}

	response := v.Message{Status: "OK", Message: vmList}
	return c.JSON(http.StatusOK, response)
}

func getTestEvents(c echo.Context) error {
	vmList := []string{"ORGID001-VM0001", "ORG002-VM0001", "ORG003-VM0001", "ORG001-VM0001", "ORG002-VM0021", "ORG012-VM0101"}
	fmt.Println(vmList)
	var orgIDS []string
	for _, l := range vmList {
		orgID := strings.Split(l, "-")[0]
		if events.Contains(orgIDS, orgID) {
			fmt.Printf("OrgID %s already in the list, skipping\n", orgID)
		} else {
			orgIDS = append(orgIDS, orgID)
		}
	}
	response := v.Message{Status: "OK", Message: orgIDS}
	return c.JSON(http.StatusOK, response)
}

func getConfig() v.Config {
	filename := "./vision.yaml"
	var config v.Config
	source, err := ioutil.ReadFile(filename)
	if err != nil {
		fmt.Println("Unable to load configuration file")
		os.Exit(1)
	}

	err = yaml.Unmarshal(source, &config)
	if err != nil {
		panic(err)
	}
	return config
}

func main() {
	e := echo.New()

	e.GET("/haReport", getEvents)
	e.GET("/haTest", getTestEvents)
	e.GET("/", func(c echo.Context) error { return c.JSON(http.StatusNotImplemented, "Not Implemented") })
	e.Start(":8000")

}
