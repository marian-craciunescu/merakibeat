package beater

import (
	"time"

	"github.com/elastic/beats/libbeat/beat"
	_ "github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/marian-craciunescu/merakibeat/config"
	"github.com/marian-craciunescu/merakibeat/merakiclient"
)

type MerakiHealthPoller struct {
	MerakiPoller
}

func NewMerakiHealthPoller(merakibeat *Merakibeat, config config.Config) *MerakiHealthPoller {
	mc := merakiclient.NewMerakiClient(config.MerakiHost, config.MerakiKey,
		config.MerakiOrgID, config.MerakiNetworkIDs, config.Period, config.VideoPeriod)
	origLen := len(config.MerakiNetworkIDs)
	if origLen == 0 {
		nwList, err := mc.GetNetworksForOrg()
		config.MerakiNetworkIDsAll = make(map[string]string)
		if err == nil {
			for _, nw := range nwList {
				logp.Info("Network Name: %s ID:%s ", nw.Name, nw.ID)
				config.MerakiNetworkIDsAll[nw.ID] = nw.Name
				// If user have not configured specific network IDs to be monitored,
				// monitor all networks in organization
				if origLen == 0 {
					config.MerakiNetworkIDs = append(config.MerakiNetworkIDs, nw.ID)
				}
			}
		} else {
			logp.Err("Failed to get network list for org %s, Err :%s ", config.MerakiOrgID, err.Error())
		}
	}

	poller := &MerakiHealthPoller{}
	poller.merakibeat = merakibeat
	poller.config = config
	poller.mc = mc
	return poller
}

// This is function that will call MerakiClient to fetch & publish data based on
// config item.  MerakiClient should have no understanding of beats framework except
// function that returns mapstr type.
func (p *MerakiHealthPoller) Run() {
	logp.Info("%+v", p.config)
	// Publish Network Connection Event
	logp.Info("Getting nw stat for network %+v", p.config.MerakiNetworkIDs)
	if p.config.NwConnStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStr, err := p.mc.GetNetworkConnectionStat(netID)
			if err == nil {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields:    mapStr,
				}
				p.merakibeat.client.Publish(event)
				logp.Info("Network Connection Stat event sent")
			}
		}
	}

	// Publish Network Latency Event
	if p.config.NwLatencyStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStr, err := p.mc.GetNetworkLatencyStat(netID)
			if err == nil {
				event := beat.Event{
					Timestamp: time.Now(),
					Fields:    mapStr,
				}
				p.merakibeat.client.Publish(event)
				logp.Info("Network Connection Stat event sent")
			}
		}
	}

	// Publish devices network stats for configured network
	if p.config.DeviceConnStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStrArr, err := p.mc.GetDevicesConnectionStat(netID)
			if err == nil {
				for j, mapStr := range mapStrArr {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    mapStr,
					}
					p.merakibeat.client.Publish(event)
					logp.Info("Device network connection Stat event sent %d", j)
				}
				logp.Info("Device network connection Stat event sent")

			}
		}
	}

	// Publish devices latency stats for configured network
	if p.config.DeviceLatencyStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStrArr, err := p.mc.GetDevicesLatencyStat(netID)
			if err == nil {
				for j, mapStr := range mapStrArr {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    mapStr,
					}
					p.merakibeat.client.Publish(event)
					logp.Info("Device network latency Stat event sent %d", j)
				}
				logp.Info("Device network latency Stat event sent")
			}
		}
	}

	// Publish client network stats for configured network
	if p.config.ClientConnStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStrArr, err := p.mc.GetClientConnectionStat(netID)
			if err == nil {
				for j, mapStr := range mapStrArr {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    mapStr,
					}
					p.merakibeat.client.Publish(event)
					logp.Info("Client network connection Stat event sent %d", j)
				}
				logp.Info("Client network connection Stat event sent")

			}
		}
	}

	// Publish client latency stats for configured network
	if p.config.ClientLatencyStat != 0 {
		for _, netID := range p.config.MerakiNetworkIDs {
			mapStrArr, err := p.mc.GetClientLatencyStat(netID)
			if err == nil {
				for j, mapStr := range mapStrArr {
					event := beat.Event{
						Timestamp: time.Now(),
						Fields:    mapStr,
					}
					p.merakibeat.client.Publish(event)
					logp.Info("Client network latency Stat event sent %d", j)
				}
				logp.Info("Client network latency Stat event sent")
			}
		}
	}
}
