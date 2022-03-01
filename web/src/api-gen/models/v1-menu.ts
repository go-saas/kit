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

import { V1PermissionRequirement } from './v1-permission-requirement';

/**
 *
 * @export
 * @interface V1Menu
 */
export interface V1Menu {
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  id?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  name?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  createdAt?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  desc?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  component?: string;
  /**
   *
   * @type {Array<V1PermissionRequirement>}
   * @memberof V1Menu
   */
  requirement?: Array<V1PermissionRequirement>;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  parent?: string;
  /**
   *
   * @type {object}
   * @memberof V1Menu
   */
  props?: object;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  fullPath?: string;
  /**
   *
   * @type {number}
   * @memberof V1Menu
   */
  priority?: number;
  /**
   *
   * @type {boolean}
   * @memberof V1Menu
   */
  ignoreAuth?: boolean;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  icon?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  iframe?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  microApp?: string;
  /**
   *
   * @type {object}
   * @memberof V1Menu
   */
  meta?: object;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  title?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  path?: string;
  /**
   *
   * @type {string}
   * @memberof V1Menu
   */
  redirect?: string;
}
