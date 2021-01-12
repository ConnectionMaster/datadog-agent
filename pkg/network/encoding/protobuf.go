package encoding

import (
	model "github.com/DataDog/agent-payload/process"
	"github.com/DataDog/datadog-agent/pkg/network"
	"github.com/DataDog/datadog-agent/pkg/network/http"
	"github.com/gogo/protobuf/proto"
)

// ContentTypeProtobuf holds the HTML content-type of a Protobuf payload
const ContentTypeProtobuf = "application/protobuf"

type protoSerializer struct{}

func (protoSerializer) Marshal(conns *network.Connections) ([]byte, error) {
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

	buf, err := proto.Marshal(payload)
	returnToPool(payload)
	return buf, err
}

func (protoSerializer) Unmarshal(blob []byte) (*model.Connections, error) {
	conns := new(model.Connections)
	if err := proto.Unmarshal(blob, conns); err != nil {
		return nil, err
	}
	return conns, nil
}

func (p protoSerializer) ContentType() string {
	return ContentTypeProtobuf
}

var _ Marshaler = protoSerializer{}
var _ Unmarshaler = protoSerializer{}
