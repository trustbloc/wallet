/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { CHAPIEventHandler } from './chapi';
import { normalizeQuery } from '../common/helper';
import jp from 'jsonpath';

/**
 * WACIRedirectHandler handles WACI events for redirect flows.
 *
 * currently handling only WACI share flow with redirect.
 *
 * TODO: to be implemented
 *
 * @param credEvent instance
 * @class
 */
export class WACIRedirectHandler {
  constructor(credEvent) {
    this.credEvent = credEvent;
  }

  message() {
    // oob message from waci request
    return null;
  }

  requestor() {
    return null;
  }

  done() {
    // redirect to `redirectURL` with status success
  }

  cancel() {
    // redirect to `redirectURL` with status error
  }
}

/**
 * WACIPolyfillHandler handles CHAPI polyfill based WACI flows.
 *
 * @param credEvent instance
 * @class
 */
export class WACIPolyfillHandler {
  constructor(credEvent) {
    this.handler = new CHAPIEventHandler(credEvent);
  }

  message() {
    let { query } = this.handler.getEventData();

    let oob = jp.query(normalizeQuery(query), '$[?(@.type=="WACIShare")].credentialQuery[*].oob');
    if (oob.length > 1) {
      throw 'invalid request';
    }

    return oob[0];
  }

  requestor() {
    return this.handler.requestor();
  }

  done() {
    return this.handler.done();
  }

  cancel() {
    return this.handler.cancel();
  }
}
