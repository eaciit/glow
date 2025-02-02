package master

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/chrislusf/glow/resource"
	"github.com/chrislusf/glow/util"
)

func (tl *TeamMaster) listAgentsHandler(w http.ResponseWriter, r *http.Request) {
	util.Json(w, r, http.StatusAccepted, tl.MasterResource.Topology)
}

func (tl *TeamMaster) requestAgentHandler(w http.ResponseWriter, r *http.Request) {
	requestBlob := []byte(r.FormValue("request"))
	var request resource.AllocationRequest
	err := json.Unmarshal(requestBlob, &request)
	if err != nil {
		util.Error(w, r, http.StatusBadRequest, fmt.Sprintf("request JSON unmarshal error:%v, json:%s", err, string(requestBlob)))
		return
	}

	// fmt.Printf("request:\n%+v\n", request)

	result := tl.allocate(&request)
	// fmt.Printf("result: %v\n%+v\n", result.Error, result.Allocations)
	if result.Error != "" {
		util.Json(w, r, http.StatusNotFound, result)
		return
	}

	util.Json(w, r, http.StatusAccepted, result)

}

func (tl *TeamMaster) updateAgentHandler(w http.ResponseWriter, r *http.Request) {
	servicePortString := r.FormValue("servicePort")
	servicePort, err := strconv.Atoi(servicePortString)
	if err != nil {
		log.Printf("Strange: servicePort not found: %s, %v", servicePortString, err)
	}
	host := r.Host
	if strings.Contains(host, ":") {
		host = host[0:strings.Index(host, ":")]
	}
	// println("received agent update from", host+":"+servicePort)
	res, alloc := resource.NewComputeResourceFromRequest(r)
	ai := &resource.AgentInformation{
		Location: resource.Location{
			DataCenter: r.FormValue("dataCenter"),
			Rack:       r.FormValue("rack"),
			Server:     host,
			Port:       servicePort,
		},
		Resource:  res,
		Allocated: alloc,
	}

	// fmt.Printf("reported allocated: %v\n", alloc)

	tl.MasterResource.UpdateAgentInformation(ai)

	w.WriteHeader(http.StatusAccepted)
}
