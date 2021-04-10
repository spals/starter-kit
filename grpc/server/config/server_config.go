package config

import (
	"log"
	"os"
	"strconv"
	"strings"

	descriptorpb "google.golang.org/protobuf/types/descriptorpb"

	"github.com/golang/protobuf/jsonpb"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/sethvargo/go-envconfig"
	"github.com/spals/starter-kit/grpc/proto"
)

// NewGrpcServerConfig ...
// Constructor for a Grpc server config. This is injected into
// the dependency graph using wire (see main.go)
func NewGrpcServerConfig(l envconfig.Lookuper) *proto.GrpcServerConfig {
	log.Print("Parsing GrpcServerConfig")

	serverConfig := proto.GrpcServerConfig{}
	md, err := desc.LoadMessageDescriptorForMessage(&serverConfig)
	if err != nil {
		log.Fatalf("GrpcServerConfig lookup failure: %s", err)
		os.Exit(1)
	}

	dynamicConfig := dynamic.NewMessage(md)
	dynamicConfig, err = makeConfig(dynamicConfig, l)

	if err != nil {
		log.Fatalf("GrpcServerConfig parse failure: %s", err)
		os.Exit(1)
	}
	dynamicConfig.ConvertTo(&serverConfig)

	json := jsonpb.Marshaler{Indent: "  "}
	configJSON, _ := json.MarshalToString(&serverConfig)
	log.Printf("GrpcServerConfig parsed as \n%s", configJSON)

	return &serverConfig
}

func makeConfig(dynamicMessage *dynamic.Message, l envconfig.Lookuper) (*dynamic.Message, error) {
	for _, fd := range dynamicMessage.GetMessageDescriptor().GetFields() {
		if fd.GetType() == descriptorpb.FieldDescriptorProto_TYPE_MESSAGE {
			// dynamicMessage.SetField(fd, makeConfig)
		} else {
			configKey := strings.ToUpper(fd.GetName())
			configValue, configFound := l.Lookup(configKey)
			log.Printf("GrpcServerConfig found %s as value \"%s\"", configKey, configValue)

			if configFound {
				convertedConfigValue := convertValue(fd, configValue)
				dynamicMessage.SetField(fd, convertedConfigValue)
			}
		}

	}

	return dynamicMessage, nil
}

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
			log.Fatalf("GrpcServerConfig parse failure: Unable to parse int value (%s) for field (%s): %s", configValue, fd.GetName(), err)
			os.Exit(1)
		}
		return int32(i)
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return configValue
	default:
		log.Fatalf("GrpcServerConfig parse failure: Unable to parse value (%s) for field (%s): Unsupported type", configValue, fd.GetName())
		os.Exit(1)
	}

	return nil
}
