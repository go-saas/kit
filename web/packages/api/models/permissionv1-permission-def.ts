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

import { V1PermissionSide } from './v1-permission-side';

/**
 *
 * @export
 * @interface Permissionv1PermissionDef
 */
export interface Permissionv1PermissionDef {
  /**
   *
   * @type {string}
   * @memberof Permissionv1PermissionDef
   */
  displayName?: string;
  /**
   *
   * @type {V1PermissionSide}
   * @memberof Permissionv1PermissionDef
   */
  side?: V1PermissionSide;
  /**
   *
   * @type {object}
   * @memberof Permissionv1PermissionDef
   */
  extra?: object;
  /**
   *
   * @type {string}
   * @memberof Permissionv1PermissionDef
   */
  namespace?: string;
  /**
   *
   * @type {string}
   * @memberof Permissionv1PermissionDef
   */
  action?: string;
  /**
   *
   * @type {boolean}
   * @memberof Permissionv1PermissionDef
   */
  granted?: boolean;
}
