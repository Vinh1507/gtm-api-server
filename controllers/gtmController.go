package controllers

import (
	"encoding/json"
	"fmt"
	gtm_etcd "go-api-server/etcd"
	"go-api-server/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConfigGtm(c *gin.Context) {

	var body struct {
		Id          string
		DomainName  string
		DataCenters []models.DataCenter
		Type        string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	dataCenterKeys := make([]string, 0)
	for _, dataCenter := range body.DataCenters {
		dataCenterKey := fmt.Sprintf("resource/datacenter/%s", dataCenter.Name)
		dataCenterJsonData, err := json.Marshal(dataCenter)
		dataCenterKeys = append(dataCenterKeys, dataCenterKey)
		if err != nil {
			fmt.Println("Error marshalling Data Center to JSON:", err)
			return
		}
		gtm_etcd.PutEntry(dataCenterKey, string(dataCenterJsonData))
	}

	domain := models.Domain{
		Id:          uuid.New().String(),
		DomainName:  body.DomainName,
		Type:        body.Type,
		DataCenters: dataCenterKeys,
	}

	jsonData, err := json.Marshal(domain)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	gtmKey := fmt.Sprintf("resource/domain/%s", body.DomainName)
	gtm_etcd.PutEntry(gtmKey, string(jsonData))

	c.JSON(http.StatusOK, gin.H{
		"body": body,
	})
}

func GetGtmConfig(c *gin.Context) {
	var body models.Domain

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	jsonData, err := json.Marshal(body)
	if err != nil {
		fmt.Println("Error marshalling to JSON:", err)
		return
	}

	gtm_etcd.PutEntry(body.DomainName, string(jsonData))

	c.JSON(http.StatusOK, gin.H{
		"body": body,
	})
}
