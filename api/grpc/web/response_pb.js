// source: response.proto
/**
 * @fileoverview
 * @enhanceable
 * @suppress {missingRequire} reports error on implicit type usages.
 * @suppress {messageConventions} JS Compiler reports an error if a variable or
 *     field starts with 'MSG_' and isn't a translatable message.
 * @public
 */
// GENERATED CODE -- DO NOT EDIT!
/* eslint-disable */
// @ts-nocheck

var jspb = require('google-protobuf');
var goog = jspb;
var global = Function('return this')();

var google_protobuf_struct_pb = require('google-protobuf/google/protobuf/struct_pb.js');
goog.object.extend(proto, google_protobuf_struct_pb);
var google_protobuf_timestamp_pb = require('google-protobuf/google/protobuf/timestamp_pb.js');
goog.object.extend(proto, google_protobuf_timestamp_pb);
goog.exportSymbol('proto.gothic.api.BearerResponse', null, global);
goog.exportSymbol('proto.gothic.api.UserResponse', null, global);
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gothic.api.UserResponse = function (opt_data) {
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gothic.api.UserResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
    /**
     * @public
     * @override
     */
    proto.gothic.api.UserResponse.displayName = 'proto.gothic.api.UserResponse';
}
/**
 * Generated by JsPbCodeGenerator.
 * @param {Array=} opt_data Optional initial data array, typically from a
 * server response, or constructed directly in Javascript. The array is used
 * in place and becomes part of the constructed object. It is not cloned.
 * If no data is provided, the constructed object will be empty, but still
 * valid.
 * @extends {jspb.Message}
 * @constructor
 */
proto.gothic.api.BearerResponse = function (opt_data) {
    jspb.Message.initialize(this, opt_data, 0, -1, null, null);
};
goog.inherits(proto.gothic.api.BearerResponse, jspb.Message);
if (goog.DEBUG && !COMPILED) {
    /**
     * @public
     * @override
     */
    proto.gothic.api.BearerResponse.displayName = 'proto.gothic.api.BearerResponse';
}


if (jspb.Message.GENERATE_TO_OBJECT) {
    /**
     * Creates an object representation of this proto.
     * Field names that are reserved in JavaScript and will be renamed to pb_name.
     * Optional fields that are not set will be set to undefined.
     * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
     * For the list of reserved names please see:
     *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
     * @param {boolean=} opt_includeInstance Deprecated. whether to include the
     *     JSPB instance for transitional soy proto support:
     *     http://goto/soy-param-migration
     * @return {!Object}
     */
    proto.gothic.api.UserResponse.prototype.toObject = function (opt_includeInstance) {
        return proto.gothic.api.UserResponse.toObject(opt_includeInstance, this);
    };


    /**
     * Static version of the {@see toObject} method.
     * @param {boolean|undefined} includeInstance Deprecated. Whether to include
     *     the JSPB instance for transitional soy proto support:
     *     http://goto/soy-param-migration
     * @param {!proto.gothic.api.UserResponse} msg The msg instance to transform.
     * @return {!Object}
     * @suppress {unusedLocalVariables} f is only used for nested messages
     */
    proto.gothic.api.UserResponse.toObject = function (includeInstance, msg) {
        var f, obj = {
            role: jspb.Message.getFieldWithDefault(msg, 1, ""),
            email: jspb.Message.getFieldWithDefault(msg, 2, ""),
            username: jspb.Message.getFieldWithDefault(msg, 3, ""),
            data: (f = msg.getData()) && google_protobuf_struct_pb.Struct.toObject(includeInstance, f),
            token: (f = msg.getToken()) && proto.gothic.api.BearerResponse.toObject(includeInstance, f)
        };

        if (includeInstance) {
            obj.$jspbMessageInstance = msg;
        }
        return obj;
    };
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gothic.api.UserResponse}
 */
proto.gothic.api.UserResponse.deserializeBinary = function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.gothic.api.UserResponse;
    return proto.gothic.api.UserResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gothic.api.UserResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gothic.api.UserResponse}
 */
