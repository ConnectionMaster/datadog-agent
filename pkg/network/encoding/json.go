package encoding

import (
	"bytes"

	model "github.com/DataDog/agent-payload/process"
	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/http"
	"github.com/gogo/protobuf/jsonpb"
)

// ContentTypeJSON holds the HTML content-type of a JSON payload
const ContentTypeJSON = "application/json"

type jsonSerializer struct {
	marshaller jsonpb.Marshaler
}

func (j jsonSerializer) Marshal(conns *network.Connections) ([]byte, error) {
	agentConns := make([]*model.Connection, len(conns.Conns))
	domainSet := make(map[string]int)

	httpKeySet := make(map[http.Key]int)
	i := 0
	for key := range conns.HTTP {
		httpKeySet[key] = i
		i++
	}

	for i, conn := range conns.Conns {
		agentConns[i] = FormatConnection(conn, domainSet, httpKeySet)
	}

	domains := make([]string, len(domainSet))
	for k, v := range domainSet {
		domains[v] = k
	}

	httpKeys := make([]*model.HTTPKey, len(httpKeySet))
	for k, v := range httpKeySet {
		httpKeys[v] = &model.HTTPKey{
			Source: &model.Addr{Ip: k.SourceIP.String()},
			Dest:   &model.Addr{Ip: k.DestIP.String(), Port: int32(k.DestPort)},
		}
	}

	payload := connsPool.Get().(*model.Connections)
	payload.Conns = agentConns
	payload.Domains = domains
	payload.Dns = FormatDNS(conns.DNS)
	payload.HttpKeys = httpKeys
	payload.Http = FormatHTTP(conns.HTTP, httpKeySet)
	payload.Telemetry = FormatTelemetry(conns.Telemetry)

	writer := new(bytes.Buffer)
	err := j.marshaller.Marshal(writer, payload)
	returnToPool(payload)
	return writer.Bytes(), err
}

func (jsonSerializer) Unmarshal(blob []byte) (*model.Connections, error) {
	conns := new(model.Connections)
	reader := bytes.NewReader(blob)
	if err := jsonpb.Unmarshal(reader, conns); err != nil {
		return nil, err
	}
	return conns, nil
}

func (j jsonSerializer) ContentType() string {
	return ContentTypeJSON
}

var _ Marshaler = jsonSerializer{}
var _ Unmarshaler = jsonSerializer{}
