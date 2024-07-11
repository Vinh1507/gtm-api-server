package controllers

import (
	"fmt"
	"go-api-server/etcd"
	"net/http"

	"github.com/gin-gonic/gin"
)

func EtcdPut(c *gin.Context) {

	var body struct {
		Key   string
		Value string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	etcd.PutEntry(body.Key, body.Value)

	resp, err := etcd.GetEntryByKey(body.Key)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get entry",
		})
		return
	}
	for _, ev := range resp.Kvs {
		fmt.Printf("Key: %s, Value: %s\n", ev.Key, ev.Value)
		c.JSON(200, gin.H{
			"key":   ev.Key,
			"value": ev.Value,
		})
		return
	}

	c.JSON(http.StatusBadRequest, gin.H{
		"message": "Failed to get entry",
	})
}

func EtcdGetByPrefix(c *gin.Context) {

	var body struct {
		Prefix string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read body",
		})
		return
	}

	resp, err := etcd.GetEntryByPrefix(body.Prefix)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get entries",
		})
		return
	}

	type Entry struct {
		Key   string
		Value string
	}

	entries := make([]Entry, 0)
	for _, ev := range resp.Kvs {
		entries = append(entries, Entry{
			Key:   string(ev.Key),
			Value: string(ev.Value),
		})
	}

	c.JSON(http.StatusOK, gin.H{
		"entries": entries,
	})
}
