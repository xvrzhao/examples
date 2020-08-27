package utils

import "google.golang.org/grpc/metadata"

// MetadataCarrier is the alias of gRPC metadata.MD that implements the
// opentracing.TextMapWriter and opentracing.TextMapReader interface.
type MetadataCarrier metadata.MD

func (mc MetadataCarrier) ForeachKey(handler func(key, val string) error) error {
	for key, values := range mc {
		for _, value := range values {
			if err := handler(key, value); err != nil {
				return err
			}
		}
	}
	return nil
}

func (mc MetadataCarrier) Set(key, value string) {
	mc[key] = append(mc[key], value)
}
