package main

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v2"
)

func GetEtcdClient() *clientv3.Client {
	cli, err := clientv3.New(clientv3.Config{
		Endpoints:   []string{"http://localhost:2379"},
		DialTimeout: 5 * time.Second,
	})
	if err != nil {
		log.Fatal(err)
	}

	return cli
}

var (
	etcdClient *clientv3.Client = GetEtcdClient()
)

func GetValueFromEtcd(etcdClient *clientv3.Client, key string) (string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := etcdClient.Get(ctx, key)
	if err != nil {
		return "", fmt.Errorf("failed to get key %s from etcd: %v", key, err)
	}

	if len(resp.Kvs) == 0 {
		return "", fmt.Errorf("key %s not found in etcd", key)
	}

	value := string(resp.Kvs[0].Value)
	return value, nil
}

func GetMicroserviceData(etcdClient *clientv3.Client) (MicroServiceData, error) {
	microservices, err := GetValuesWithPrefix(etcdClient, "/microservices/")
	if err != nil {
		log.Warn(err)
	}

	msData := MicroServiceData{
		Services: make(map[string]MicroServiceDetails),
	}

	for key, value := range microservices {
		var msDataDetails MicroServiceDetails

		err = json.Unmarshal([]byte(value), &msDataDetails)
		if err != nil {
			log.Printf("Error unmarshalling JSON: %v", err)
			return msData, err
		}

		// Trim the '/microservices/' prefix from the key
		trimmedKey := strings.TrimPrefix(key, "/microservices/")
		msData.Services[trimmedKey] = msDataDetails
	}

	return msData, nil
}

func main() {

	msdata, err := GetMicroserviceData(etcdClient)
	if err != nil {
		log.Errorf("failed  %s", err)
	}

	fmt.Println(msdata)

}

func GetValuesWithPrefix(etcdClient *clientv3.Client, prefix string) (map[string]string, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	resp, err := etcdClient.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		log.Errorf("failed to get keys with prefix %s from etcd: %v", prefix, err)
		return nil, err
	}

	if len(resp.Kvs) == 0 {
		log.Errorf("no keys with prefix %s found in etcd", prefix)
		return nil, err
	}

	values := make(map[string]string)
	for _, kv := range resp.Kvs {
		values[string(kv.Key)] = string(kv.Value)
	}
	return values, nil
}

func SplitImageAndTag(fullImageName string) (string, string) {
	splitted := strings.Split(fullImageName, ":")
	if len(splitted) == 1 {
		return splitted[0], "latest"
	}
	return splitted[0], splitted[1]
}

func (s *MicroServiceData) UnmarshalYAML(unmarshal func(interface{}) error) error {
	temp := struct {
		Services map[string]struct {
			Image    string            `yaml:"image"`
			EnvVars  map[string]string `yaml:"environment"`
			Networks []string          `yaml:"networks"`
			Secrets  []string          `yaml:"secrets"`
			Volumes  []string          `yaml:"volumes"`
			Ports    []string          `yaml:"ports,omitempty"`
			Deploy   Deploy            `yaml:"deploy"`
		} `yaml:"services"`
	}{}

	err := unmarshal(&temp)
	if err != nil {
		return err
	}

	s.Services = make(map[string]MicroServiceDetails)

	for serviceName, serviceDetails := range temp.Services {
		imageName, tag := SplitImageAndTag(serviceDetails.Image)

		volumes := make(map[string]string)
		for _, volume := range serviceDetails.Volumes {
			parts := strings.Split(volume, ":")
			if len(parts) == 2 {
				volumes[parts[0]] = parts[1]
			}
		}

		ports := make(map[string]string)
		for _, port := range serviceDetails.Ports {
			parts := strings.Split(port, ":")
			if len(parts) == 2 {
				ports[parts[0]] = parts[1]
			}
		}

		payload := MicroServiceDetails{
			Image: imageName,
			Tag:   tag,
			// EnvVars: serviceDetails.EnvVars,
			// Secrets: serviceDetails.Secrets,
			// Volumes: volumes,
			// Ports:   ports,
			// Deploy:  serviceDetails.Deploy,
		}

		s.Services[serviceName] = payload
	}

	return nil
}

// Take a given docker stack yaml file, and save all pertinent info, like the
// required env variable and volumes etc. Into etcd.
func SetMicroservicesEtcd(etcdClient EtcdClient, fileLocation string) (map[string]MicroServiceDetails, error) {
	yamlFile, err := ioutil.ReadFile(fileLocation)
	if err != nil {
		log.Errorf("Failed to read the YAML file: %v", err)
	}

	service := MicroServiceData{}
	err = yaml.Unmarshal(yamlFile, &service)
	if err != nil {
		log.Errorf("Failed to unmarshal the YAML file: %v", err)
	}

	processedServices := make(map[string]MicroServiceDetails)

	for serviceName, payload := range service.Services {
		jsonPayload, err := json.Marshal(payload)
		if err != nil {
			log.Errorf("Failed to marshal the payload to JSON: %v", err)
		}

		_, err = etcdClient.Put(context.Background(), fmt.Sprintf("/microservices/%s", serviceName), string(jsonPayload))
		if err != nil {
			log.Errorf("Failed creating service config in etcd: %s", err)
		}
		processedServices[serviceName] = payload
	}

	return processedServices, nil
}

// 	for serviceName, payload := range service.Services {
// 		payload.ImageName, payload.Tag = SplitImageAndTag(payload.ImageName)

// 		jsonPayload, err := json.Marshal(payload)
// 		if err != nil {
// 			log.Errorf("Failed to marshal the payload to JSON: %v", err)
// 		}

// 		_, err = etcdClient.Put(context.Background(), fmt.Sprintf("/microservices/%s", serviceName), string(jsonPayload))
// 		if err != nil {
// 			log.Errorf("Failed creating service config in etcd: %s", err)
// 		}
// 		processedServices[serviceName] = payload
// 	}

// 	return processedServices, nil
// }
