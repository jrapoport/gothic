/**
 * @fileoverview gRPC-Web generated client stub for account
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

var api_pb = require('./api_pb.js')
const proto = {};
proto.account = require('./account_pb.js');

/**
 * @param {string} hostname
 * @param {?Object} credentials
 * @param {?Object} options
 * @constructor
 * @struct
 * @final
 */
proto.account.AccountClient =
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
proto.account.AccountPromiseClient =
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
 *   !proto.account.SignupRequest,
 *   !proto.api.UserResponse>}
 */
const methodDescriptor_Account_Signup = new grpc.web.MethodDescriptor(
  '/account.Account/Signup',
  grpc.web.MethodType.UNARY,
  proto.account.SignupRequest,
  api_pb.UserResponse,
  /**
   * @param {!proto.account.SignupRequest} request
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
 *   !proto.account.SignupRequest,
 *   !proto.api.UserResponse>}
 */
const methodInfo_Account_Signup = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.UserResponse,
  /**
   * @param {!proto.account.SignupRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.account.SignupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.signup =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/Signup',
      request,
      metadata || {},
      methodDescriptor_Account_Signup,
      callback);
};


/**
 * @param {!proto.account.SignupRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.signup =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/Signup',
      request,
      metadata || {},
      methodDescriptor_Account_Signup);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.SendConfirmRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_SendConfirmUser = new grpc.web.MethodDescriptor(
  '/account.Account/SendConfirmUser',
  grpc.web.MethodType.UNARY,
  proto.account.SendConfirmRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.SendConfirmRequest} request
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
 *   !proto.account.SendConfirmRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Account_SendConfirmUser = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.SendConfirmRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.account.SendConfirmRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.sendConfirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_SendConfirmUser,
      callback);
};


/**
 * @param {!proto.account.SendConfirmRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.sendConfirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/SendConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_SendConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.ConfirmUserRequest,
 *   !proto.api.BearerResponse>}
 */
const methodDescriptor_Account_ConfirmUser = new grpc.web.MethodDescriptor(
  '/account.Account/ConfirmUser',
  grpc.web.MethodType.UNARY,
  proto.account.ConfirmUserRequest,
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.ConfirmUserRequest} request
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
 *   !proto.account.ConfirmUserRequest,
 *   !proto.api.BearerResponse>}
 */
const methodInfo_Account_ConfirmUser = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.ConfirmUserRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.account.ConfirmUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.confirmUser =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/ConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmUser,
      callback);
};


/**
 * @param {!proto.account.ConfirmUserRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.confirmUser =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/ConfirmUser',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmUser);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.LoginRequest,
 *   !proto.api.UserResponse>}
 */
