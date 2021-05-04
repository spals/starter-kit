package config

import (
	"fmt"
	nativelog "log"
	"os"
	"strconv"
	"strings"
	"time"

	descriptorpb "google.golang.org/protobuf/types/descriptorpb"

	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/proto"
)

// NewGrpcServerConfig ...
// Constructor for a Grpc server config. This is injected into
// the dependency graph using wire (see main.go)
func NewGrpcServerConfig(l envconfig.Lookuper) *proto.GrpcServerConfig {
	serverConfig := proto.GrpcServerConfig{}
	md, initErr := desc.LoadMessageDescriptorForMessage(&serverConfig)
	if initErr != nil {
		nativelog.Fatalf("GrpcServerConfig initialization failure: %s", initErr)
	}

	dynamicConfig := dynamic.NewMessage(md)
	dynamicConfig, parseErr := makeConfig(dynamicConfig, l)
	if parseErr != nil {
		nativelog.Fatalf("GrpcServerConfig parse failure: %s", parseErr)
	}

	convertErr := dynamicConfig.ConvertTo(&serverConfig)
	if convertErr != nil {
		nativelog.Fatalf("GrpcServerConfig covert failure: %s", convertErr)
	}

	log.Logger = makeLogger(&serverConfig)
	log.Debug().Interface("config", &serverConfig).Msg("HTTPServerConfig parsed")

	return &serverConfig
}

// ========== Private Helpers ==========

// Yet another string -> type converter. This is a minimal implementation which
// accounts for the types needed in the out-of-the-box configuration. More type
// convertions will be required if the configuration is fancier.
func convertValue(fd *desc.FieldDescriptor, configValue string) (val interface{}) {
	switch fd.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return configValue == "true"
	case descriptorpb.FieldDescriptorProto_TYPE_INT32,
		descriptorpb.FieldDescriptorProto_TYPE_SINT32,
		descriptorpb.FieldDescriptorProto_TYPE_UINT32:
		i, err := strconv.Atoi(configValue)
		if err != nil {
			nativelog.Fatalf("GrpcServerConfig parse failure: Unable to parse int value (%s) for field (%s): %s", configValue, fd.GetName(), err)
		}
		return int32(i)
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return configValue
	default:
		nativelog.Fatalf("GrpcServerConfig parse failure: Unable to parse value (%s) for field (%s): Unsupported type", configValue, fd.GetName())
	}

	return nil
}

func makeConfig(dynamicMessage *dynamic.Message, l envconfig.Lookuper) (*dynamic.Message, error) {
	for _, fd := range dynamicMessage.GetMessageDescriptor().GetFields() {
		if fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// dynamicMessage.SetField(fd, makeConfig)
		} else {
			configKey := strings.ToUpper(fd.GetName())
			configValue, configFound := l.Lookup(configKey)
			nativelog.Printf("GrpcServerConfig found %s as value \"%s\"", configKey, configValue)

			if configFound {
				convertedConfigValue := convertValue(fd, configValue)
				dynamicMessage.SetField(fd, convertedConfigValue)
			}
		}

	}

	return dynamicMessage, nil
}

func makeLogger(config *proto.GrpcServerConfig) zerolog.Logger {
	logLevel, err := zerolog.ParseLevel(config.GetLogLevel())
	if err != nil {
		nativelog.Fatalf("Error while parsing log level: %s. Available log levels are (trace|debug|info|warn|error|fatal|panic)", err)
	} else if logLevel == zerolog.NoLevel {
		nativelog.Fatalf("No log level configured. Please specify a log level (trace|debug|info|warn|error|fatal|panic)")
	}
	zerolog.SetGlobalLevel(logLevel)

	if config.GetDev() {
		output := zerolog.ConsoleWriter{Out: os.Stdout, TimeFormat: time.RFC3339}
		output.FormatFieldName = func(i interface{}) string {
			return fmt.Sprintf("%s:", i)
		}

		return zerolog.New(output).With().Timestamp().Logger()
	} else {
		zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
		return zerolog.New(os.Stderr).With().Timestamp().Logger()
	}
}
