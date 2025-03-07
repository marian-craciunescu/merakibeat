package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/marian-craciunescu/merakibeat/config"
	"github.com/marian-craciunescu/merakibeat/merakiclient"
)

type Merakibeat struct {
	done   chan struct{}
	config config.Config
	client beat.Client
	b      *beat.Beat
}

// Creates beater
func New(b *beat.Beat, cfg *common.Config) (beat.Beater, error) {
	c := config.DefaultConfig
	if err := cfg.Unpack(&c); err != nil {
		return nil, fmt.Errorf("Error reading config file: %v", err)
	}

	bt := &Merakibeat{
		done:   make(chan struct{}),
		config: c,
		b:      b,
	}
	return bt, nil
}

func (bt *Merakibeat) Run(b *beat.Beat) error {
	logp.Info("merakibeat is running! Hit CTRL-C to stop it.")

	var err error
	bt.client, err = b.Publisher.Connect()
	if err != nil {
		return err
	}
	if bt.config.ScanEnable == 1 {
		receiver := merakiclient.NewScanReceiver(bt.config, bt.client)
		go receiver.Run()
	}
	// health api poller
	ticker := time.NewTicker(bt.config.Period)
	poller := NewMerakiHealthPoller(bt, bt.config)

	// video api poller
	videoTicker := time.NewTicker(bt.config.VideoPeriod)
	videoPoller := NewMerakiVideoPoller(bt, bt.config)
	fmt.Printf("Period health %+v Period video %+v", bt.config.Period, bt.config.VideoPeriod)
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
			poller.Run()
		case <-videoTicker.C:
			videoPoller.Run()
		}
	}
}

func (bt *Merakibeat) Stop() {
	bt.client.Close()
	close(bt.done)
}
