package events

import (
	"context"
	"fmt"
	"net/url"
	"strings"
	"time"

	"github.com/tparker00/Vision/lookup"
	v "github.com/tparker00/Vision/types"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/event"
	"github.com/vmware/govmomi/vim25/types"
)

//Contains function iterates over the slice and returns true if the item already exists in it
func Contains(slice []string, item string) bool {
	set := make(map[string]struct{}, len(slice))
	for _, s := range slice {
		set[s] = struct{}{}
	}

	_, ok := set[item]
	return ok
}

//GetEvents return the list of VMs that have recently HA restarted (last 12 hours)
func GetEvents(ctx context.Context, config v.Config) ([]string, error) {
	var vmList []string
	escapedPassword := url.QueryEscape(config.Password)
	escapedUsername := url.QueryEscape(config.Username)

	//Use the SRV record for vCenter to lookup all the vCenters in the environment.
	servicelist, err := lookup.ServiceLookup(fmt.Sprintf("_vcenter._tcp.%s", config.Domain), fmt.Sprintf("%s", config.DnsServer))
	if err != nil {
		fmt.Printf("Error: %s", err)
		return nil, err
	}

	for _, service := range servicelist {
		fmt.Printf("Connecting to %s\n", service)
		connectionURI, err := url.Parse(fmt.Sprintf("https://%s:%s@%s/sdk", escapedUsername, escapedPassword, service))
		if err != nil {
			fmt.Printf("Error connecting: %s", err)
			return nil, err
		}
		conn, err := govmomi.NewClient(ctx, connectionURI, true)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return nil, err
		}
		defer conn.Logout(ctx)
		events := event.NewManager(conn.Client)
		filter := types.EventFilterSpec{}
		filter.EventTypeId = []string{"com.vmware.vc.ha.VmRestartedByHAEvent"}
		timeFilter := types.EventFilterSpecByTime{}
		beginTime := time.Now().Add(-12 * time.Hour)
		timeFilter.BeginTime = &beginTime
		filter.Time = &timeFilter

		eventList, err := events.QueryEvents(ctx, filter)
		if err != nil {
			fmt.Printf("Error: %s", err)
			return nil, err
		}
		for _, event := range eventList {
			orgID := strings.Split(event.GetEvent().Vm.Name, "-")[0]
			if Contains(vmList, orgID) {
				fmt.Printf("OrgID %s already in list, skipping\n", orgID)
			} else {
				vmList = append(vmList, orgID)
			}
		}
	}
	return vmList, nil
}
