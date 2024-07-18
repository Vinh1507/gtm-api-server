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
		TTL:         body.TTL,
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

func getDataCenterHistoryByDomain(domain string) ([]models.DataCenterHistory, error) {
	key := fmt.Sprintf("resource/datacenterhistory/%s", domain)

	resp, err := gtm_etcd.GetEntryByKey(key)

	if err != nil {
		return nil, err
	}

	var history []models.DataCenterHistory

	if len(resp.Kvs) > 0 {
		err = json.Unmarshal([]byte(resp.Kvs[len(resp.Kvs)-1].Value), &history)
		if err != nil {
			fmt.Println("Error Unmarshal Data Center History JSON:", err)
		}
	}

	sort.Slice(history, func(i, j int) bool {
		return history[i].TimeStamp > history[j].TimeStamp
	})

	return history, nil
}

func GetDataCenterHistory(c *gin.Context) {
	domain := c.Query("domain")
	history, err := getDataCenterHistoryByDomain(domain)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get history",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"history": history,
	})
}

func GetResolverList(c *gin.Context) {

	prefix := fmt.Sprintf("resource/domain/")

	resp, err := gtm_etcd.GetEntryByPrefix(prefix)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get entries",
		})
		return
	}

	var resolvers []models.Domain
	for _, ev := range resp.Kvs {
		var resolver models.Domain
		err = json.Unmarshal([]byte(ev.Value), &resolver)
		if err != nil {
			fmt.Println("Error unmarshalling resolver JSON:", err)
		}
		resolvers = append(resolvers, resolver)
	}

	c.JSON(http.StatusOK, gin.H{
		"resolvers": resolvers,
	})
}

func GetResolverDetail(c *gin.Context) {
	id := c.Param("id")
	prefix := fmt.Sprintf("resource/domain/")

	resp, err := gtm_etcd.GetEntryByPrefix(prefix)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get entries",
		})
		return
	}

	for _, ev := range resp.Kvs {
		var resolver models.Domain
		err = json.Unmarshal([]byte(ev.Value), &resolver)
		if err != nil {
			fmt.Println("Error unmarshalling resolver JSON:", err)
		}

		if resolver.Id == id {
			history, err := getDataCenterHistoryByDomain(resolver.DomainName)
			if err != nil {
				c.JSON(http.StatusOK, gin.H{
					"resolver":          resolver,
					"dataCenterHistory": nil,
				})
			}
			c.JSON(http.StatusOK, gin.H{
				"resolver":          resolver,
				"dataCenterHistory": history,
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"resolver":          nil,
		"dataCenterHistory": nil,
	})
}
