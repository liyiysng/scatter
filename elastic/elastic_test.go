package elastic

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"strings"
	"sync/atomic"
	"testing"
	"time"

	"golang.org/x/sync/errgroup"

	"github.com/liyiysng/scatter/config"
	"github.com/olivere/elastic/v7"
)

const (
	// es url
	url = "http://localhost:9200"
	// index name
	index = "tests"
	// need sniff
	sniff = false
	// Number of documents to bulk insert
	n = 10000
	// Number of documents to collect before committing
	bulkSize = 1000
)

func TestElastic(t *testing.T) {
	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff))
	if err != nil {
		log.Fatal(err)
	}

	// Setup a group of goroutines from the excellent errgroup package
	g, ctx := errgroup.WithContext(context.TODO())

	// The first goroutine will emit documents and send it to the second goroutine
	// via the docsc channel.
	// The second Goroutine will simply bulk insert the documents.
	type doc struct {
		ID        string    `json:"id"`
		Timestamp time.Time `json:"@timestamp"`
	}

	docsc := make(chan doc)

	begin := time.Now()

	// Goroutine to create documents
	g.Go(func() error {
		defer close(docsc)

		buf := make([]byte, 32)
		for i := 0; i < n; i++ {
			// Generate a random ID
			_, err := rand.Read(buf)
			if err != nil {
				return err
			}
			id := base64.URLEncoding.EncodeToString(buf)

			// Construct the document
			d := doc{
				ID:        id,
				Timestamp: time.Now(),
			}

			// Send over to 2nd goroutine, or cancel
			select {
			case docsc <- d:
			case <-ctx.Done():
				return ctx.Err()
			}
		}
		return nil
	})

	// Second goroutine will consume the documents sent from the first and bulk insert into ES
	var total uint64
	g.Go(func() error {
		bulk := client.Bulk().Index(index).Type("_doc")
		for d := range docsc {
			// Simple progress
			current := atomic.AddUint64(&total, 1)
			dur := time.Since(begin).Seconds()
			sec := int(dur)
			pps := int64(float64(current) / dur)
			fmt.Printf("%10d | %6d req/s | %02d:%02d\r", current, pps, sec/60, sec%60)

			// Enqueue the document
			bulk.Add(elastic.NewBulkIndexRequest().Id(d.ID).Doc(d))
			if bulk.NumberOfActions() >= bulkSize {
				// Commit
				res, err := bulk.Do(ctx)
				if err != nil {
					return err
				}
				if res.Errors {
					// Look up the failed documents with res.Failed(), and e.g. recommit
					return fmt.Errorf("bulk commit failed %v", res.Items)
				}
				// "bulk" is reset after Do, so you can reuse it
			}

			select {
			default:
			case <-ctx.Done():
				return ctx.Err()
			}
		}

		// Commit the final batch before exiting
		if bulk.NumberOfActions() > 0 {
			_, err = bulk.Do(ctx)
			if err != nil {
				return err
			}
		}
		return nil
	})

	// Wait until all goroutines are finished
	if err := g.Wait(); err != nil {
		log.Fatal(err)
	}

	// Final results
	dur := time.Since(begin).Seconds()
	sec := int(dur)
	pps := int64(float64(total) / dur)
	fmt.Printf("%10d | %6d req/s | %02d:%02d\n", total, pps, sec/60, sec%60)
}

func TestElasticCreateIndex(t *testing.T) {
	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff))
	if err != nil {
		log.Fatal(err)
	}

	indexes := client.CatIndices()

	indexRes, err := indexes.Do(context.Background())
	if err != nil {
		t.Fatal(err)
	}
	t.Log(indexRes)
	exist := false
	for i := 0; i < len(indexRes); i++ {
		ii := indexRes[i]
		if strings.HasPrefix(ii.Index, "test-") {
			exist = true
			break
		}
	}

	t.Log(exist)
	t.Log(exist)
	if !exist {
		i := client.CreateIndex("<test-{now/d}-1>")

		i.BodyString(`{"aliases":{ "test_write":{ "is_write_index": true} } }`)

		res, err := i.Do(context.Background())
		if err != nil {
			t.Fatal(err)
		}
		t.Log(res)
	}

}

func TestElasticCreateLifecyclePolicy(t *testing.T) {
	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff), elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags)))
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig()
	err = createLifecycle(client, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElasticCreateTemplate(t *testing.T) {
	// Create an Elasticsearch client
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff), elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags)))
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig()
	err = createTemplate(client, cfg)
	if err != nil {
		t.Fatal(err)
	}
}

func TestElasticSetupEnv(t *testing.T) {
	client, err := elastic.NewClient(elastic.SetURL(url), elastic.SetSniff(sniff), elastic.SetTraceLog(log.New(os.Stderr, "", log.LstdFlags)))
	if err != nil {
		log.Fatal(err)
	}
	cfg := config.NewConfig()
	err = SetupEnv(client, cfg)
	if err != nil {
		t.Fatal(err)
	}
}
