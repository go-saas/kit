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

import { V1UpdateMenu } from './v1-update-menu';

/**
 *
 * @export
 * @interface V1UpdateMenuRequest
 */
export interface V1UpdateMenuRequest {
  /**
   *
   * @type {V1UpdateMenu}
   * @memberof V1UpdateMenuRequest
   */
  menu?: V1UpdateMenu;
  /**
   *
   * @type {string}
   * @memberof V1UpdateMenuRequest
   */
  updateMask?: string;
}