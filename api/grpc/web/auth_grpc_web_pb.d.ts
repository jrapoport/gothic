import * as grpcWeb from 'grpc-web';

import * as response_pb from './response_pb';
import * as auth_pb from './auth_pb';


export class AuthClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  refreshBearerToken(
    request: auth_pb.RefreshTokenRequest,
    metadata: grpcWeb.Metadata | undefined,
    callback: (err: grpcWeb.RpcError,
               response: response_pb.BearerResponse) => void
  ): grpcWeb.ClientReadableStream<response_pb.BearerResponse>;

}

export class AuthPromiseClient {
  constructor (hostname: string,
               credentials?: null | { [index: string]: string; },
               options?: null | { [index: string]: any; });

  refreshBearerToken(
    request: auth_pb.RefreshTokenRequest,
    metadata?: grpcWeb.Metadata
  ): Promise<response_pb.BearerResponse>;

}

