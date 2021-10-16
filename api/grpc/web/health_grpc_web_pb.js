/**
 * @fileoverview gRPC-Web generated client stub for gothic.api
 * @enhanceable
 * @public
 */

// GENERATED CODE -- DO NOT EDIT!


/* eslint-disable */
// @ts-nocheck



const grpc = {};
grpc.web = require('grpc-web');


var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')
const proto = {};
proto.gothic = {};
proto.gothic.api = require('./health_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gothic.api.HealthClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gothic.api.HealthPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options.format = 'text';

  /**
   * @private @const {!grpc.web.GrpcWebClientBase} The client
   */
  this.client_ = new grpc.web.GrpcWebClientBase(options);

  /**
   * @private @const {string} The hostname
   */
  this.hostname_ = hostname;

};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.gothic.api.HealthCheckResponse>}
 */
const methodDescriptor_Health_HealthCheck = new grpc.web.MethodDescriptor(
  '/gothic.api.Health/HealthCheck',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  proto.gothic.api.HealthCheckResponse,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  proto.gothic.api.HealthCheckResponse.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.HealthCheckResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.HealthCheckResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.HealthClient.prototype.healthCheck =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Health/HealthCheck',
      request,
      metadata || {},
      methodDescriptor_Health_HealthCheck,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.HealthCheckResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.HealthPromiseClient.prototype.healthCheck =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Health/HealthCheck',
      request,
      metadata || {},
      methodDescriptor_Health_HealthCheck);
};


module.exports = proto.gothic.api;

