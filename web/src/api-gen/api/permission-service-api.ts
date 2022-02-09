/* tslint:disable */
/* eslint-disable */
/**
 * Saas Service
 * No description provided (generated by Openapi Generator https://github.com/openapitools/openapi-generator)
 *
 * The version of the OpenAPI document: 1.0
 *
 *
 * NOTE: This class is auto generated by OpenAPI Generator (https://openapi-generator.tech).
 * https://openapi-generator.tech
 * Do not edit the class manually.
 */

import globalAxios, { AxiosPromise, AxiosInstance, AxiosRequestConfig } from 'axios';
import { Configuration } from '../configuration';
// Some imports not used depending on template conditions
// @ts-ignore
import {
  DUMMY_BASE_URL,
  assertParamExists,
  setApiKeyToObject,
  setBasicAuthToObject,
  setBearerAuthToObject,
  setOAuthToObject,
  setSearchParams,
  serializeDataIfNeeded,
  toPathString,
  createRequestFunction,
} from '../common';
// @ts-ignore
import { BASE_PATH, COLLECTION_FORMATS, RequestArgs, BaseAPI, RequiredError } from '../base';
// @ts-ignore
import { RpcStatus } from '../models';
// @ts-ignore
import { V1CheckPermissionReply } from '../models';
// @ts-ignore
import { V1CheckPermissionRequest } from '../models';
// @ts-ignore
import { V1CheckSubjectsPermissionReply } from '../models';
// @ts-ignore
import { V1CheckSubjectsPermissionRequest } from '../models';
// @ts-ignore
import { V1GetCurrentPermissionReply } from '../models';
// @ts-ignore
import { V1UpdateSubjectPermissionRequest } from '../models';
// @ts-ignore
import { V1UpdateSubjectPermissionResponse } from '../models';
/**
 * PermissionServiceApi - axios parameter creator
 * @export
 */
