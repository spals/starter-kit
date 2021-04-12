package impl_test

// func TestGetLive(t *testing.T) {
// 	assert := assert.New(t)

// 	livenessConfig := proto.LivenessConfig{}
// 	config := proto.GrpcServerConfig{LivenessConfig: &livenessConfig}

// 	healthServer := impl.NewHealthServer(&config)
// 	resp, err := healthServer.GetLive(context.Background(), &proto.LiveRequest{})
// 	if assert.NoError(err) {
// 		assert.True(resp.IsLive)
// 	}
// }

// func TestGetReady(t *testing.T) {
// 	assert := assert.New(t)

// 	livenessConfig := proto.LivenessConfig{}
// 	readinessConfig := proto.ReadinessConfig{}
// 	config := proto.GrpcServerConfig{LivenessConfig: &livenessConfig, ReadinessConfig: &readinessConfig}

// 	healthServer := impl.NewHealthServer(&config)
// 	resp, err := healthServer.GetReady(context.Background(), &proto.ReadyRequest{})
// 	if assert.NoError(err) {
// 		assert.True(resp.IsReady)
// 	}
// }
