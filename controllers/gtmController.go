package controllers

import (
	"encoding/json"
	"fmt"
	gtm_etcd "go-api-server/etcd"
	"go-api-server/models"
	"net/http"
	"sort"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func ConfigGtm(c *gin.Context) {

	var body struct {
		Id          string
		DomainName  string
		DataCenters []models.DataCenter
		Policy      string
		TTL         int
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	dataCenterKeys := make([]string, 0)
	for _, dataCenter := range body.DataCenters {
		dataCenterKey := fmt.Sprintf("resource/datacenter/%s_%s", body.DomainName, dataCenter.Name)
		dataCenter.Domain = body.DomainName
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
		Policy:      body.Policy,
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
		"domain": domain,
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

func GetDataCenterHistory(c *gin.Context) {
	domain := c.Query("domain")

	prefix := fmt.Sprintf("resource/datacenterhistory/%s_", domain)

	fmt.Println(prefix)
	resp, err := gtm_etcd.GetEntryByPrefix(prefix)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get entries",
		})
		return
	}

	var history []models.DataCenterHistory
	for _, ev := range resp.Kvs {
		value := ev.Value
		var records []models.DataCenterHistory
		err = json.Unmarshal([]byte(value), &records)
		if err != nil {
			fmt.Println("Error marshalling Data Center History to JSON:", err)
		}
		history = append(history, records...)
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].TimeStamp > history[j].TimeStamp
	})

	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}
