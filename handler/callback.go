package handler

import (
	"context"
	"fmt"

	pb "github.com/dapr/dapr/pkg/proto/runtime/v1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/encoding/protojson"
	"google.golang.org/protobuf/types/known/emptypb"

	pbsub "github.com/el-zacharoo/pubsub/gen/proto/go/person/v1"
)

type CallbackServer struct {
	Server Server
	pb.UnimplementedAppCallbackServer
}

// Dapr will call this method to get the list of topics the app wants to subscribe to.
func (d CallbackServer) ListTopicSubscriptions(ctx context.Context, in *emptypb.Empty) (*pb.ListTopicSubscriptionsResponse, error) {

	fmt.Println("ListTopicSubscriptions")
	return &pb.ListTopicSubscriptionsResponse{
		Subscriptions: []*pb.TopicSubscription{{
			PubsubName: "pubsubsrv",
			Topic:      "create",
			Routes:     &pb.TopicRoutes{Default: "/create"},
		}},
	}, nil
}

// OnTopicEvent is fired for events subscribed to.
// Dapr sends published messages in a CloudEvents 0.3 envelope.
func (d CallbackServer) OnTopicEvent(ctx context.Context, in *pb.TopicEventRequest) (*pb.TopicEventResponse, error) {
	var person pbsub.Person
	if err := protojson.Unmarshal(in.Data, &person); err != nil {
		return &pb.TopicEventResponse{Status: pb.TopicEventResponse_DROP},
			status.Errorf(codes.Aborted, "issue unmarshalling data: %v", err)
	}

	switch in.Path {
	case "/create":
		fmt.Println(in, "created")
	case "/update":
		fmt.Println(in, "changed")
	default:
		return &pb.TopicEventResponse{},
			status.Errorf(codes.Aborted, "unexpected path in OnTopicEvent: %s", in.Path)
	}

	return &pb.TopicEventResponse{}, nil
}
