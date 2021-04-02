import * as jspb from 'google-protobuf'

import * as google_protobuf_struct_pb from 'google-protobuf/google/protobuf/struct_pb';


export class SignupRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: SignupRequest): SignupRequest.AsObject;

    static serializeBinaryToWriter(message: SignupRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): SignupRequest;

    static deserializeBinaryFromReader(message: SignupRequest, reader: jspb.BinaryReader): SignupRequest;

    getEmail(): string;

    setEmail(value: string): SignupRequest;

    getPassword(): string;

    setPassword(value: string): SignupRequest;

    getUsername(): string;

    setUsername(value: string): SignupRequest;

    getCode(): string;

    setCode(value: string): SignupRequest;

    getData(): google_protobuf_struct_pb.Struct | undefined;

    setData(value?: google_protobuf_struct_pb.Struct): SignupRequest;

    hasData(): boolean;

    clearData(): SignupRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): SignupRequest.AsObject;
}

export namespace SignupRequest {
    export type AsObject = {
        email: string,
        password: string,
        username: string,
        code: string,
        data?: google_protobuf_struct_pb.Struct.AsObject,
    }
}

export class SendConfirmRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: SendConfirmRequest): SendConfirmRequest.AsObject;

    static serializeBinaryToWriter(message: SendConfirmRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): SendConfirmRequest;

    static deserializeBinaryFromReader(message: SendConfirmRequest, reader: jspb.BinaryReader): SendConfirmRequest;

    getEmail(): string;

    setEmail(value: string): SendConfirmRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): SendConfirmRequest.AsObject;
}

export namespace SendConfirmRequest {
    export type AsObject = {
        email: string,
    }
}

export class ConfirmUserRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: ConfirmUserRequest): ConfirmUserRequest.AsObject;

    static serializeBinaryToWriter(message: ConfirmUserRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): ConfirmUserRequest;

    static deserializeBinaryFromReader(message: ConfirmUserRequest, reader: jspb.BinaryReader): ConfirmUserRequest;

    getToken(): string;

    setToken(value: string): ConfirmUserRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): ConfirmUserRequest.AsObject;
}

export namespace ConfirmUserRequest {
    export type AsObject = {
        token: string,
    }
}

export class LoginRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: LoginRequest): LoginRequest.AsObject;

    static serializeBinaryToWriter(message: LoginRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): LoginRequest;

    static deserializeBinaryFromReader(message: LoginRequest, reader: jspb.BinaryReader): LoginRequest;

    getEmail(): string;

    setEmail(value: string): LoginRequest;

    getPassword(): string;

    setPassword(value: string): LoginRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): LoginRequest.AsObject;
}

export namespace LoginRequest {
    export type AsObject = {
        email: string,
        password: string,
    }
}

export class LogoutRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: LogoutRequest): LogoutRequest.AsObject;

    static serializeBinaryToWriter(message: LogoutRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): LogoutRequest;

    static deserializeBinaryFromReader(message: LogoutRequest, reader: jspb.BinaryReader): LogoutRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): LogoutRequest.AsObject;
}

export namespace LogoutRequest {
    export type AsObject = {}
}

export class ResetPasswordRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: ResetPasswordRequest): ResetPasswordRequest.AsObject;

    static serializeBinaryToWriter(message: ResetPasswordRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): ResetPasswordRequest;

    static deserializeBinaryFromReader(message: ResetPasswordRequest, reader: jspb.BinaryReader): ResetPasswordRequest;

    getEmail(): string;

    setEmail(value: string): ResetPasswordRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): ResetPasswordRequest.AsObject;
}

export namespace ResetPasswordRequest {
    export type AsObject = {
        email: string,
    }
}

export class ConfirmPasswordRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: ConfirmPasswordRequest): ConfirmPasswordRequest.AsObject;

    static serializeBinaryToWriter(message: ConfirmPasswordRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): ConfirmPasswordRequest;

    static deserializeBinaryFromReader(message: ConfirmPasswordRequest, reader: jspb.BinaryReader): ConfirmPasswordRequest;

    getPassword(): string;

    setPassword(value: string): ConfirmPasswordRequest;

    getToken(): string;

    setToken(value: string): ConfirmPasswordRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): ConfirmPasswordRequest.AsObject;
}

export namespace ConfirmPasswordRequest {
    export type AsObject = {
        password: string,
        token: string,
    }
}

export class RefreshTokenRequest extends jspb.Message {
    static toObject(includeInstance: boolean, msg: RefreshTokenRequest): RefreshTokenRequest.AsObject;

    static serializeBinaryToWriter(message: RefreshTokenRequest, writer: jspb.BinaryWriter): void;

    static deserializeBinary(bytes: Uint8Array): RefreshTokenRequest;

    static deserializeBinaryFromReader(message: RefreshTokenRequest, reader: jspb.BinaryReader): RefreshTokenRequest;

    getToken(): string;

    setToken(value: string): RefreshTokenRequest;

    serializeBinary(): Uint8Array;

    toObject(includeInstance?: boolean): RefreshTokenRequest.AsObject;
}

export namespace RefreshTokenRequest {
    export type AsObject = {
        token: string,
    }
}