const methodDescriptor_Account_Login = new grpc.web.MethodDescriptor(
  '/account.Account/Login',
  grpc.web.MethodType.UNARY,
  proto.account.LoginRequest,
  api_pb.UserResponse,
  /**
   * @param {!proto.account.LoginRequest} request
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
 *   !proto.account.LoginRequest,
 *   !proto.api.UserResponse>}
 */
const methodInfo_Account_Login = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.UserResponse,
  /**
   * @param {!proto.account.LoginRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.UserResponse.deserializeBinary
);


/**
 * @param {!proto.account.LoginRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.UserResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.UserResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.login =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/Login',
      request,
      metadata || {},
      methodDescriptor_Account_Login,
      callback);
};


/**
 * @param {!proto.account.LoginRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.UserResponse>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.login =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/Login',
      request,
      metadata || {},
      methodDescriptor_Account_Login);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.LogoutRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_Logout = new grpc.web.MethodDescriptor(
  '/account.Account/Logout',
  grpc.web.MethodType.UNARY,
  proto.account.LogoutRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.LogoutRequest} request
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
 *   !proto.account.LogoutRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Account_Logout = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.LogoutRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.account.LogoutRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.logout =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/Logout',
      request,
      metadata || {},
      methodDescriptor_Account_Logout,
      callback);
};


/**
 * @param {!proto.account.LogoutRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.logout =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/Logout',
      request,
      metadata || {},
      methodDescriptor_Account_Logout);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.ResetPasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodDescriptor_Account_SendResetPassword = new grpc.web.MethodDescriptor(
  '/account.Account/SendResetPassword',
  grpc.web.MethodType.UNARY,
  proto.account.ResetPasswordRequest,
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.ResetPasswordRequest} request
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
 *   !proto.account.ResetPasswordRequest,
 *   !proto.google.protobuf.Empty>}
 */
const methodInfo_Account_SendResetPassword = new grpc.web.AbstractClientBase.MethodInfo(
  google_protobuf_empty_pb.Empty,
  /**
   * @param {!proto.account.ResetPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  google_protobuf_empty_pb.Empty.deserializeBinary
);


/**
 * @param {!proto.account.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.google.protobuf.Empty)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.google.protobuf.Empty>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.sendResetPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/SendResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_SendResetPassword,
      callback);
};


/**
 * @param {!proto.account.ResetPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.google.protobuf.Empty>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.sendResetPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/SendResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_SendResetPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.ConfirmPasswordRequest,
 *   !proto.api.BearerResponse>}
 */
const methodDescriptor_Account_ConfirmResetPassword = new grpc.web.MethodDescriptor(
  '/account.Account/ConfirmResetPassword',
  grpc.web.MethodType.UNARY,
  proto.account.ConfirmPasswordRequest,
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.ConfirmPasswordRequest} request
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
 *   !proto.account.ConfirmPasswordRequest,
 *   !proto.api.BearerResponse>}
 */
const methodInfo_Account_ConfirmResetPassword = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.ConfirmPasswordRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.account.ConfirmPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.confirmResetPassword =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/ConfirmResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmResetPassword,
      callback);
};


/**
 * @param {!proto.account.ConfirmPasswordRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.confirmResetPassword =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/ConfirmResetPassword',
      request,
      metadata || {},
      methodDescriptor_Account_ConfirmResetPassword);
};


/**
 * @const
 * @type {!grpc.web.MethodDescriptor<
 *   !proto.account.RefreshTokenRequest,
 *   !proto.api.BearerResponse>}
 */
const methodDescriptor_Account_RefreshBearerToken = new grpc.web.MethodDescriptor(
  '/account.Account/RefreshBearerToken',
  grpc.web.MethodType.UNARY,
  proto.account.RefreshTokenRequest,
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.RefreshTokenRequest} request
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
 *   !proto.account.RefreshTokenRequest,
 *   !proto.api.BearerResponse>}
 */
const methodInfo_Account_RefreshBearerToken = new grpc.web.AbstractClientBase.MethodInfo(
  api_pb.BearerResponse,
  /**
   * @param {!proto.account.RefreshTokenRequest} request
   * @return {!Uint8Array}
   */
  function(request) {
    return request.serializeBinary();
  },
  api_pb.BearerResponse.deserializeBinary
);


/**
 * @param {!proto.account.RefreshTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @param {function(?grpc.web.Error, ?proto.api.BearerResponse)}
 *     callback The callback function(error, response)
 * @return {!grpc.web.ClientReadableStream<!proto.api.BearerResponse>|undefined}
 *     The XHR Node Readable Stream
 */
proto.account.AccountClient.prototype.refreshBearerToken =
    function(request, metadata, callback) {
  return this.client_.rpcCall(this.hostname_ +
      '/account.Account/RefreshBearerToken',
      request,
      metadata || {},
      methodDescriptor_Account_RefreshBearerToken,
      callback);
};


/**
 * @param {!proto.account.RefreshTokenRequest} request The
 *     request proto
 * @param {?Object<string, string>} metadata User defined
 *     call metadata
 * @return {!Promise<!proto.api.BearerResponse>}
 *     Promise that resolves to the response
 */
proto.account.AccountPromiseClient.prototype.refreshBearerToken =
    function(request, metadata) {
  return this.client_.unaryCall(this.hostname_ +
      '/account.Account/RefreshBearerToken',
      request,
      metadata || {},
      methodDescriptor_Account_RefreshBearerToken);
};


module.exports = proto.account;

