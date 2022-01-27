package config

import (
	"errors"
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
)

type (
	RelayAuthMethod     int
	RelayConnectionType int
)

const (
	AuthNone RelayAuthMethod = iota
	AuthPlain
	AuthCramMD5
)

const (
	ConnectionPlain RelayConnectionType = iota
	ConnectionStartTLS
	ConnectionTLS
)

var (
	ErrInvalidAuthMethod     = errors.New("invalid auth method")
	ErrInvalidConnectionType = errors.New("invalid address protocol")
)

type RelayOutgoingServer struct {
	AuthMethod     RelayAuthMethod
	ConnectionType RelayConnectionType
	FromEmail      string
	Host           string
	Password       string
	Port           int
	Username       string
	VerifyTLS      bool
}

type Relay struct {
	OutgoingServer RelayOutgoingServer
}

func RelayHook(dataType reflect.Type, targetDataType reflect.Type, rawData interface{}) (interface{}, error) {
	var (
		data map[string]interface{}
		ok   bool
	)

	if dataType.Kind() != reflect.Map {
		return rawData, nil
	}

	if targetDataType != reflect.TypeOf(Relay{}) {
		return rawData, nil
	}

	if data, ok = rawData.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	relayOutgoingServer, err := buildOutgoingServer(data["outgoing_server"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	relayConfig := Relay{
		OutgoingServer: *relayOutgoingServer,
	}

	return relayConfig, nil
}

//nolint:cyclop
func buildOutgoingServer(outgoingServerMap interface{}) (relayOutgoingServer *RelayOutgoingServer, err error) {
	var (
		address, authMethod string
		outgoingServer      map[string]interface{}
		ok                  bool
	)

	if outgoingServer, ok = outgoingServerMap.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	if address, ok = outgoingServer["address"].(string); !ok {
		return nil, ErrUnserializing
	}

	relayOutgoingServer = &RelayOutgoingServer{}

	if address != "" {
		err = buildProtoHostPortFromAddress(relayOutgoingServer, address)
	} else {
		err = buildProtoHostPortFromData(relayOutgoingServer, outgoingServer)
	}

	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if authMethod, ok = outgoingServer["auth_method"].(string); !ok {
		return nil, ErrUnserializing
	}

	relayOutgoingServer.AuthMethod, err = buildAuthMethod(authMethod)
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	if relayOutgoingServer.FromEmail, ok = outgoingServer["from_email"].(string); !ok {
		return nil, ErrUnserializing
	}

	if relayOutgoingServer.Password, ok = outgoingServer["password"].(string); !ok {
		return nil, ErrUnserializing
	}

	if relayOutgoingServer.Username, ok = outgoingServer["username"].(string); !ok {
		return nil, ErrUnserializing
	}

	switch outgoingTLS := outgoingServer["verify_tls"].(type) {
	case bool:
		relayOutgoingServer.VerifyTLS = outgoingTLS
	case string:
		relayOutgoingServer.VerifyTLS = parseBoolString(outgoingTLS)
	default:
		return nil, ErrUnserializing
	}

	return relayOutgoingServer, nil
}

func buildAuthMethod(authMethod string) (RelayAuthMethod, error) {
	switch authMethod {
	case "none":
		return AuthNone, nil
	case "plain":
		return AuthPlain, nil
	case "crammd5":
		return AuthCramMD5, nil
	}

	return -1, ErrInvalidAuthMethod
}

func buildConnectionType(connectionType string) (RelayConnectionType, error) {
	switch connectionType {
	case "plain":
		return ConnectionPlain, nil
	case "starttls":
		return ConnectionStartTLS, nil
	case "tls":
		return ConnectionTLS, nil
	}

	return -1, ErrInvalidConnectionType
}

func buildProtoHostPortFromAddress(relayOutgoingServer *RelayOutgoingServer, address string) (err error) {
	var portString string

	addr, err := url.Parse(address)
	if err != nil {
		return fmt.Errorf("%w", err)
	}

	relayOutgoingServer.Host, portString, err = net.SplitHostPort(addr.Host)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	relayOutgoingServer.ConnectionType, err = buildConnectionType(addr.Scheme)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	relayOutgoingServer.Port, err = strconv.Atoi(portString)

	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}

func buildProtoHostPortFromData(relayOutgoingServer *RelayOutgoingServer, data map[string]interface{}) (err error) {
	var (
		connectionTypeString string
		ok                   bool
	)

	if relayOutgoingServer.Host, ok = data["host"].(string); !ok {
		return ErrUnserializing
	}

	switch portItem := data["port"].(type) {
	case int:
		relayOutgoingServer.Port = portItem
	case string:
		relayOutgoingServer.Port, err = strconv.Atoi(portItem)
		if err != nil {
			return fmt.Errorf("%w", err)
		}
	default:
		return ErrUnserializing
	}

	if connectionTypeString, ok = data["connection_type"].(string); !ok {
		return ErrUnserializing
	}

	relayOutgoingServer.ConnectionType, err = buildConnectionType(connectionTypeString)
	if err != nil {
		return fmt.Errorf("invalid address: %w", err)
	}

	return nil
}
