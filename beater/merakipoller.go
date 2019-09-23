package beater

import (
	"time"

	_ "github.com/elastic/beats/libbeat/common"
	"github.com/marian-craciunescu/merakibeat/config"
	"github.com/marian-craciunescu/merakibeat/merakiclient"
)

type MerakiPoller struct {
	merakibeat *Merakibeat
	config     config.Config
	timeout    time.Duration
	mc         merakiclient.MerakiClient
}

type MerakiPolleriIntf interface {
	Run()
}
