package consul

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/liyiysng/scatter/cluster/registry"
)

func encode(buf []byte) string {
	var b bytes.Buffer
	defer b.Reset()

	w := zlib.NewWriter(&b)
	if _, err := w.Write(buf); err != nil {
		return ""
	}
	w.Close()

	return hex.EncodeToString(b.Bytes())
}

func decode(d string) []byte {
	hr, err := hex.DecodeString(d)
	if err != nil {
		return nil
	}

	br := bytes.NewReader(hr)
	zr, err := zlib.NewReader(br)
	if err != nil {
		return nil
	}

	rbuf, err := ioutil.ReadAll(zr)
	if err != nil {
		return nil
	}
	zr.Close()

	return rbuf
}

func encodeEndpoints(endpoints []*registry.Endpoint, needCms bool) string {
	buf, err := json.Marshal(endpoints)
	if err != nil {
		myLog.Errorf("encode endpoints error %v", err)
		return ""
	}

	if needCms {
		return encode(buf)
	}

	return string(buf)
}

func encodeMeta(meta map[string]string, needCms bool) string {
	buf, err := json.Marshal(meta)
	if err != nil {
		myLog.Errorf("encode metas error %v", err)
		return ""
	}

	if needCms {
		return encode(buf)
	}

	return string(buf)
}

func decodeEndpoints(meta map[string]string) []*registry.Endpoint {
	compressed := false
	if cp, ok := meta["compress"]; ok {
		fmt.Sscanf(cp, "%v", &compressed)
	}

	if ep, ok := meta["endpoints"]; ok {
		var buf []byte
		if compressed {
			buf = decode(ep)
		} else {
			buf = []byte(ep)
		}
		var endpoints []*registry.Endpoint
		err := json.Unmarshal(buf, &endpoints)
		if err != nil {
			myLog.Errorf("unmarshal endpoints error %v", err)
			return nil
		}
		return endpoints
	}

	return nil
}

func decodeMeta(meta map[string]string) map[string]string {
	compressed := false
	if cp, ok := meta["compress"]; ok {
		fmt.Sscanf(cp, "%v", &compressed)
	}

	if ep, ok := meta["meta"]; ok {
		var buf []byte
		if compressed {
			buf = decode(ep)
		} else {
			buf = []byte(ep)
		}
		var nMeta map[string]string
		err := json.Unmarshal(buf, &nMeta)
		if err != nil {
			myLog.Errorf("unmarshal meta error %v", err)
			return nil
		}
		return nMeta
	}

	return nil
}
