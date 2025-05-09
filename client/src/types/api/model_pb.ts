// @generated by protoc-gen-es v2.2.5 with parameter "target=ts"
// @generated from file api/model.proto (package model, syntax proto3)
/* eslint-disable */

import type { GenFile, GenMessage, GenService } from "@bufbuild/protobuf/codegenv1";
import { fileDesc, messageDesc, serviceDesc } from "@bufbuild/protobuf/codegenv1";
import type { Timestamp } from "@bufbuild/protobuf/wkt";
import { file_google_protobuf_timestamp } from "@bufbuild/protobuf/wkt";
import type { Message } from "@bufbuild/protobuf";

/**
 * Describes the file api/model.proto.
 */
export const file_api_model: GenFile = /*@__PURE__*/
  fileDesc("Cg9hcGkvbW9kZWwucHJvdG8SBW1vZGVsIjcKIEdldFJlZ2lzdHJ5VG9rZW5QYXNzd29yZHNSZXF1ZXN0EhMKC3dvcmtzcGFjZUlkGAEgASgJIqcBCiFHZXRSZWdpc3RyeVRva2VuUGFzc3dvcmRzUmVzcG9uc2USMgoJcGFzc3dvcmQxGAEgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcEgAiAEBEjIKCXBhc3N3b3JkMhgCIAEoCzIaLmdvb2dsZS5wcm90b2J1Zi5UaW1lc3RhbXBIAYgBAUIMCgpfcGFzc3dvcmQxQgwKCl9wYXNzd29yZDIiTAoiQ3JlYXRlUmVnaXN0cnlUb2tlblBhc3N3b3JkUmVxdWVzdBITCgt3b3Jrc3BhY2VJZBgBIAEoCRIRCglwYXNzd29yZDIYAiABKAgiZgojQ3JlYXRlUmVnaXN0cnlUb2tlblBhc3N3b3JkUmVzcG9uc2USEAoIcGFzc3dvcmQYASABKAkSLQoJY3JlYXRlZEF0GAIgASgLMhouZ29vZ2xlLnByb3RvYnVmLlRpbWVzdGFtcCInChBHZXRJbWFnZXNSZXF1ZXN0EhMKC3dvcmtzcGFjZUlkGAEgASgJIjEKEUdldEltYWdlc1Jlc3BvbnNlEhwKBmltYWdlcxgBIAMoCzIMLm1vZGVsLkltYWdlIqYBChJDcmVhdGVNb2RlbFJlcXVlc3QSHAoUaW5wdXRTcGVjaWZpY2F0aW9uSWQYASABKAkSHQoVb3V0cHV0U3BlY2lmaWNhdGlvbklkGAIgASgJEgwKBG5hbWUYAyABKAkSHQoQcGFyYW1ldGVyc1NjaGVtYRgEIAEoCUgAiAEBEhEKCWltYWdlTmFtZRgFIAEoCUITChFfcGFyYW1ldGVyc1NjaGVtYSIhChNDcmVhdGVNb2RlbFJlc3BvbnNlEgoKAmlkGAEgASgJIicKEEdldE1vZGVsc1JlcXVlc3QSEwoLd29ya3NwYWNlSWQYASABKAkiNQoRR2V0TW9kZWxzUmVzcG9uc2USIAoGbW9kZWxzGAEgAygLMhAubW9kZWwuTW9kZWxOYW1lIh0KD0dldE1vZGVsUmVxdWVzdBIKCgJpZBgBIAEoCSIvChBHZXRNb2RlbFJlc3BvbnNlEhsKBW1vZGVsGAEgASgLMgwubW9kZWwuTW9kZWwieAoVQ3JlYXRlTW9kZWxSdW5SZXF1ZXN0Eg8KB21vZGVsSWQYASABKAkSGAoQaW5wdXRGaWxlR3JvdXBJZBgCIAEoCRIXCgpwYXJhbWV0ZXJzGAMgASgJSACIAQESDAoEbmFtZRgEIAEoCUINCgtfcGFyYW1ldGVycyIkChZDcmVhdGVNb2RlbFJ1blJlc3BvbnNlEgoKAmlkGAEgASgJIioKE0dldE1vZGVsUnVuc1JlcXVlc3QSEwoLd29ya3NwYWNlSWQYASABKAkiOgoUR2V0TW9kZWxSdW5zUmVzcG9uc2USIgoJbW9kZWxSdW5zGAEgAygLMg8ubW9kZWwuTW9kZWxSdW4iJAoFSW1hZ2USDQoFc2NvcGUYASABKAkSDAoEbmFtZRgCIAEoCSKFAQoITW9kZWxSdW4SCgoCaWQYASABKAkSDwoHbW9kZWxJZBgCIAEoCRIYChBpbnB1dEZpbGVHcm91cElkGAMgASgJEh4KEW91dHB1dEZpbGVHcm91cElkGAQgASgJSACIAQESDAoEbmFtZRgFIAEoCUIUChJfb3V0cHV0RmlsZUdyb3VwSWQipQEKBU1vZGVsEgoKAmlkGAEgASgJEgwKBG5hbWUYAiABKAkSHAoUaW5wdXRTcGVjaWZpY2F0aW9uSWQYAyABKAkSHQoVb3V0cHV0U3BlY2lmaWNhdGlvbklkGAQgASgJEh0KEHBhcmFtZXRlcnNTY2hlbWEYBSABKAlIAIgBARIRCglpbWFnZU5hbWUYBiABKAlCEwoRX3BhcmFtZXRlcnNTY2hlbWEiJQoJTW9kZWxOYW1lEgoKAmlkGAEgASgJEgwKBG5hbWUYAiABKAkyjwUKDE1vZGVsU2VydmljZRJuChlHZXRSZWdpc3RyeVRva2VuUGFzc3dvcmRzEicubW9kZWwuR2V0UmVnaXN0cnlUb2tlblBhc3N3b3Jkc1JlcXVlc3QaKC5tb2RlbC5HZXRSZWdpc3RyeVRva2VuUGFzc3dvcmRzUmVzcG9uc2USdAobQ3JlYXRlUmVnaXN0cnlUb2tlblBhc3N3b3JkEikubW9kZWwuQ3JlYXRlUmVnaXN0cnlUb2tlblBhc3N3b3JkUmVxdWVzdBoqLm1vZGVsLkNyZWF0ZVJlZ2lzdHJ5VG9rZW5QYXNzd29yZFJlc3BvbnNlEj4KCUdldEltYWdlcxIXLm1vZGVsLkdldEltYWdlc1JlcXVlc3QaGC5tb2RlbC5HZXRJbWFnZXNSZXNwb25zZRJECgtDcmVhdGVNb2RlbBIZLm1vZGVsLkNyZWF0ZU1vZGVsUmVxdWVzdBoaLm1vZGVsLkNyZWF0ZU1vZGVsUmVzcG9uc2USPgoJR2V0TW9kZWxzEhcubW9kZWwuR2V0TW9kZWxzUmVxdWVzdBoYLm1vZGVsLkdldE1vZGVsc1Jlc3BvbnNlEjsKCEdldE1vZGVsEhYubW9kZWwuR2V0TW9kZWxSZXF1ZXN0GhcubW9kZWwuR2V0TW9kZWxSZXNwb25zZRJNCg5DcmVhdGVNb2RlbFJ1bhIcLm1vZGVsLkNyZWF0ZU1vZGVsUnVuUmVxdWVzdBodLm1vZGVsLkNyZWF0ZU1vZGVsUnVuUmVzcG9uc2USRwoMR2V0TW9kZWxSdW5zEhoubW9kZWwuR2V0TW9kZWxSdW5zUmVxdWVzdBobLm1vZGVsLkdldE1vZGVsUnVuc1Jlc3BvbnNlQilaJ2dpdGh1Yi5jb20vTW9yaGFmQWxzaGlibHkvaXVudmkvZ2VuL2FwaWIGcHJvdG8z", [file_google_protobuf_timestamp]);

