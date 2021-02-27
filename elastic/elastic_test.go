package elastic

import (
	"testing"

	"github.com/elastic/go-elasticsearch"
)

func TestElastic(t *testing.T) {
	es, err := elasticsearch.NewDefaultClient()
	if err != nil {
		t.Fatalf("Error creating the client: %s", err)
	}
	res, err := es.Info()
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	defer res.Body.Close()

	t.Logf("%v", res)

	iCreatRes, err := es.Indices.Create("scatter")
	if err != nil {
		t.Fatalf("Error getting response: %s", err)
	}
	defer iCreatRes.Body.Close()
}
