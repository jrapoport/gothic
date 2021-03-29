import * as grpcWeb from 'grpc-web';

import * as response_pb from './response_pb';
import * as google_protobuf_empty_pb from 'google-protobuf/google/protobuf/empty_pb';
import * as user_pb from './user_pb';


export class UserClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getUser(
    request: user_pb.UserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: response_pb.UserResponse) => void
  ): grpcWeb.ClientReadableStream<response_pb.UserResponse>;

  updateUser(
    request: user_pb.UpdateUserRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: response_pb.UserResponse) => void
  ): grpcWeb.ClientReadableStream<response_pb.UserResponse>;

  sendConfirmUser(
    request: google_protobuf_empty_pb.Empty,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: google_protobuf_empty_pb.Empty) => void
  ): grpcWeb.ClientReadableStream<google_protobuf_empty_pb.Empty>;

  changePassword(
    request: user_pb.ChangePasswordRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.Error,
               response: response_pb.BearerResponse) => void
  ): grpcWeb.ClientReadableStream<response_pb.BearerResponse>;

}

export class UserPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  getUser(
    request: user_pb.UserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<response_pb.UserResponse>;

  updateUser(
    request: user_pb.UpdateUserRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<response_pb.UserResponse>;

  sendConfirmUser(
    request: google_protobuf_empty_pb.Empty,
    metadata?: grpcWeb.Metadata
  ): Promise<google_protobuf_empty_pb.Empty>;

  changePassword(
    request: user_pb.ChangePasswordRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<response_pb.BearerResponse>;

}

