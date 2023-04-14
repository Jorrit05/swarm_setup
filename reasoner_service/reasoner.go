package main

import (
	"encoding/json"
	"io"

	"net/http"
	"os"

	"github.com/Jorrit05/GoLib"

	clientv3 "go.etcd.io/etcd/client/v3"
)

var (
	serviceName                   = "reasoner_service"
	log, logFile                  = GoLib.InitLogger(serviceName)
	etcdClient   *clientv3.Client = GoLib.GetEtcdClient()
	hostname                      = os.Getenv("HOSTNAME")
)

type updateRequest struct {
	Type      string          `json:"type"`
	Requestor string          `json:"requestor"`
	Archetype GoLib.ArcheType `json:"archetype"`
	Reasoner  GoLib.Requestor `json:"requestor"`
}

func main() {
	defer logFile.Close()
	defer etcdClient.Close()

	registerReasonerConfiguration()

	mux := http.NewServeMux()
	mux.HandleFunc("/update", updateHandler)
	go func() {
		if err := http.ListenAndServe(":8081", mux); err != nil {
			log.Fatalf("Error starting HTTP server: %s", err)
		}
	}()

	select {}

}

func registerReasonerConfiguration() {

	archetypesJSON, err := os.ReadFile("/var/log/stack-files/config/archetype_config.json")
	if err != nil {
		log.Fatalf("Failed to read the JSON archetype config file: %v", err)
	}

	GoLib.RegisterJSONArray[GoLib.ArcheType](archetypesJSON, &GoLib.ArcheTypes{}, etcdClient, "/reasoner/archetype_config")

	reasonerConfigJSON, err := os.ReadFile("/var/log/stack-files/config/requestor_config.json")
	if err != nil {
		log.Fatalf("Failed to read the JSON requestor config file: %v", err)
	}

	GoLib.RegisterJSONArray[GoLib.Requestor](reasonerConfigJSON, &GoLib.RequestorConfig{}, etcdClient, "/reasoner/requestor_config")
}

func updateHandler(w http.ResponseWriter, req *http.Request) {
	body, err := io.ReadAll(req.Body)
	if err != nil {
		log.Printf("handler: Error reading body: %v", err)
		http.Error(w, "handler: Error reading request body", http.StatusBadRequest)
		return
	}
	defer req.Body.Close()

	// Unmarshal request body
	var updateReq updateRequest
	err = json.Unmarshal(body, &updateReq)
	if err != nil {
		http.Error(w, "Failed to unmarshal request body", http.StatusBadRequest)
		return
	}

	// Process the request based on the "type" field
	switch updateReq.Type {
	case "archeTypeUpdate":
		updateArchetypeHandler(updateReq, w, req)
	case "reasonerUpdate":
		return

	default:
		log.Printf("Unknown message type: %s", updateReq.Type)
		http.Error(w, "Unknown request", http.StatusNotFound)
		return
	}
}

func updateArchetypeHandler(updateReq updateRequest, w http.ResponseWriter, req *http.Request) {
	// Load archetypes from etcd
	archeTypeMap, err := GoLib.GetAndUnmarshalJSONMap[GoLib.ArcheType](etcdClient, "/reasoner/archetype_config/")
	if err != nil {
		http.Error(w, "Failed get data", http.StatusInternalServerError)
	}
	// Find and update the target archetype
	oldArcheType, ok := archeTypeMap[updateReq.Archetype.Name]
	if ok {
		archeTypeMap[updateReq.Archetype.Name] = updateArchetype(oldArcheType, updateReq.Archetype)

		// Save updated archetypes back to etcd
		err = GoLib.SaveStructToEtcd(etcdClient, "/reasoner/archetype_config/", archeTypeMap[updateReq.Archetype.Name])
		if err != nil {
			http.Error(w, "Failed to save updated archetypes to etcd", http.StatusInternalServerError)
			return
		}
	} else {
		err = GoLib.SaveStructToEtcd(etcdClient, "/reasoner/archetype_config/", updateReq.Archetype)
		if err != nil {
			http.Error(w, "Failed to save updated archetypes to etcd", http.StatusInternalServerError)
			return
		}
		log.Info("Saved update with name %s as new config in '/reasoner/archetype_config/'.", updateReq.Archetype.Name)
	}

	w.WriteHeader(http.StatusOK)
	w.Write([]byte("Archetype updated successfully"))

	http.Error(w, "Unknown request type", http.StatusBadRequest)
}

func updateArchetype(old, new GoLib.ArcheType) GoLib.ArcheType {
	if new.RequestType != "" {
		old.RequestType = new.RequestType
	}

	if len(new.IoConfig.ServiceIO) > 0 {
		old.IoConfig.ServiceIO = new.IoConfig.ServiceIO
	}

	if len(new.IoConfig.ThirdParty) > 0 {
		old.IoConfig.ThirdParty = new.IoConfig.ThirdParty
	}

	if new.IoConfig.Finish != "" {
		old.IoConfig.Finish = new.IoConfig.Finish
	}

	if new.IoConfig.ThirdPartyName != "" {
		old.IoConfig.ThirdPartyName = new.IoConfig.ThirdPartyName
	}

	return old
}