/**
 * @generated from message model.GetRegistryTokenPasswordsRequest
 */
export type GetRegistryTokenPasswordsRequest = Message<"model.GetRegistryTokenPasswordsRequest"> & {
  /**
   * @generated from field: string workspaceId = 1;
   */
  workspaceId: string;
};

/**
 * Describes the message model.GetRegistryTokenPasswordsRequest.
 * Use `create(GetRegistryTokenPasswordsRequestSchema)` to create a new message.
 */
export const GetRegistryTokenPasswordsRequestSchema: GenMessage<GetRegistryTokenPasswordsRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 0);

/**
 * @generated from message model.GetRegistryTokenPasswordsResponse
 */
export type GetRegistryTokenPasswordsResponse = Message<"model.GetRegistryTokenPasswordsResponse"> & {
  /**
   * @generated from field: optional google.protobuf.Timestamp password1 = 1;
   */
  password1?: Timestamp;

  /**
   * @generated from field: optional google.protobuf.Timestamp password2 = 2;
   */
  password2?: Timestamp;
};

/**
 * Describes the message model.GetRegistryTokenPasswordsResponse.
 * Use `create(GetRegistryTokenPasswordsResponseSchema)` to create a new message.
 */
export const GetRegistryTokenPasswordsResponseSchema: GenMessage<GetRegistryTokenPasswordsResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 1);

