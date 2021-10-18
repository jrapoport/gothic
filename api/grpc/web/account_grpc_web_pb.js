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


var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js')

var google_protobuf_empty_pb = require('google-protobuf/google/protobuf/empty_pb.js')

var response_pb = require('./response_pb.js')
const proto = {};
proto.gothic = {};
proto.gothic.api = require('./account_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?grpc.web.ClientOptions} options
 * @constructor
 * @struct
 * @final
 */
proto.gothic.api.AccountClient =
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
proto.gothic.api.AccountPromiseClient =
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
 *   !proto.gothic.api.SignupRequest,
 *   !proto.gothic.api.UserResponse>}
 */
const methodDescriptor_Account_Signup = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/Signup',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.SignupRequest,
  response_pb.UserResponse,
  /**
   * @param {!proto.gothic.api.SignupRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.SignupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.signup =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/Signup',
      request,
      metadata || {},
      methodDescriptor_Account_Signup,
      callback);
};


/**
 * @param {!proto.gothic.api.SignupRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.signup =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/Signup',
      request,
      metadata || {},
      methodDescriptor_Account_Signup);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.SendConfirmRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_SendConfirmUser = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/SendConfirmUser',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.SendConfirmRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.gothic.api.SendConfirmRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.gothic.api.SendConfirmRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.sendConfirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_SendConfirmUser,
      callback);
};


/**
 * @param {!proto.gothic.api.SendConfirmRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.sendConfirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_SendConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.ConfirmUserRequest,
 *   !proto.gothic.api.BearerResponse>}
 */
const methodDescriptor_Account_ConfirmUser = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/ConfirmUser',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.ConfirmUserRequest,
  response_pb.BearerResponse,
  /**
   * @param {!proto.gothic.api.ConfirmUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.ConfirmUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.confirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/ConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmUser,
      callback);
};


/**
 * @param {!proto.gothic.api.ConfirmUserRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.confirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/ConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.LoginRequest,
 *   !proto.gothic.api.UserResponse>}
 */
const methodDescriptor_Account_Login = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/Login',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.LoginRequest,
  response_pb.UserResponse,
  /**
   * @param {!proto.gothic.api.LoginRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.LoginRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.login =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/Login',
      request,
      metadata || {},
      methodDescriptor_Account_Login,
      callback);
};


/**
 * @param {!proto.gothic.api.LoginRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.login =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/Login',
      request,
      metadata || {},
      methodDescriptor_Account_Login);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.LogoutRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_Logout = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/Logout',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.LogoutRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.gothic.api.LogoutRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.gothic.api.LogoutRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.logout =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/Logout',
      request,
      metadata || {},
      methodDescriptor_Account_Logout,
      callback);
};


/**
 * @param {!proto.gothic.api.LogoutRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.logout =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/Logout',
      request,
      metadata || {},
      methodDescriptor_Account_Logout);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.ResetPasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_SendResetPassword = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/SendResetPassword',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.ResetPasswordRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.gothic.api.ResetPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.gothic.api.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.sendResetPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/SendResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_SendResetPassword,
      callback);
};


/**
 * @param {!proto.gothic.api.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.sendResetPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/SendResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_SendResetPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.gothic.api.ConfirmPasswordRequest,
 *   !proto.gothic.api.BearerResponse>}
 */
const methodDescriptor_Account_ConfirmResetPassword = new grpc.web.MethodDescriptor(
  '/gothic.api.Account/ConfirmResetPassword',
  grpc.web.MethodType.UNARY,
  proto.gothic.api.ConfirmPasswordRequest,
  response_pb.BearerResponse,
  /**
   * @param {!proto.gothic.api.ConfirmPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  response_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.gothic.api.ConfirmPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.RpcError, ?proto.gothic.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.gothic.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.gothic.api.AccountClient.prototype.confirmResetPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/gothic.api.Account/ConfirmResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmResetPassword,
      callback);
};


/**
 * @param {!proto.gothic.api.ConfirmPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>=} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.gothic.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.gothic.api.AccountPromiseClient.prototype.confirmResetPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/gothic.api.Account/ConfirmResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmResetPassword);
};


module.exports = proto.gothic.api;

