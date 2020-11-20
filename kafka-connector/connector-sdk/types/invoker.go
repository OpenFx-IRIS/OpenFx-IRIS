// Copyright (c) OpenFaaS Author(s) 2019. All rights reserved.
// Licensed under the MIT license. See LICENSE file in the project root for full license information.

package types

import (
	"context"
	"fmt"
	"log"

	"github.com/pkg/errors"
	"github.com/kafka-connector/connector-sdk/pb"
	"google.golang.org/grpc"
)

type Invoker struct {
	PrintResponse bool
	GatewayURL    string
	Responses     chan InvokerResponse
}

type InvokerResponse struct {
	Context  context.Context
	Body     string
	Error    error
	Topic    string
	Function string
}

func NewInvoker(gatewayURL string, printResponse bool) *Invoker {
	return &Invoker{
		PrintResponse: printResponse,
		GatewayURL:    gatewayURL,
		Responses:     make(chan InvokerResponse),
	}
}

// Invoke triggers a function by accessing the API Gateway
func (i *Invoker) Invoke(topicMap *TopicMap, topic string, message *[]byte) {
	i.InvokeWithContext(context.Background(), topicMap, topic, message)
}

//InvokeWithContext triggers a function by accessing the API Gateway while propagating context
func (i *Invoker) InvokeWithContext(ctx context.Context, topicMap *TopicMap, topic string, message *[]byte) {
	if len(*message) == 0 {
		i.Responses <- InvokerResponse{
			Context: ctx,
			Error:   fmt.Errorf("no message to send"),
		}
	}

	matchedFunctions := topicMap.Match(topic)
	for _, matchedFunction := range matchedFunctions {
		log.Printf("Invoke function: %s", matchedFunction)

		res, doErr := invokefunction(ctx, matchedFunction, i.GatewayURL, *message)

		if doErr != nil {
			i.Responses <- InvokerResponse{
				Context: ctx,
				Error:   errors.Wrap(doErr, fmt.Sprintf("unable to invoke %s", matchedFunction)),
			}
			continue
		}

		i.Responses <- InvokerResponse{
			Context:  ctx,
			Body:     res,
			Function: matchedFunction,
			Topic:    topic,
		}
	}
}

func invokefunction(ctx context.Context, functionName string, gwURL string, message []byte) (string, error) {

	conn, err := grpc.Dial(gwURL, grpc.WithInsecure())
	if err != nil {
		return "", errors.New("did not connect: " + err.Error())
	}
	client := pb.NewFxGatewayClient(conn)

	msg, statusErr := client.Invoke(ctx, &pb.InvokeServiceRequest{Service: functionName, Input: message})
	if statusErr != nil {
		return "", errors.New("did not invoke: " + statusErr.Error())
	}

	return msg.Msg, nil
}