/**
 * @generated from message model.CreateRegistryTokenPasswordRequest
 */
export type CreateRegistryTokenPasswordRequest = Message<"model.CreateRegistryTokenPasswordRequest"> & {
  /**
   * @generated from field: string workspaceId = 1;
   */
  workspaceId: string;

  /**
   * @generated from field: bool password2 = 2;
   */
  password2: boolean;
};

/**
 * Describes the message model.CreateRegistryTokenPasswordRequest.
 * Use `create(CreateRegistryTokenPasswordRequestSchema)` to create a new message.
 */
export const CreateRegistryTokenPasswordRequestSchema: GenMessage<CreateRegistryTokenPasswordRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 2);

/**
 * @generated from message model.CreateRegistryTokenPasswordResponse
 */
export type CreateRegistryTokenPasswordResponse = Message<"model.CreateRegistryTokenPasswordResponse"> & {
  /**
   * @generated from field: string password = 1;
   */
  password: string;

  /**
   * @generated from field: google.protobuf.Timestamp createdAt = 2;
   */
  createdAt?: Timestamp;
};

/**
 * Describes the message model.CreateRegistryTokenPasswordResponse.
 * Use `create(CreateRegistryTokenPasswordResponseSchema)` to create a new message.
 */
export const CreateRegistryTokenPasswordResponseSchema: GenMessage<CreateRegistryTokenPasswordResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 3);

/**
 * @generated from message model.GetImagesRequest
 */
export type GetImagesRequest = Message<"model.GetImagesRequest"> & {
  /**
   * @generated from field: string workspaceId = 1;
   */
  workspaceId: string;
};

/**
 * Describes the message model.GetImagesRequest.
 * Use `create(GetImagesRequestSchema)` to create a new message.
 */
export const GetImagesRequestSchema: GenMessage<GetImagesRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 4);

/**
 * @generated from message model.GetImagesResponse
 */
export type GetImagesResponse = Message<"model.GetImagesResponse"> & {
  /**
   * @generated from field: repeated model.Image images = 1;
   */
  images: Image[];
};

/**
 * Describes the message model.GetImagesResponse.
 * Use `create(GetImagesResponseSchema)` to create a new message.
 */
export const GetImagesResponseSchema: GenMessage<GetImagesResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 5);

/**
 * @generated from message model.CreateModelRequest
 */
export type CreateModelRequest = Message<"model.CreateModelRequest"> & {
  /**
   * @generated from field: string inputSpecificationId = 1;
   */
  inputSpecificationId: string;

  /**
   * @generated from field: string outputSpecificationId = 2;
   */
  outputSpecificationId: string;

  /**
   * @generated from field: string name = 3;
   */
  name: string;

  /**
   * @generated from field: optional string parametersSchema = 4;
   */
  parametersSchema?: string;

  /**
   * @generated from field: string imageName = 5;
   */
  imageName: string;
};

/**
 * Describes the message model.CreateModelRequest.
 * Use `create(CreateModelRequestSchema)` to create a new message.
 */
