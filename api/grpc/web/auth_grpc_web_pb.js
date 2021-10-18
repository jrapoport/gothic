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


var response_pb = require('./response_pb.js')
const proto = {};
proto.gothic = {};
proto.gothic.api = require('./auth_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gothic.api.AuthClient =
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
proto.gothic.api.AuthPromiseClient =
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
 *   !proto.gothic.api.RefreshTokenRequest,
 *   !proto.gothic.api.BearerResponse>}
 */
const methodDescriptor_Auth_RefreshBearerToken = new grpc.web.MethodDescriptor(
  '/gothic.api.Auth/RefreshBearerToken',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.RefreshTokenRequest,
  response_pb.BearerResponse,
  /**
   * @param {!proto.gothic.api.RefreshTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.RefreshTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AuthClient.prototype.refreshBearerToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Auth/RefreshBearerToken',
      request,
      metadata || {},
      methodDescriptor_Auth_RefreshBearerToken,
      callback);
};


/**
 * @param {!proto.gothic.api.RefreshTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AuthPromiseClient.prototype.refreshBearerToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Auth/RefreshBearerToken',
      request,
      metadata || {},
      methodDescriptor_Auth_RefreshBearerToken);
};


module.exports = proto.gothic.api;

