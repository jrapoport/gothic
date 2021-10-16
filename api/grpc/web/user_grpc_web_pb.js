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

var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js')

var response_pb = require('./response_pb.js')
const proto = {};
proto.gothic = {};
proto.gothic.api = require('./user_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gothic.api.UserClient =
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
proto.gothic.api.UserPromiseClient =
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
 *   !proto.gothic.api.UserRequest,
 *   !proto.gothic.api.UserResponse>}
 */
const methodDescriptor_User_GetUser = new grpc.web.MethodDescriptor(
  '/gothic.api.User/GetUser',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.UserRequest,
  response_pb.UserResponse,
  /**
   * @param {!proto.gothic.api.UserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.UserClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.User/GetUser',
      request,
      metadata || {},
      methodDescriptor_User_GetUser,
      callback);
};


/**
 * @param {!proto.gothic.api.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.UserPromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.User/GetUser',
      request,
      metadata || {},
      methodDescriptor_User_GetUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.UpdateUserRequest,
 *   !proto.gothic.api.UserResponse>}
 */
const methodDescriptor_User_UpdateUser = new grpc.web.MethodDescriptor(
  '/gothic.api.User/UpdateUser',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.UpdateUserRequest,
  response_pb.UserResponse,
  /**
   * @param {!proto.gothic.api.UpdateUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.UpdateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.UserClient.prototype.updateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.User/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_User_UpdateUser,
      callback);
};


/**
 * @param {!proto.gothic.api.UpdateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.UserPromiseClient.prototype.updateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.User/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_User_UpdateUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_User_SendConfirmUser = new grpc.web.MethodDescriptor(
  '/gothic.api.User/SendConfirmUser',
  grpc.web.MethodType.UNARY,
  google_protobuf_empty_pb.Empty,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.google.protobuf.Empty} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.UserClient.prototype.sendConfirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.User/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_User_SendConfirmUser,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.gothic.api.UserPromiseClient.prototype.sendConfirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.User/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_User_SendConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.ChangePasswordRequest,
 *   !proto.gothic.api.BearerResponse>}
 */
const methodDescriptor_User_ChangePassword = new grpc.web.MethodDescriptor(
  '/gothic.api.User/ChangePassword',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.ChangePasswordRequest,
  response_pb.BearerResponse,
  /**
   * @param {!proto.gothic.api.ChangePasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.ChangePasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.UserClient.prototype.changePassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.User/ChangePassword',
      request,
      metadata || {},
      methodDescriptor_User_ChangePassword,
      callback);
};


/**
 * @param {!proto.gothic.api.ChangePasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.UserPromiseClient.prototype.changePassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.User/ChangePassword',
      request,
      metadata || {},
      methodDescriptor_User_ChangePassword);
};


module.exports = proto.gothic.api;

