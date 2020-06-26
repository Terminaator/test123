package resp

import (
	"errors"
	"fmt"
	"log"
	"redis-proxy/src/clients"
	"redis-proxy/src/clients/util"
	"reflect"
	"regexp"
	"strconv"
)

type Initializer struct {
	list *clients.Clients
}

func (c *Initializer) add(m *map[string]interface{}) error {

	redis := NewRedis()

	defer redis.Close()
	re := regexp.MustCompile("\r\n")

	for k, v := range *m {
		if reflect.ValueOf(v).Kind() == reflect.Map {
			for k2, v2 := range v.(map[string]interface{}) {
				command := []byte(fmt.Sprintf("*4\r\n$4\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "HSET", len(k), k, len(k2), k2, len(strconv.Itoa(v2.(int))), v2.(int)))
				log.Println("redis add", re.ReplaceAllString(string(command), " "))

				if b := redis.Do(&command); (*b)[0] == '-' {
					return errors.New("add failed")
				}
			}
		} else {
			command := []byte(fmt.Sprintf("*3\r\n$3\r\n%s\r\n$%d\r\n%s\r\n$%d\r\n%d\r\n", "SET", len(k), k, len(strconv.Itoa(v.(int))), v.(int)))
			log.Println("redis add", re.ReplaceAllString(string(command), " "))

			if b := redis.Do(&command); (*b)[0] == '-' {
				return errors.New("add failed")
			}
		}
	}
	return nil
}

func (i *Initializer) Initialize() {
	i.add(i.list.Values())
}

func NewInitializer(file *string) *Initializer {
	return &Initializer{list: clients.NewClients(util.NewFile(file))}
}