export const CreateModelRequestSchema: GenMessage<CreateModelRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 6);

/**
 * @generated from message model.CreateModelResponse
 */
export type CreateModelResponse = Message<"model.CreateModelResponse"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message model.CreateModelResponse.
 * Use `create(CreateModelResponseSchema)` to create a new message.
 */
export const CreateModelResponseSchema: GenMessage<CreateModelResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 7);

/**
 * @generated from message model.GetModelsRequest
 */
export type GetModelsRequest = Message<"model.GetModelsRequest"> & {
  /**
   * @generated from field: string workspaceId = 1;
   */
  workspaceId: string;
};

/**
 * Describes the message model.GetModelsRequest.
 * Use `create(GetModelsRequestSchema)` to create a new message.
 */
export const GetModelsRequestSchema: GenMessage<GetModelsRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 8);

/**
 * @generated from message model.GetModelsResponse
 */
export type GetModelsResponse = Message<"model.GetModelsResponse"> & {
  /**
   * @generated from field: repeated model.ModelName models = 1;
   */
  models: ModelName[];
};

/**
 * Describes the message model.GetModelsResponse.
 * Use `create(GetModelsResponseSchema)` to create a new message.
 */
export const GetModelsResponseSchema: GenMessage<GetModelsResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 9);

/**
 * @generated from message model.GetModelRequest
 */
export type GetModelRequest = Message<"model.GetModelRequest"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message model.GetModelRequest.
 * Use `create(GetModelRequestSchema)` to create a new message.
 */
export const GetModelRequestSchema: GenMessage<GetModelRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 10);

/**
 * @generated from message model.GetModelResponse
 */
export type GetModelResponse = Message<"model.GetModelResponse"> & {
  /**
   * @generated from field: model.Model model = 1;
   */
  model?: Model;
};

/**
 * Describes the message model.GetModelResponse.
 * Use `create(GetModelResponseSchema)` to create a new message.
 */
export const GetModelResponseSchema: GenMessage<GetModelResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 11);

/**
 * @generated from message model.CreateModelRunRequest
 */
export type CreateModelRunRequest = Message<"model.CreateModelRunRequest"> & {
  /**
   * @generated from field: string modelId = 1;
   */
  modelId: string;

  /**
   * @generated from field: string inputFileGroupId = 2;
   */
  inputFileGroupId: string;

  /**
   * @generated from field: optional string parameters = 3;
   */
  parameters?: string;

  /**
   * @generated from field: string name = 4;
   */
  name: string;
};

/**
 * Describes the message model.CreateModelRunRequest.
 * Use `create(CreateModelRunRequestSchema)` to create a new message.
 */
export const CreateModelRunRequestSchema: GenMessage<CreateModelRunRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 12);

/**
 * @generated from message model.CreateModelRunResponse
 */
export type CreateModelRunResponse = Message<"model.CreateModelRunResponse"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;
};

/**
 * Describes the message model.CreateModelRunResponse.
 * Use `create(CreateModelRunResponseSchema)` to create a new message.
 */
export const CreateModelRunResponseSchema: GenMessage<CreateModelRunResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 13);

/**
 * @generated from message model.GetModelRunsRequest
 */
export type GetModelRunsRequest = Message<"model.GetModelRunsRequest"> & {
  /**
   * @generated from field: string workspaceId = 1;
   */
  workspaceId: string;
};

/**
 * Describes the message model.GetModelRunsRequest.
 * Use `create(GetModelRunsRequestSchema)` to create a new message.
 */
export const GetModelRunsRequestSchema: GenMessage<GetModelRunsRequest> = /*@__PURE__*/
  messageDesc(file_api_model, 14);

/**
 * @generated from message model.GetModelRunsResponse
 */
export type GetModelRunsResponse = Message<"model.GetModelRunsResponse"> & {
  /**
   * @generated from field: repeated model.ModelRun modelRuns = 1;
   */
  modelRuns: ModelRun[];
};

/**
 * Describes the message model.GetModelRunsResponse.
 * Use `create(GetModelRunsResponseSchema)` to create a new message.
 */
