import * as jspb from 'google-protobuf'

import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';
import * as google_protobuf_timestamp_pb from 'google-protobuf/google/protobuf/timestamp_pb';


export class UserResponse extends jspb.Message {
    static toObject(includeInstance: boolean, msg: UserResponse): UserResponse.AsObject;

    static serializeBinaryToWriter(message: UserResponse, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): UserResponse;

    static deserializeBinaryFromReader(message: UserResponse, reader: jspb.BinaryReader): UserResponse;

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

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): UserResponse.AsObject;
}

export namespace UserResponse {
    export type AsObject = {
        role: string,
        email: string,
        username: string,
        data?: google_protobuf_struct_pb.Struct.AsObject,
        token?: BearerResponse.AsObject,
    }
}

export class BearerResponse extends jspb.Message {
    static toObject(includeInstance: boolean, msg: BearerResponse): BearerResponse.AsObject;

    static serializeBinaryToWriter(message: BearerResponse, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): BearerResponse;

    static deserializeBinaryFromReader(message: BearerResponse, reader: jspb.BinaryReader): BearerResponse;

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