proto.gothic.api.UserResponse.deserializeBinaryFromReader = function (msg, reader) {
    while (reader.nextField()) {
        if (reader.isEndGroup()) {
            break;
        }
        var field = reader.getFieldNumber();
        switch (field) {
            case 1:
                var value = /** @type {string} */ (reader.readString());
                msg.setRole(value);
                break;
            case 2:
                var value = /** @type {string} */ (reader.readString());
                msg.setEmail(value);
                break;
            case 3:
                var value = /** @type {string} */ (reader.readString());
                msg.setUsername(value);
                break;
            case 4:
                var value = new google_protobuf_struct_pb.Struct;
                reader.readMessage(value, google_protobuf_struct_pb.Struct.deserializeBinaryFromReader);
                msg.setData(value);
                break;
            case 5:
                var value = new proto.gothic.api.BearerResponse;
                reader.readMessage(value, proto.gothic.api.BearerResponse.deserializeBinaryFromReader);
                msg.setToken(value);
                break;
            default:
                reader.skipField();
                break;
        }
    }
    return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gothic.api.UserResponse.prototype.serializeBinary = function () {
    var writer = new jspb.BinaryWriter();
    proto.gothic.api.UserResponse.serializeBinaryToWriter(this, writer);
    return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gothic.api.UserResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gothic.api.UserResponse.serializeBinaryToWriter = function (message, writer) {
    var f = undefined;
    f = message.getRole();
    if (f.length > 0) {
        writer.writeString(
            1,
            f
        );
    }
    f = message.getEmail();
    if (f.length > 0) {
        writer.writeString(
            2,
            f
        );
    }
    f = message.getUsername();
    if (f.length > 0) {
        writer.writeString(
            3,
            f
        );
    }
    f = message.getData();
    if (f != null) {
        writer.writeMessage(
            4,
            f,
            google_protobuf_struct_pb.Struct.serializeBinaryToWriter
        );
    }
    f = message.getToken();
    if (f != null) {
        writer.writeMessage(
            5,
            f,
            proto.gothic.api.BearerResponse.serializeBinaryToWriter
        );
    }
};


/**
 * optional string role = 1;
 * @return {string}
 */
proto.gothic.api.UserResponse.prototype.getRole = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.setRole = function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string email = 2;
 * @return {string}
 */
proto.gothic.api.UserResponse.prototype.getEmail = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.setEmail = function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string username = 3;
 * @return {string}
 */
proto.gothic.api.UserResponse.prototype.getUsername = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.setUsername = function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional google.protobuf.Struct data = 4;
 * @return {?proto.google.protobuf.Struct}
 */
proto.gothic.api.UserResponse.prototype.getData = function () {
    return /** @type{?proto.google.protobuf.Struct} */ (
        jspb.Message.getWrapperField(this, google_protobuf_struct_pb.Struct, 4));
};


/**
 * @param {?proto.google.protobuf.Struct|undefined} value
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.setData = function (value) {
    return jspb.Message.setWrapperField(this, 4, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.clearData = function () {
    return this.setData(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.gothic.api.UserResponse.prototype.hasData = function () {
    return jspb.Message.getField(this, 4) != null;
};


/**
 * optional BearerResponse token = 5;
 * @return {?proto.gothic.api.BearerResponse}
 */
proto.gothic.api.UserResponse.prototype.getToken = function () {
    return /** @type{?proto.gothic.api.BearerResponse} */ (
        jspb.Message.getWrapperField(this, proto.gothic.api.BearerResponse, 5));
};


/**
 * @param {?proto.gothic.api.BearerResponse|undefined} value
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.setToken = function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.gothic.api.UserResponse} returns this
 */
proto.gothic.api.UserResponse.prototype.clearToken = function () {
    return this.setToken(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.gothic.api.UserResponse.prototype.hasToken = function () {
    return jspb.Message.getField(this, 5) != null;
};


if (jspb.Message.GENERATE_TO_OBJECT) {
    /**
     * Creates an object representation of this proto.
     * Field names that are reserved in JavaScript and will be renamed to pb_name.
     * Optional fields that are not set will be set to undefined.
     * To access a reserved field use, foo.pb_<name>, eg, foo.pb_default.
     * For the list of reserved names please see:
     *     net/proto2/compiler/js/internal/generator.cc#kKeyword.
     * @param {boolean=} opt_includeInstance Deprecated. whether to include the
     *     JSPB instance for transitional soy proto support:
     *     http://goto/soy-param-migration
     * @return {!Object}
     */
    proto.gothic.api.BearerResponse.prototype.toObject = function (opt_includeInstance) {
        return proto.gothic.api.BearerResponse.toObject(opt_includeInstance, this);
    };


    /**
     * Static version of the {@see toObject} method.
     * @param {boolean|undefined} includeInstance Deprecated. Whether to include
     *     the JSPB instance for transitional soy proto support:
     *     http://goto/soy-param-migration
     * @param {!proto.gothic.api.BearerResponse} msg The msg instance to transform.
     * @return {!Object}
     * @suppress {unusedLocalVariables} f is only used for nested messages
     */
    proto.gothic.api.BearerResponse.toObject = function (includeInstance, msg) {
        var f, obj = {
            type: jspb.Message.getFieldWithDefault(msg, 1, ""),
            access: jspb.Message.getFieldWithDefault(msg, 2, ""),
            refresh: jspb.Message.getFieldWithDefault(msg, 3, ""),
            id: jspb.Message.getFieldWithDefault(msg, 4, ""),
            expiresAt: (f = msg.getExpiresAt()) && google_protobuf_timestamp_pb.Timestamp.toObject(includeInstance, f)
        };

        if (includeInstance) {
            obj.$jspbMessageInstance = msg;
        }
        return obj;
    };
}


/**
 * Deserializes binary data (in protobuf wire format).
 * @param {jspb.ByteSource} bytes The bytes to deserialize.
 * @return {!proto.gothic.api.BearerResponse}
 */
proto.gothic.api.BearerResponse.deserializeBinary = function (bytes) {
    var reader = new jspb.BinaryReader(bytes);
    var msg = new proto.gothic.api.BearerResponse;
    return proto.gothic.api.BearerResponse.deserializeBinaryFromReader(msg, reader);
};


/**
 * Deserializes binary data (in protobuf wire format) from the
 * given reader into the given message object.
 * @param {!proto.gothic.api.BearerResponse} msg The message object to deserialize into.
 * @param {!jspb.BinaryReader} reader The BinaryReader to use.
 * @return {!proto.gothic.api.BearerResponse}
 */
proto.gothic.api.BearerResponse.deserializeBinaryFromReader = function (msg, reader) {
    while (reader.nextField()) {
        if (reader.isEndGroup()) {
            break;
        }
        var field = reader.getFieldNumber();
        switch (field) {
            case 1:
                var value = /** @type {string} */ (reader.readString());
                msg.setType(value);
                break;
            case 2:
                var value = /** @type {string} */ (reader.readString());
                msg.setAccess(value);
                break;
            case 3:
                var value = /** @type {string} */ (reader.readString());
                msg.setRefresh(value);
                break;
            case 4:
                var value = /** @type {string} */ (reader.readString());
                msg.setId(value);
                break;
            case 5:
                var value = new google_protobuf_timestamp_pb.Timestamp;
                reader.readMessage(value, google_protobuf_timestamp_pb.Timestamp.deserializeBinaryFromReader);
                msg.setExpiresAt(value);
                break;
            default:
                reader.skipField();
                break;
        }
    }
    return msg;
};


/**
 * Serializes the message to binary data (in protobuf wire format).
 * @return {!Uint8Array}
 */
proto.gothic.api.BearerResponse.prototype.serializeBinary = function () {
    var writer = new jspb.BinaryWriter();
    proto.gothic.api.BearerResponse.serializeBinaryToWriter(this, writer);
    return writer.getResultBuffer();
};


/**
 * Serializes the given message to binary data (in protobuf wire
 * format), writing to the given BinaryWriter.
 * @param {!proto.gothic.api.BearerResponse} message
 * @param {!jspb.BinaryWriter} writer
 * @suppress {unusedLocalVariables} f is only used for nested messages
 */
proto.gothic.api.BearerResponse.serializeBinaryToWriter = function (message, writer) {
    var f = undefined;
    f = message.getType();
    if (f.length > 0) {
        writer.writeString(
            1,
            f
        );
    }
    f = message.getAccess();
    if (f.length > 0) {
        writer.writeString(
            2,
            f
        );
    }
    f = message.getRefresh();
    if (f.length > 0) {
        writer.writeString(
            3,
            f
        );
    }
    f = message.getId();
    if (f.length > 0) {
        writer.writeString(
            4,
            f
        );
    }
    f = message.getExpiresAt();
    if (f != null) {
        writer.writeMessage(
            5,
            f,
            google_protobuf_timestamp_pb.Timestamp.serializeBinaryToWriter
        );
    }
};


/**
 * optional string type = 1;
 * @return {string}
 */
proto.gothic.api.BearerResponse.prototype.getType = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 1, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.setType = function (value) {
    return jspb.Message.setProto3StringField(this, 1, value);
};


/**
 * optional string access = 2;
 * @return {string}
 */
proto.gothic.api.BearerResponse.prototype.getAccess = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 2, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.setAccess = function (value) {
    return jspb.Message.setProto3StringField(this, 2, value);
};


/**
 * optional string refresh = 3;
 * @return {string}
 */
proto.gothic.api.BearerResponse.prototype.getRefresh = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 3, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.setRefresh = function (value) {
    return jspb.Message.setProto3StringField(this, 3, value);
};


/**
 * optional string id = 4;
 * @return {string}
 */
proto.gothic.api.BearerResponse.prototype.getId = function () {
    return /** @type {string} */ (jspb.Message.getFieldWithDefault(this, 4, ""));
};


/**
 * @param {string} value
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.setId = function (value) {
    return jspb.Message.setProto3StringField(this, 4, value);
};


/**
 * optional google.protobuf.Timestamp expires_at = 5;
 * @return {?proto.google.protobuf.Timestamp}
 */
proto.gothic.api.BearerResponse.prototype.getExpiresAt = function () {
    return /** @type{?proto.google.protobuf.Timestamp} */ (
        jspb.Message.getWrapperField(this, google_protobuf_timestamp_pb.Timestamp, 5));
};


/**
 * @param {?proto.google.protobuf.Timestamp|undefined} value
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.setExpiresAt = function (value) {
    return jspb.Message.setWrapperField(this, 5, value);
};


/**
 * Clears the message field making it undefined.
 * @return {!proto.gothic.api.BearerResponse} returns this
 */
proto.gothic.api.BearerResponse.prototype.clearExpiresAt = function () {
    return this.setExpiresAt(undefined);
};


/**
 * Returns whether this field is set.
 * @return {boolean}
 */
proto.gothic.api.BearerResponse.prototype.hasExpiresAt = function () {
    return jspb.Message.getField(this, 5) != null;
};


goog.object.extend(exports, proto.gothic.api);