import * as jspb from 'google-protobuf'
import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';


export class UserRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: UserRequest): UserRequest.AsObject;

    static serializeBinaryToWriter(message: UserRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): UserRequest;

    static deserializeBinaryFromReader(message: UserRequest, reader: jspb.BinaryReader): UserRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): UserRequest.AsObject;
}

export namespace UserRequest {
    export type AsObject = {}
}

export class UpdateUserRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: UpdateUserRequest): UpdateUserRequest.AsObject;

    static serializeBinaryToWriter(message: UpdateUserRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): UpdateUserRequest;

    static deserializeBinaryFromReader(message: UpdateUserRequest, reader: jspb.BinaryReader): UpdateUserRequest;

    getUsername(): string;

    setUsername(value: string): UpdateUserRequest;

    getData(): google_protobuf_struct_pb.Struct | undefined;

    setData(value?: google_protobuf_struct_pb.Struct): UpdateUserRequest;

    hasData(): boolean;

    clearData(): UpdateUserRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): UpdateUserRequest.AsObject;
}

export namespace UpdateUserRequest {
    export type AsObject = {
        username: string,
        data?: google_protobuf_struct_pb.Struct.AsObject,
    }
}

export class ChangePasswordRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: ChangePasswordRequest): ChangePasswordRequest.AsObject;

    static serializeBinaryToWriter(message: ChangePasswordRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): ChangePasswordRequest;

    static deserializeBinaryFromReader(message: ChangePasswordRequest, reader: jspb.BinaryReader): ChangePasswordRequest;

    getPassword(): string;

    setPassword(value: string): ChangePasswordRequest;

    getNewPassword(): string;

    setNewPassword(value: string): ChangePasswordRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): ChangePasswordRequest.AsObject;
}

export namespace ChangePasswordRequest {
    export type AsObject = {
        password: string,
        newPassword: string,
    }
}

