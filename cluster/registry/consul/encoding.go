package consul

import (
	"bytes"
	"compress/zlib"
	"encoding/hex"
	"encoding/json"
	"io/ioutil"
	"strings"

	"github.com/liyiysng/scatter/cluster/registry"
)

const (
	_endPointPrefix = "_ep_"
	_nodeMetaPrefix = "_nm_"
	_srvMetaPrefix  = "_sm_"
	_srvVersion     = "_version_"
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

// 合并所有元数据
func encodeMetaData(srv *registry.Service) map[string]string {

	if len(srv.Nodes) != 1 {
		myLog.Error("encodeMetaData len(srv.Nodes) != 1")
		return nil
	}

	meta := make(map[string]string)
	// endpoints meta
	for _, v := range srv.Endpoints {
		if len(v.Metadata) == 0 {
			continue
		}
		buf, err := json.Marshal(v.Metadata)
		if err != nil {
			myLog.Errorf("encode metas error %v", err)
			continue
		}
		meta[_endPointPrefix+v.Name] = string(buf)
	}
	node := srv.Nodes[0]
	// node meta
	for k, v := range node.Metadata {
		meta[_nodeMetaPrefix+k] = v
	}
	// srv meta
	for k, v := range srv.Metadata {
		meta[_srvMetaPrefix+k] = v
	}

	meta[_srvVersion] = srv.Version

	return meta
}

func decodeVersion(meta map[string]string) string {
	return meta[_srvVersion]
}

func decodeEndpoints(meta map[string]string) []*registry.Endpoint {

	eps := make([]*registry.Endpoint, 0)

	for k, v := range meta {
		if strings.HasPrefix(k, _endPointPrefix) {
			var eMeta map[string]string
			err := json.Unmarshal([]byte(v), &eMeta)
			if err != nil {
				myLog.Errorf("unmarshal endpoint meta error %v", err)
				continue
			}
			eps = append(eps, &registry.Endpoint{
				Name:     strings.TrimLeft(k, _endPointPrefix),
				Metadata: eMeta,
			})
		}
	}

	return eps
}

func decodeNodeMeta(meta map[string]string) map[string]string {

	nMeta := map[string]string{}

	for k, v := range meta {
		if strings.HasPrefix(k, _nodeMetaPrefix) {
			nMeta[strings.TrimLeft(k, _nodeMetaPrefix)] = v
		}
	}

	return nMeta
}

func decodeSrvMeta(meta map[string]string) map[string]string {

	sMeta := map[string]string{}

	for k, v := range meta {
		if strings.HasPrefix(k, _srvMetaPrefix) {
			sMeta[strings.TrimLeft(k, _srvMetaPrefix)] = v
		}
	}

	return sMeta
}
