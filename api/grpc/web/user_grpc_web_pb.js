/**
 * @fileoverview gRPC-Web generated client stub for user
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

var api_pb = require('./api_pb.js')
const proto = {};
proto.user = require('./user_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.user.UserClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.user.UserPromiseClient =
    function(hostname, credentials, options) {
  if (!options) options = {};
  options['format'] = 'text';

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
 *   !proto.user.UserRequest,
 *   !proto.api.UserResponse>}
 */
const methodDescriptor_User_GetUser = new grpc.web.MethodDescriptor(
  '/user.User/GetUser',
  grpc.web.MethodType.UNARY,
  proto.user.UserRequest,
  api_pb.UserResponse,
  /**
   * @param {!proto.user.UserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.user.UserRequest,
 *   !proto.api.UserResponse>}
 */
const methodInfo_User_GetUser = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.UserResponse,
  /**
   * @param {!proto.user.UserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.user.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.user.UserClient.prototype.getUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/user.User/GetUser',
      request,
      metadata || {},
      methodDescriptor_User_GetUser,
      callback);
};


/**
 * @param {!proto.user.UserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.user.UserPromiseClient.prototype.getUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/user.User/GetUser',
      request,
      metadata || {},
      methodDescriptor_User_GetUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.user.UpdateUserRequest,
 *   !proto.api.UserResponse>}
 */
const methodDescriptor_User_UpdateUser = new grpc.web.MethodDescriptor(
  '/user.User/UpdateUser',
  grpc.web.MethodType.UNARY,
  proto.user.UpdateUserRequest,
  api_pb.UserResponse,
  /**
   * @param {!proto.user.UpdateUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.user.UpdateUserRequest,
 *   !proto.api.UserResponse>}
 */
const methodInfo_User_UpdateUser = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.UserResponse,
  /**
   * @param {!proto.user.UpdateUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.user.UpdateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.user.UserClient.prototype.updateUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/user.User/UpdateUser',
      request,
      metadata || {},
      methodDescriptor_User_UpdateUser,
      callback);
};


/**
 * @param {!proto.user.UpdateUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.user.UserPromiseClient.prototype.updateUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/user.User/UpdateUser',
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
  '/user.User/SendConfirmUser',
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
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.google.protobuf.Empty,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_User_SendConfirmUser = new grpc.web.AbstractClientBase.MethodInfo(
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
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.user.UserClient.prototype.sendConfirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/user.User/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_User_SendConfirmUser,
      callback);
};


/**
 * @param {!proto.google.protobuf.Empty} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.user.UserPromiseClient.prototype.sendConfirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/user.User/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_User_SendConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.user.ChangePasswordRequest,
 *   !proto.api.BearerResponse>}
 */
const methodDescriptor_User_ChangePassword = new grpc.web.MethodDescriptor(
  '/user.User/ChangePassword',
  grpc.web.MethodType.UNARY,
  proto.user.ChangePasswordRequest,
  api_pb.BearerResponse,
  /**
   * @param {!proto.user.ChangePasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.BearerResponse.deserializeBinary
);


/**
 * @const
 * @type {!grpc.web.AbstractClientBase.MethodInfo<
 *   !proto.user.ChangePasswordRequest,
 *   !proto.api.BearerResponse>}
 */
const methodInfo_User_ChangePassword = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.BearerResponse,
  /**
   * @param {!proto.user.ChangePasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.user.ChangePasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.user.UserClient.prototype.changePassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/user.User/ChangePassword',
      request,
      metadata || {},
      methodDescriptor_User_ChangePassword,
      callback);
};


/**
 * @param {!proto.user.ChangePasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.user.UserPromiseClient.prototype.changePassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/user.User/ChangePassword',
      request,
      metadata || {},
      methodDescriptor_User_ChangePassword);
};


module.exports = proto.user;

