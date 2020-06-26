package clients

import (
	"log"
	"redis-proxy/src/clients/util"

	"redis-proxy/src/clients/client.go"
)

type Clients struct {
	file    *util.File
	clients *[]client.Client
}

func (c *Clients) convertMap(m *map[string]interface{}) *map[string]interface{} {
	var r = make(map[string]interface{})
	for k, v := range *m {
		r[k] = int(v.(float64))
	}
	return &r
}

func (c *Clients) addBiggestValuesIntoMap(final *map[string]interface{}, client *map[string]interface{}) {
	for key, value := range *client {
		switch v := value.(type) {
		case map[string]interface{}:
			m := c.convertMap(&v)
			if finalValue, ok := (*final)[key]; ok {
				finalSubValue := finalValue.(map[string]interface{})
				c.addBiggestValuesIntoMap(&finalSubValue, m)
			} else {
				(*final)[key] = *m
			}
		case float64:
			c.addBiggestValueIntoMap(final, &key, int(v))
		case int:
			c.addBiggestValueIntoMap(final, &key, v)
		default:
			return
		}
	}
}

func (c *Clients) addBiggestValueIntoMap(final *map[string]interface{}, key *string, value int) {
	if finalValue, ok := (*final)[*key]; ok {
		if finalValue.(int) < value {
			(*final)[*key] = value
		}
	} else {
		(*final)[*key] = value
	}
}

func (c *Clients) Values() *map[string]interface{} {
	final := make(map[string]interface{})

	for _, v := range *c.clients {
		if m, err := v.Value(); err == nil {
			c.addBiggestValuesIntoMap(&final, m)
		} else {
			log.Println("bad client", *v.Url)
		}
	}

	return &final
}

func (c *Clients) init() {
	for i, v := range c.file.Clients {
		(*c.clients)[i].Url = &v
	}
}

func NewClients(file *util.File) *Clients {
	list := make([]client.Client, len(file.Clients))
	clients := &Clients{file: file, clients: &list}
	clients.init()
	return clients
}
