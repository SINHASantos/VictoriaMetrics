package prometheusimport

import (
	"bytes"
	"flag"
	"log"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"github.com/VictoriaMetrics/VictoriaMetrics/app/vmagent/remotewrite"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/fs"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/logger"
	"github.com/VictoriaMetrics/VictoriaMetrics/lib/protoparser/protoparserutil"
)

var (
	srv        *httptest.Server
	testOutput *bytes.Buffer
)

func TestInsertHandler(t *testing.T) {
	setUp()
	defer tearDown()
	req := httptest.NewRequest(http.MethodPost, "/insert/0/api/v1/import/prometheus", bytes.NewBufferString(`{"foo":"bar"}
go_memstats_alloc_bytes_total 1`))
	if err := InsertHandler(nil, req); err != nil {
		t.Fatalf("unexpected error %s", err)
	}
	expectedMsg := "cannot unmarshal Prometheus line"
	if !strings.Contains(testOutput.String(), expectedMsg) {
		t.Fatalf("output %q should contain %q", testOutput.String(), expectedMsg)
	}
}

func setUp() {
	srv = httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, _ *http.Request) {
		w.WriteHeader(204)
	}))
	flag.Parse()
	remoteWriteFlag := "remoteWrite.url"
	if err := flag.Lookup(remoteWriteFlag).Value.Set(srv.URL); err != nil {
		log.Fatalf("unable to set %q with value %q, err: %v", remoteWriteFlag, srv.URL, err)
	}
	logger.Init()
	protoparserutil.StartUnmarshalWorkers()
	remotewrite.Init()
	testOutput = &bytes.Buffer{}
	logger.SetOutputForTests(testOutput)
}

func tearDown() {
	protoparserutil.StopUnmarshalWorkers()
	srv.Close()
	logger.ResetOutputForTest()
	tmpDataDir := flag.Lookup("remoteWrite.tmpDataPath").Value.String()
	fs.MustRemoveDir(tmpDataDir)

}
