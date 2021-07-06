import * as jspb from 'google-protobuf'

import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class UserResponse extends jspb.Message {
  getUserId(): string;
  setUserId(value: string): UserResponse;

  getRole(): string;
  setRole(value: string): UserResponse;

  getEmail(): string;
  setEmail(value: string): UserResponse;

  getUsername(): string;
  setUsername(value: string): UserResponse;

  getData(): google_protobuf_struct_pb.Struct | undefined;
  setData(value?: google_protobuf_struct_pb.Struct): UserResponse;
  hasData(): boolean;
  clearData(): UserResponse;

  getToken(): BearerResponse | undefined;
  setToken(value?: BearerResponse): UserResponse;
  hasToken(): boolean;
  clearToken(): UserResponse;

  getTokenCase(): UserResponse.TokenCase;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): UserResponse.AsObject;
  static toObject(includeInstance: boolean, msg: UserResponse): UserResponse.AsObject;
  static serializeBinaryToWriter(message: UserResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): UserResponse;
  static deserializeBinaryFromReader(message: UserResponse, reader: jspb.BinaryReader): UserResponse;
}

export namespace UserResponse {
  export type AsObject = {
    userId: string,
    role: string,
    email: string,
    username: string,
    data?: google_protobuf_struct_pb.Struct.AsObject,
    token?: BearerResponse.AsObject,
  }

  export enum TokenCase { 
    _TOKEN_NOT_SET = 0,
    TOKEN = 6,
  }
}

export class BearerResponse extends jspb.Message {
  getType(): string;
  setType(value: string): BearerResponse;

  getAccess(): string;
  setAccess(value: string): BearerResponse;

  getRefresh(): string;
  setRefresh(value: string): BearerResponse;

  getId(): string;
  setId(value: string): BearerResponse;

  getExpiresAt(): google_protobuf_timestamp_pb.Timestamp | undefined;
  setExpiresAt(value?: google_protobuf_timestamp_pb.Timestamp): BearerResponse;
  hasExpiresAt(): boolean;
  clearExpiresAt(): BearerResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): BearerResponse.AsObject;
  static toObject(includeInstance: boolean, msg: BearerResponse): BearerResponse.AsObject;
  static serializeBinaryToWriter(message: BearerResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): BearerResponse;
  static deserializeBinaryFromReader(message: BearerResponse, reader: jspb.BinaryReader): BearerResponse;
}

export namespace BearerResponse {
  export type AsObject = {
    type: string,
    access: string,
    refresh: string,
    id: string,
    expiresAt?: google_protobuf_timestamp_pb.Timestamp.AsObject,
  }
}

export class PagedResponse extends jspb.Message {
  getIndex(): number;
  setIndex(value: number): PagedResponse;

  getSize(): number;
  setSize(value: number): PagedResponse;

  getFirst(): number;
  setFirst(value: number): PagedResponse;

  getPrev(): number;
  setPrev(value: number): PagedResponse;

  getNext(): number;
  setNext(value: number): PagedResponse;

  getLast(): number;
  setLast(value: number): PagedResponse;

  getCount(): number;
  setCount(value: number): PagedResponse;

  getTotal(): number;
  setTotal(value: number): PagedResponse;

  serializeBinary(): Uint8Array;
  toObject(includeInstance?: boolean): PagedResponse.AsObject;
  static toObject(includeInstance: boolean, msg: PagedResponse): PagedResponse.AsObject;
  static serializeBinaryToWriter(message: PagedResponse, writer: jspb.BinaryWriter): void;
  static deserializeBinary(bytes: Uint8Array): PagedResponse;
  static deserializeBinaryFromReader(message: PagedResponse, reader: jspb.BinaryReader): PagedResponse;
}

export namespace PagedResponse {
  export type AsObject = {
    index: number,
    size: number,
    first: number,
    prev: number,
    next: number,
    last: number,
    count: number,
    total: number,
  }
}

