package beater

import (
	"fmt"
	"time"

	"github.com/elastic/beats/libbeat/beat"
	"github.com/elastic/beats/libbeat/cfgfile"
	"github.com/elastic/beats/libbeat/common"
	"github.com/elastic/beats/libbeat/logp"

	"github.com/kussj/jolokiabeat/config"
	jbcommon "github.com/kussj/jolokiabeat/common"
)

type Jolokiabeat struct {
	beatConfig 	*config.Config
	done       	chan struct{}
	period     	time.Duration
	queries		[]jbcommon.QueryConfig
	url			string
}

// Creates beater
func New() *Jolokiabeat {
	return &Jolokiabeat{
		done: make(chan struct{}),
	}
}

/// *** Beater interface methods ***///

func (bt *Jolokiabeat) Config(b *beat.Beat) error {

	// Load beater beatConfig
	err := cfgfile.Read(&bt.beatConfig, "")
	if err != nil {
		return fmt.Errorf("Error reading config file: %v", err)
	}

	var url string
	if bt.beatConfig.Jolokiabeat.Url != "" {
		bt.url = bt.beatConfig.Jolokiabeat.Url
		logp.Debug("Configured URL: %v\n", url)
	} else {
		logp.Err("URL endpoint not configured")
	}

	var queries []jbcommon.QueryConfig
	if bt.beatConfig.Jolokiabeat.Queries != nil {
		queries = bt.beatConfig.Jolokiabeat.Queries
	} else {
		logp.Err("No JMX queries configured")
	}

	bt.queries = make([]jbcommon.QueryConfig, len(queries))
//	logp.Debug("Found %d queries\n", len(queries))
	for i := 0; i < len(queries); i++ {
		q := queries[i]
		bt.queries[i] = q
	}

	return nil
}

func (bt *Jolokiabeat) Setup(b *beat.Beat) error {

	// Setting default period if not set
	if bt.beatConfig.Jolokiabeat.Period == "" {
		bt.beatConfig.Jolokiabeat.Period = "1s"
	}

	var err error
	bt.period, err = time.ParseDuration(bt.beatConfig.Jolokiabeat.Period)
	if err != nil {
		return err
	}

	return nil
}

func (bt *Jolokiabeat) Run(b *beat.Beat) error {
	logp.Info("jolokiabeat is running! Hit CTRL-C to stop it.")

	ticker := time.NewTicker(bt.period)
	counter := 1
	for {
		select {
		case <-bt.done:
			return nil
		case <-ticker.C:
		}

		metrics, err := bt.GetJMXMetrics(bt.url, bt.queries)
		if err != nil {
			logp.Err("Error reading metrics from: %v", err)
		} else {
			event := common.MapStr{
				"@timestamp": common.Time(time.Now()),
				"type":       	b.Name,
				"counter":    	counter,
				"metrics":		metrics,
				"url":			bt.url,
			}
//			b.Events.PublishEvent(event)
			logp.Info("Event sent: %v\n", event)
		}
		counter++
	}
}

func (bt *Jolokiabeat) Cleanup(b *beat.Beat) error {
	return nil
}

func (bt *Jolokiabeat) Stop() {
	close(bt.done)
}