export const PermissionServiceApiAxiosParamCreator = function (configuration?: Configuration) {
  return {
    /**
     *
     * @param {V1CheckPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceCheckCurrent: async (
      body: V1CheckPermissionRequest,
      options: AxiosRequestConfig = {},
    ): Promise<RequestArgs> => {
      // verify required parameter 'body' is not null or undefined
      assertParamExists('permissionServiceCheckCurrent', 'body', body);
      const localVarPath = `/v1/permission/check`;
      // use dummy base URL string because the URL constructor only accepts absolute URLs.
      const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
      let baseOptions;
      if (configuration) {
        baseOptions = configuration.baseOptions;
      }

      const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options };
      const localVarHeaderParameter = {} as any;
      const localVarQueryParameter = {} as any;

      localVarHeaderParameter['Content-Type'] = 'application/json';

      setSearchParams(localVarUrlObj, localVarQueryParameter);
      let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
      localVarRequestOptions.headers = {
        ...localVarHeaderParameter,
        ...headersFromBaseOptions,
        ...options.headers,
      };
      localVarRequestOptions.data = serializeDataIfNeeded(
        body,
        localVarRequestOptions,
        configuration,
      );

      return {
        url: toPathString(localVarUrlObj),
        options: localVarRequestOptions,
      };
    },
    /**
     *
     * @param {V1CheckSubjectsPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceCheckForSubjects: async (
      body: V1CheckSubjectsPermissionRequest,
      options: AxiosRequestConfig = {},
    ): Promise<RequestArgs> => {
      // verify required parameter 'body' is not null or undefined
      assertParamExists('permissionServiceCheckForSubjects', 'body', body);
      const localVarPath = `/v1/permission/check-subjects`;
      // use dummy base URL string because the URL constructor only accepts absolute URLs.
      const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
      let baseOptions;
      if (configuration) {
        baseOptions = configuration.baseOptions;
      }

      const localVarRequestOptions = { method: 'POST', ...baseOptions, ...options };
      const localVarHeaderParameter = {} as any;
      const localVarQueryParameter = {} as any;

      localVarHeaderParameter['Content-Type'] = 'application/json';

      setSearchParams(localVarUrlObj, localVarQueryParameter);
      let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
      localVarRequestOptions.headers = {
        ...localVarHeaderParameter,
        ...headersFromBaseOptions,
        ...options.headers,
      };
      localVarRequestOptions.data = serializeDataIfNeeded(
        body,
        localVarRequestOptions,
        configuration,
      );

      return {
        url: toPathString(localVarUrlObj),
        options: localVarRequestOptions,
      };
    },
    /**
     *
     * @summary Get current permission
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceGetCurrent: async (options: AxiosRequestConfig = {}): Promise<RequestArgs> => {
      const localVarPath = `/v1/permission/current`;
      // use dummy base URL string because the URL constructor only accepts absolute URLs.
      const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
      let baseOptions;
      if (configuration) {
        baseOptions = configuration.baseOptions;
      }

      const localVarRequestOptions = { method: 'GET', ...baseOptions, ...options };
      const localVarHeaderParameter = {} as any;
      const localVarQueryParameter = {} as any;

      setSearchParams(localVarUrlObj, localVarQueryParameter);
      let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
      localVarRequestOptions.headers = {
        ...localVarHeaderParameter,
        ...headersFromBaseOptions,
        ...options.headers,
      };

      return {
        url: toPathString(localVarUrlObj),
        options: localVarRequestOptions,
      };
    },
    /**
     *
     * @param {V1UpdateSubjectPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceUpdateSubjectPermission: async (
      body: V1UpdateSubjectPermissionRequest,
      options: AxiosRequestConfig = {},
    ): Promise<RequestArgs> => {
      // verify required parameter 'body' is not null or undefined
      assertParamExists('permissionServiceUpdateSubjectPermission', 'body', body);
      const localVarPath = `/v1/permission/subject`;
      // use dummy base URL string because the URL constructor only accepts absolute URLs.
      const localVarUrlObj = new URL(localVarPath, DUMMY_BASE_URL);
      let baseOptions;
      if (configuration) {
        baseOptions = configuration.baseOptions;
      }

      const localVarRequestOptions = { method: 'PUT', ...baseOptions, ...options };
      const localVarHeaderParameter = {} as any;
      const localVarQueryParameter = {} as any;

      localVarHeaderParameter['Content-Type'] = 'application/json';

      setSearchParams(localVarUrlObj, localVarQueryParameter);
      let headersFromBaseOptions = baseOptions && baseOptions.headers ? baseOptions.headers : {};
      localVarRequestOptions.headers = {
        ...localVarHeaderParameter,
        ...headersFromBaseOptions,
        ...options.headers,
      };
      localVarRequestOptions.data = serializeDataIfNeeded(
        body,
        localVarRequestOptions,
        configuration,
      );

      return {
        url: toPathString(localVarUrlObj),
        options: localVarRequestOptions,
      };
    },
  };
};

/**
 * PermissionServiceApi - functional programming interface
 * @export
 */
export const PermissionServiceApiFp = function (configuration?: Configuration) {
  const localVarAxiosParamCreator = PermissionServiceApiAxiosParamCreator(configuration);
  return {
    /**
     *
     * @param {V1CheckPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    async permissionServiceCheckCurrent(
      body: V1CheckPermissionRequest,
      options?: AxiosRequestConfig,
    ): Promise<(axios?: AxiosInstance, basePath?: string) => AxiosPromise<V1CheckPermissionReply>> {
      const localVarAxiosArgs = await localVarAxiosParamCreator.permissionServiceCheckCurrent(
        body,
        options,
      );
      return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
    },
    /**
     *
     * @param {V1CheckSubjectsPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    async permissionServiceCheckForSubjects(
      body: V1CheckSubjectsPermissionRequest,
      options?: AxiosRequestConfig,
    ): Promise<
      (axios?: AxiosInstance, basePath?: string) => AxiosPromise<V1CheckSubjectsPermissionReply>
    > {
      const localVarAxiosArgs = await localVarAxiosParamCreator.permissionServiceCheckForSubjects(
        body,
        options,
      );
      return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
    },
    /**
     *
     * @summary Get current permission
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    async permissionServiceGetCurrent(
      options?: AxiosRequestConfig,
    ): Promise<
      (axios?: AxiosInstance, basePath?: string) => AxiosPromise<V1GetCurrentPermissionReply>
    > {
      const localVarAxiosArgs = await localVarAxiosParamCreator.permissionServiceGetCurrent(
        options,
      );
      return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
    },
    /**
     *
     * @param {V1UpdateSubjectPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    async permissionServiceUpdateSubjectPermission(
      body: V1UpdateSubjectPermissionRequest,
      options?: AxiosRequestConfig,
    ): Promise<
      (axios?: AxiosInstance, basePath?: string) => AxiosPromise<V1UpdateSubjectPermissionResponse>
    > {
      const localVarAxiosArgs =
        await localVarAxiosParamCreator.permissionServiceUpdateSubjectPermission(body, options);
      return createRequestFunction(localVarAxiosArgs, globalAxios, BASE_PATH, configuration);
    },
  };
};

/**
 * PermissionServiceApi - factory interface
 * @export
 */
export const PermissionServiceApiFactory = function (
  configuration?: Configuration,
  basePath?: string,
  axios?: AxiosInstance,
) {
  const localVarFp = PermissionServiceApiFp(configuration);
  return {
    /**
     *
     * @param {V1CheckPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceCheckCurrent(
      body: V1CheckPermissionRequest,
      options?: any,
    ): AxiosPromise<V1CheckPermissionReply> {
      return localVarFp
        .permissionServiceCheckCurrent(body, options)
        .then((request) => request(axios, basePath));
    },
    /**
     *
     * @param {V1CheckSubjectsPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceCheckForSubjects(
      body: V1CheckSubjectsPermissionRequest,
      options?: any,
    ): AxiosPromise<V1CheckSubjectsPermissionReply> {
      return localVarFp
        .permissionServiceCheckForSubjects(body, options)
        .then((request) => request(axios, basePath));
    },
    /**
     *
     * @summary Get current permission
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceGetCurrent(options?: any): AxiosPromise<V1GetCurrentPermissionReply> {
      return localVarFp
        .permissionServiceGetCurrent(options)
        .then((request) => request(axios, basePath));
    },
    /**
     *
     * @param {V1UpdateSubjectPermissionRequest} body
     * @param {*} [options] Override http request option.
     * @throws {RequiredError}
     */
    permissionServiceUpdateSubjectPermission(
      body: V1UpdateSubjectPermissionRequest,
      options?: any,
    ): AxiosPromise<V1UpdateSubjectPermissionResponse> {
      return localVarFp
        .permissionServiceUpdateSubjectPermission(body, options)
        .then((request) => request(axios, basePath));
    },
  };
};

/**
 * Request parameters for permissionServiceCheckCurrent operation in PermissionServiceApi.
 * @export
 * @interface PermissionServiceApiPermissionServiceCheckCurrentRequest
 */
export interface PermissionServiceApiPermissionServiceCheckCurrentRequest {
  /**
   *
   * @type {V1CheckPermissionRequest}
   * @memberof PermissionServiceApiPermissionServiceCheckCurrent
   */
  readonly body: V1CheckPermissionRequest;
}

/**
 * Request parameters for permissionServiceCheckForSubjects operation in PermissionServiceApi.
 * @export
 * @interface PermissionServiceApiPermissionServiceCheckForSubjectsRequest
 */
export interface PermissionServiceApiPermissionServiceCheckForSubjectsRequest {
  /**
   *
   * @type {V1CheckSubjectsPermissionRequest}
   * @memberof PermissionServiceApiPermissionServiceCheckForSubjects
   */
  readonly body: V1CheckSubjectsPermissionRequest;
}

/**
 * Request parameters for permissionServiceUpdateSubjectPermission operation in PermissionServiceApi.
 * @export
 * @interface PermissionServiceApiPermissionServiceUpdateSubjectPermissionRequest
 */
export interface PermissionServiceApiPermissionServiceUpdateSubjectPermissionRequest {
  /**
   *
   * @type {V1UpdateSubjectPermissionRequest}
   * @memberof PermissionServiceApiPermissionServiceUpdateSubjectPermission
   */
  readonly body: V1UpdateSubjectPermissionRequest;
}

/**
 * PermissionServiceApi - object-oriented interface
 * @export
 * @class PermissionServiceApi
 * @extends {BaseAPI}
 */
export class PermissionServiceApi extends BaseAPI {
  /**
   *
   * @param {PermissionServiceApiPermissionServiceCheckCurrentRequest} requestParameters Request parameters.
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof PermissionServiceApi
   */
  public permissionServiceCheckCurrent(
    requestParameters: PermissionServiceApiPermissionServiceCheckCurrentRequest,
    options?: AxiosRequestConfig,
  ) {
    return PermissionServiceApiFp(this.configuration)
      .permissionServiceCheckCurrent(requestParameters.body, options)
      .then((request) => request(this.axios, this.basePath));
  }

  /**
   *
   * @param {PermissionServiceApiPermissionServiceCheckForSubjectsRequest} requestParameters Request parameters.
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof PermissionServiceApi
   */
  public permissionServiceCheckForSubjects(
    requestParameters: PermissionServiceApiPermissionServiceCheckForSubjectsRequest,
    options?: AxiosRequestConfig,
  ) {
    return PermissionServiceApiFp(this.configuration)
      .permissionServiceCheckForSubjects(requestParameters.body, options)
      .then((request) => request(this.axios, this.basePath));
  }

  /**
   *
   * @summary Get current permission
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof PermissionServiceApi
   */
  public permissionServiceGetCurrent(options?: AxiosRequestConfig) {
    return PermissionServiceApiFp(this.configuration)
      .permissionServiceGetCurrent(options)
      .then((request) => request(this.axios, this.basePath));
  }

  /**
   *
   * @param {PermissionServiceApiPermissionServiceUpdateSubjectPermissionRequest} requestParameters Request parameters.
   * @param {*} [options] Override http request option.
   * @throws {RequiredError}
   * @memberof PermissionServiceApi
   */
  public permissionServiceUpdateSubjectPermission(
    requestParameters: PermissionServiceApiPermissionServiceUpdateSubjectPermissionRequest,
    options?: AxiosRequestConfig,
  ) {
    return PermissionServiceApiFp(this.configuration)
      .permissionServiceUpdateSubjectPermission(requestParameters.body, options)
      .then((request) => request(this.axios, this.basePath));
  }
}