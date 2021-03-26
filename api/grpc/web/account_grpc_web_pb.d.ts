import * as grpcWeb from 'grpc-web';

import * as api_pb from './api_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as account_pb from './account_pb';


export class AccountClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  signup(
    request: account_pb.SignupRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: api_pb.UserResponse) => void
  ): grpcWeb.ClientReadableStream<api_pb.UserResponse>;

  sendConfirmUser(
    request: account_pb.SendConfirmRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  confirmUser(
    request: account_pb.ConfirmUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: api_pb.BearerResponse) => void
  ): grpcWeb.ClientReadableStream<api_pb.BearerResponse>;

  login(
    request: account_pb.LoginRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: api_pb.UserResponse) => void
  ): grpcWeb.ClientReadableStream<api_pb.UserResponse>;

  logout(
    request: account_pb.LogoutRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  sendResetPassword(
    request: account_pb.ResetPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  confirmResetPassword(
    request: account_pb.ConfirmPasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: api_pb.BearerResponse) => void
  ): grpcWeb.ClientReadableStream<api_pb.BearerResponse>;

  refreshBearerToken(
    request: account_pb.RefreshTokenRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: api_pb.BearerResponse) => void
  ): grpcWeb.ClientReadableStream<api_pb.BearerResponse>;

}

export class AccountPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  signup(
    request: account_pb.SignupRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_pb.UserResponse>;

  sendConfirmUser(
    request: account_pb.SendConfirmRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  confirmUser(
    request: account_pb.ConfirmUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_pb.BearerResponse>;

  login(
    request: account_pb.LoginRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_pb.UserResponse>;

  logout(
    request: account_pb.LogoutRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  sendResetPassword(
    request: account_pb.ResetPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  confirmResetPassword(
    request: account_pb.ConfirmPasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_pb.BearerResponse>;

  refreshBearerToken(
    request: account_pb.RefreshTokenRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<api_pb.BearerResponse>;

}

