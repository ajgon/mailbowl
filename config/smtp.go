package config

import (
	"fmt"
	"net"
	"net/url"
	"reflect"
	"strconv"
	"strings"
	"time"
)

type SMTPAuthUser struct {
	Email        string
	PasswordHash string
}

type SMTPAuth struct {
	Enabled bool
	Users   []SMTPAuthUser
}

type SMTPLimit struct {
	Connections int
	MessageSize int
	Recipients  int
}

type SMTPListen struct {
	Proto string
	Host  string
	Port  string
}

type SMTPTimeout struct {
	Read  time.Duration
	Write time.Duration
	Data  time.Duration
}

type SMTPTLS struct {
	Key              string
	Certificate      string
	KeyFile          string
	CertificateFile  string
	ForceForStartTLS bool
}

type SMTP struct {
	Auth      SMTPAuth
	Hostname  string
	Limit     SMTPLimit
	Listen    []SMTPListen
	Timeout   SMTPTimeout
	TLS       SMTPTLS
	Whitelist []string
}

//nolint:funlen,cyclop
func SMTPHook(dataType reflect.Type, targetDataType reflect.Type, rawData interface{}) (interface{}, error) {
	var (
		data map[string]interface{}
		ok   bool
	)

	if dataType.Kind() != reflect.Map {
		return rawData, nil
	}

	if targetDataType != reflect.TypeOf(SMTP{}) {
		return rawData, nil
	}

	if data, ok = rawData.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	smtpAuth, err := buildSMTPAuth(data["auth"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtpLimit, err := buildSMTPLimit(data["limit"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtpListen, err := buildSMTPListen(data["listen"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtpTimeout, err := buildSMTPTimeout(data["timeout"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtpTLS, err := buildSMTPTLS(data["tls"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtpWhitelist, err := buildSMTPWhitelist(data["whitelist"])
	if err != nil {
		return nil, fmt.Errorf("%w", err)
	}

	smtp := SMTP{
		Auth:      *smtpAuth,
		Limit:     *smtpLimit,
		Listen:    smtpListen,
		Timeout:   *smtpTimeout,
		TLS:       *smtpTLS,
		Whitelist: smtpWhitelist,
	}

	if smtp.Hostname, ok = data["hostname"].(string); !ok {
		return nil, ErrUnserializing
	}

	return smtp, nil
}

func buildSMTPAuth(data interface{}) (*SMTPAuth, error) {
	var (
		smtpAuthMap map[string]interface{}
		ok          bool
		err         error
	)

	if smtpAuthMap, ok = data.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	smtpAuth := &SMTPAuth{}

	switch enabledValue := smtpAuthMap["enabled"].(type) {
	case bool:
		smtpAuth.Enabled = enabledValue
	case string:
		smtpAuth.Enabled = parseBoolString(enabledValue)
	default:
		return nil, ErrUnserializing
	}

	smtpAuth.Users = make([]SMTPAuthUser, 0)

	switch usersItem := smtpAuthMap["users"].(type) {
	case []interface{}:
		smtpAuth.Users, err = buildAuthUsersFromInterface(usersItem)
		if err != nil {
			return nil, fmt.Errorf("error parsing smtp.auth.users: %w", err)
		}
	case string:
		smtpAuth.Users = buildAuthUsersFromString(usersItem)
	default:
		return nil, ErrUnserializing
	}

	return smtpAuth, nil
}

func buildAuthUsersFromInterface(usersInterface []interface{}) ([]SMTPAuthUser, error) {
	smtpAuthUsers := make([]SMTPAuthUser, 0)

	for _, userInterface := range usersInterface {
		var (
			email, passwordHash string
			user                map[interface{}]interface{}
			ok                  bool
		)

		if user, ok = userInterface.(map[interface{}]interface{}); !ok {
			continue
		}

		if email, ok = user["email"].(string); !ok {
			return nil, ErrUnserializing
		}

		if passwordHash, ok = user["password_hash"].(string); !ok {
			return nil, ErrUnserializing
		}

		if email != "" && passwordHash != "" {
			smtpAuthUsers = append(smtpAuthUsers, SMTPAuthUser{Email: email, PasswordHash: passwordHash})
		}
	}

	return smtpAuthUsers, nil
}

func buildAuthUsersFromString(usersString string) []SMTPAuthUser {
	smtpAuthUsers := make([]SMTPAuthUser, 0)

	users := strings.Split(usersString, " ")
	for _, user := range users {
		lastIndex := strings.LastIndex(user, ":")

		smtpAuthUsers = append(smtpAuthUsers, SMTPAuthUser{Email: user[:lastIndex], PasswordHash: user[lastIndex+1:]})
	}

	return smtpAuthUsers
}

//nolint:cyclop
func buildSMTPLimit(limitInterface interface{}) (*SMTPLimit, error) {
	var (
		ok    bool
		limit map[string]interface{}
		err   error
	)

	smtpLimit := &SMTPLimit{}

	if limit, ok = limitInterface.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	switch conn := limit["connections"].(type) {
	case int:
		smtpLimit.Connections = conn
	case string:
		smtpLimit.Connections, err = strconv.Atoi(conn)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		return nil, ErrUnserializing
	}

	switch msgSize := limit["message_size"].(type) {
	case int:
		smtpLimit.MessageSize = msgSize
	case string:
		smtpLimit.MessageSize, err = strconv.Atoi(msgSize)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		return nil, ErrUnserializing
	}

	switch rcpts := limit["recipients"].(type) {
	case int:
		smtpLimit.Recipients = rcpts
	case string:
		smtpLimit.Recipients, err = strconv.Atoi(rcpts)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}
	default:
		return nil, ErrUnserializing
	}

	return smtpLimit, nil
}

func buildSMTPListen(urisInterface interface{}) ([]SMTPListen, error) {
	var (
		uris []string
		ok   bool
	)

	switch urisDecoded := urisInterface.(type) {
	case []string:
		uris = urisDecoded
	case []interface{}:
		for _, uriDecoded := range urisDecoded {
			var uri string

			if uri, ok = uriDecoded.(string); !ok {
				return nil, ErrUnserializing
			}

			uris = append(uris, uri)
		}
	case string:
		uris = strings.Split(urisDecoded, " ")
	default:
		return nil, ErrUnserializing
	}

	smtpListen := make([]SMTPListen, 0)

	for _, uri := range uris {
		url, err := url.Parse(uri)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		host, port, err := net.SplitHostPort(url.Host)
		if err != nil {
			return nil, fmt.Errorf("%w", err)
		}

		smtpListen = append(smtpListen, SMTPListen{Proto: url.Scheme, Host: host, Port: port})
	}

	return smtpListen, nil
}

func buildSMTPTimeout(timeoutInterface interface{}) (*SMTPTimeout, error) {
	var (
		readTimeout, writeTimeout, dataTimeout string
		ok                                     bool
		timeout                                map[string]interface{}
		err                                    error
	)

	smtpTimeout := &SMTPTimeout{}

	if timeout, ok = timeoutInterface.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	if readTimeout, ok = timeout["read"].(string); !ok {
		return nil, ErrUnserializing
	}

	smtpTimeout.Read, err = time.ParseDuration(readTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout.read: `%s`: %w", readTimeout, err)
	}

	if writeTimeout, ok = timeout["write"].(string); !ok {
		return nil, ErrUnserializing
	}

	smtpTimeout.Write, err = time.ParseDuration(writeTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout.write: `%s`: %w", writeTimeout, err)
	}

	if dataTimeout, ok = timeout["data"].(string); !ok {
		return nil, ErrUnserializing
	}

	smtpTimeout.Data, err = time.ParseDuration(dataTimeout)
	if err != nil {
		return nil, fmt.Errorf("invalid timeout.data: `%s`: %w", dataTimeout, err)
	}

	return smtpTimeout, nil
}

func buildSMTPTLS(tlsInterface interface{}) (*SMTPTLS, error) {
	var (
		tls map[string]interface{}
		ok  bool
	)

	smtpTLS := &SMTPTLS{}

	if tls, ok = tlsInterface.(map[string]interface{}); !ok {
		return nil, ErrUnserializing
	}

	if smtpTLS.Key, ok = tls["key"].(string); !ok {
		return nil, ErrUnserializing
	}

	if smtpTLS.Certificate, ok = tls["certificate"].(string); !ok {
		return nil, ErrUnserializing
	}

	if smtpTLS.KeyFile, ok = tls["key_file"].(string); !ok {
		return nil, ErrUnserializing
	}

	if smtpTLS.CertificateFile, ok = tls["certificate_file"].(string); !ok {
		return nil, ErrUnserializing
	}

	switch forceStartTLS := tls["force_for_starttls"].(type) {
	case bool:
		smtpTLS.ForceForStartTLS = forceStartTLS
	case string:
		smtpTLS.ForceForStartTLS = parseBoolString(forceStartTLS)
	default:
		return nil, ErrUnserializing
	}

	return smtpTLS, nil
}

func buildSMTPWhitelist(whitelistInterface interface{}) ([]string, error) {
	var (
		whitelist []string
		ok        bool
	)

	switch whitelistDecoded := whitelistInterface.(type) {
	case []string:
		whitelist = whitelistDecoded
	case []interface{}:
		for _, itemDecoded := range whitelistDecoded {
			var item string

			if item, ok = itemDecoded.(string); !ok {
				return nil, ErrUnserializing
			}

			whitelist = append(whitelist, item)
		}
	case string:
		whitelist = strings.Split(whitelistDecoded, " ")
	default:
		return nil, ErrUnserializing
	}

	return whitelist, nil
}
