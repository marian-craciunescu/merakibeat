package beater

import (
	"strings"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	_ "github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"
	"github.com/marian-craciunescu/merakibeat/config"
	"github.com/marian-craciunescu/merakibeat/merakiclient"
)

type MerakiVideoPoller struct {
	MerakiPoller
}

func NewMerakiVideoPoller(merakibeat *Merakibeat, config config.Config) *MerakiVideoPoller {
	mc := merakiclient.NewMerakiClient(config.MerakiHost, config.MerakiKey,
		config.MerakiOrgID, config.MerakiNetworkIDs, config.Period, config.VideoPeriod)

	poller := &MerakiVideoPoller{}
	poller.merakibeat = merakibeat
	poller.config = config
	poller.mc = mc
	return poller
}

// This is function that will call MerakiClient to fetch & publish data based on
// config item.  MerakiClient should have no understanding of beats framework except
// function that returns mapstr type.
func (p *MerakiVideoPoller) Run() {
	logp.Info("%+v", p.config)

	// Publish Network Connection Event
	logp.Info("Getting Camera history for zone %+v", p.config.CameraZoneList)

	if len(p.config.CameraZoneList) > 0 {
		for _, cameraZone := range p.config.CameraZoneList {
			serialzone := strings.Split(cameraZone, ":")
			//mapStrArr, err := p.mc.GetZoneHistory(serialzone[0], serialzone[1])
			mapStrArr, err := p.mc.GetZoneRecentInfo(serialzone[0], serialzone[1])
			if err == nil {
				for _, mapStr := range mapStrArr {
					ts := time.Now()
					tsa, err := mapStr.GetValue("timestamp")
					if err == nil {
						ts, err = time.Parse("2006-01-02T15:04:05.999Z", tsa.(string))
						if err == nil {
							event := beat.Event{
								Timestamp: ts,
								Fields:    mapStr,
							}
							p.merakibeat.client.Publish(event)
							logp.Info("Camera Zone Recent info event sent")
						}
					} else {
						logp.Info("Camera Zone History info Stat event failed err %s", err.Error())
					}
				}
			}
		}
	}

}