export const GetModelRunsResponseSchema: GenMessage<GetModelRunsResponse> = /*@__PURE__*/
  messageDesc(file_api_model, 15);

/**
 * @generated from message model.Image
 */
export type Image = Message<"model.Image"> & {
  /**
   * @generated from field: string scope = 1;
   */
  scope: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;
};

/**
 * Describes the message model.Image.
 * Use `create(ImageSchema)` to create a new message.
 */
export const ImageSchema: GenMessage<Image> = /*@__PURE__*/
  messageDesc(file_api_model, 16);

/**
 * @generated from message model.ModelRun
 */
export type ModelRun = Message<"model.ModelRun"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string modelId = 2;
   */
  modelId: string;

  /**
   * @generated from field: string inputFileGroupId = 3;
   */
  inputFileGroupId: string;

  /**
   * @generated from field: optional string outputFileGroupId = 4;
   */
  outputFileGroupId?: string;

  /**
   * @generated from field: string name = 5;
   */
  name: string;
};

/**
 * Describes the message model.ModelRun.
 * Use `create(ModelRunSchema)` to create a new message.
 */
export const ModelRunSchema: GenMessage<ModelRun> = /*@__PURE__*/
  messageDesc(file_api_model, 17);

/**
 * @generated from message model.Model
 */
export type Model = Message<"model.Model"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;

  /**
   * @generated from field: string inputSpecificationId = 3;
   */
  inputSpecificationId: string;

  /**
   * @generated from field: string outputSpecificationId = 4;
   */
  outputSpecificationId: string;

  /**
   * @generated from field: optional string parametersSchema = 5;
   */
  parametersSchema?: string;

  /**
   * @generated from field: string imageName = 6;
   */
  imageName: string;
};

/**
 * Describes the message model.Model.
 * Use `create(ModelSchema)` to create a new message.
 */
export const ModelSchema: GenMessage<Model> = /*@__PURE__*/
  messageDesc(file_api_model, 18);

/**
 * @generated from message model.ModelName
 */
export type ModelName = Message<"model.ModelName"> & {
  /**
   * @generated from field: string id = 1;
   */
  id: string;

  /**
   * @generated from field: string name = 2;
   */
  name: string;
};

/**
 * Describes the message model.ModelName.
 * Use `create(ModelNameSchema)` to create a new message.
 */
export const ModelNameSchema: GenMessage<ModelName> = /*@__PURE__*/
  messageDesc(file_api_model, 19);

/**
 * @generated from service model.ModelService
 */
export const ModelService: GenService<{
  /**
   * @generated from rpc model.ModelService.GetRegistryTokenPasswords
   */
  getRegistryTokenPasswords: {
    methodKind: "unary";
    input: typeof GetRegistryTokenPasswordsRequestSchema;
    output: typeof GetRegistryTokenPasswordsResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.CreateRegistryTokenPassword
   */
  createRegistryTokenPassword: {
    methodKind: "unary";
    input: typeof CreateRegistryTokenPasswordRequestSchema;
    output: typeof CreateRegistryTokenPasswordResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.GetImages
   */
  getImages: {
    methodKind: "unary";
    input: typeof GetImagesRequestSchema;
    output: typeof GetImagesResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.CreateModel
   */
  createModel: {
    methodKind: "unary";
    input: typeof CreateModelRequestSchema;
    output: typeof CreateModelResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.GetModels
   */
  getModels: {
    methodKind: "unary";
    input: typeof GetModelsRequestSchema;
    output: typeof GetModelsResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.GetModel
   */
  getModel: {
    methodKind: "unary";
    input: typeof GetModelRequestSchema;
    output: typeof GetModelResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.CreateModelRun
   */
  createModelRun: {
    methodKind: "unary";
    input: typeof CreateModelRunRequestSchema;
    output: typeof CreateModelRunResponseSchema;
  },
  /**
   * @generated from rpc model.ModelService.GetModelRuns
   */
  getModelRuns: {
    methodKind: "unary";
    input: typeof GetModelRunsRequestSchema;
    output: typeof GetModelRunsResponseSchema;
  },
}> = /*@__PURE__*/
  serviceDesc(file_api_model, 0);

