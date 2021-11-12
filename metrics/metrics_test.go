package metrics

import (
	"testing"
	"time"

	"github.com/liyiysng/scatter/config"
	_ "github.com/liyiysng/scatter/logger/glogger"
	"github.com/spf13/viper"
)

func TestReadWriteCount(t *testing.T) {

	c := viper.New()
	c.AddConfigPath(".")
	c.SetConfigName("testing_config")
	err := c.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(c)

	reporter, err := NewPrometheusReporter(cfg.Viper)

	if err != nil {
		t.Fatal(err)
		return
	}

	for i := 0; i < 10000; i++ {
		myLog.Info("report data")
		err = reporter.ReportCount("read_bytes_count", map[string]string{"nid": "game", "sid": "10005"}, 10)
		if err != nil {
			t.Fatal(err)
		}
		err = reporter.ReportCount("read_bytes_count", map[string]string{"nid": "game", "sid": "10006"}, 20)
		if err != nil {
			t.Fatal(err)
		}
		err = reporter.ReportCount("read_bytes_count", map[string]string{"nid": "game", "sid": "10007"}, 15)
		if err != nil {
			t.Fatal(err)
		}
		time.Sleep(time.Second)
	}
}

func TestSummer(t *testing.T) {

	c := viper.New()
	c.AddConfigPath(".")
	c.SetConfigName("testing_config")
	err := c.ReadInConfig()
	if err != nil {
		panic(err)
	}

	cfg := config.NewConfig(c)

	reporter, err := NewPrometheusReporter(cfg.Viper)

	if err != nil {
		t.Fatal(err)
		return
	}

	for i := 0; i < 10000; i++ {
		err = reporter.ReportSummary("rpc_delay", map[string]string{"nid": "game", "method": "foo1"}, 10)
		if err != nil {
			t.Fatal(err)
		}
		err = reporter.ReportSummary("rpc_delay", map[string]string{"nid": "game", "method": "foo2"}, 5)
		if err != nil {
			t.Fatal(err)
		}
		err = reporter.ReportSummary("rpc_delay", map[string]string{"nid": "game", "method": "foo3"}, 15)
		if err != nil {
			t.Fatal(err)
		}
		err = reporter.ReportSummary("rpc_delay", map[string]string{"nid": "game", "method": "foo3"}, 2)
		if err != nil {
			t.Fatal(err)
		}

		time.Sleep(time.Second)
	}
}

func TestSys(t *testing.T) {

	cfg := config.NewConfig()

	reporter, err := NewPrometheusReporter(cfg.Viper)

	if err != nil {
		t.Fatal(err)
		return
	}

	ReportSysMetrics([]Reporter{reporter}, time.Second)
}
