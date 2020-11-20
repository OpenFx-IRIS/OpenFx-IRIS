// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"fmt"
	"strings"
	"context"

	"github.com/pkg/errors"
	"google.golang.org/grpc"
	"github.com/kafka-connector/connector-sdk/pb"
)

// FunctionLookupBuilder builds a list of OpenFaaS functions
type FunctionLookupBuilder struct {
	GatewayURL     string
	TopicDelimiter string
}

func (s *FunctionLookupBuilder) getFunctions() (*pb.Functions, error) {
	conn, err := grpc.Dial(s.GatewayURL, grpc.WithInsecure())
	if err != nil {
		return nil, errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	functions, statusErr := client.List(context.Background(), &pb.Empty{})
	if statusErr != nil {
		return nil, errors.New("did not listed: " + statusErr.Error())
	}

	return functions, nil
}

// Build compiles a map of topic names and functions that have
// advertised to receive messages on said topic
func (s *FunctionLookupBuilder) Build() (map[string][]string, error) {
	var (
		err error
	)

	serviceMap := make(map[string][]string)

	functions, err := s.getFunctions()

	if err != nil {
		return map[string][]string{}, err
	}

	serviceMap = buildServiceMap(functions, s.TopicDelimiter, serviceMap)

	return serviceMap, err
}

func buildServiceMap(functions *pb.Functions, topicDelimiter string, serviceMap map[string][]string) map[string][]string {

	for _, function := range functions.Functions {

		if function.Annotations != nil {

			annotations := function.Annotations

			if topicNames, exist := annotations["topic"]; exist {
				if len(topicDelimiter) > 0 && strings.Count(topicNames, topicDelimiter) > 0 {

					topicSlice := strings.Split(topicNames, topicDelimiter)

					for _, topic := range topicSlice {
						serviceMap = appendServiceMap(topic, function.Name, serviceMap)
					}
				} else {
					serviceMap = appendServiceMap(topicNames, function.Name, serviceMap)
				}
			}
		}
	}
	return serviceMap
}

func appendServiceMap(key string, function string, sm map[string][]string) map[string][]string {

	key = strings.TrimSpace(key)

	if len(key) > 0 {

		if sm[key] == nil {
			sm[key] = []string{}
		}
		sep := ""

		functionPath := fmt.Sprintf("%s%s", function, sep)
		sm[key] = append(sm[key], functionPath)
	}

	return sm
}
