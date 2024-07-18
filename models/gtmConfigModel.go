package models

type Domain struct {
	Id          string   `json:"id"`
	DomainName  string   `json:"domainName"`
	Policy      string   `json:"policy"`
	DataCenters []string `json:"dataCenters"`
	TTL         int      `json:"ttl"`
}

type DataCenter struct {
	Id             string       `json:"id"`
	Domain         string       `json:"domain"`
	Name           string       `json:"name"`
	IP             string       `json:"ip"`
	Status         string       `json:"status"`
	HealthCheckUrl string       `json:"healthCheckUrl"`
	Port           int          `json:"port"`
	Count          int          `json:"count"`
	Weight         int          `json:"weight"`
	IsPrimary      bool         `json:"isPrimary"`
	FailoverDelay  int          `json:"failoverDelay"`
	FailbackDelay  int          `json:"failbackDelay"`
	RankFailover   int          `json:"rankFailover"`
	LoadFeedbacks  []LoadObject `json:"loadFeedbacks"`
}

type LoadObject struct {
	DataCenterId string     `json:"dataCenterId"`
	RelativeUrl  string     `json:"relativeUrl"`
	Port         string     `json:"port"`
	TimeStamp    string     `json:"timeStamp"`
	Tag          string     `json:"tag"`
	Resources    []Resource `json:"resources"`
}

type Resource struct {
	Name        string `json:"name"`
	CurrentLoad int    `json:"currentLoad"`
	TargetLoad  int    `json:"targetLoad"`
	MaxLoad     int    `json:"maxLoad"`
}

type DataCenterHistory struct {
	DataCenterName string `json:"dataCenterName"`
	HealthCheckUrl string `json:"healthCheckUrl"`
	Domain         string `json:"domain"`
	Status         string `json:"status"`
	Reason         string `json:"reason"`
	TimeStamp      string `json:"timestamp"`
}
