/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { CHAPIEventHandler } from './chapi';
import { normalizeQuery, wait } from '../common/helper';
import jp from 'jsonpath';

/**
 * WACIRedirectHandler handles WACI events for redirect flows.
 *
 * currently handling only WACI share flow with redirect.
 *
 * @param credEvent instance
 * @class
 */
export class WACIRedirectHandler {
  constructor(oob) {
    this.oob = oob;
  }

  message() {
    // oob message from waci request
    return this.oob;
  }

  requestor() {
    //TODO to be removed, requestor info currently unavailable in WACI
    return 'requestor';
  }

  done(url = '/') {
    // redirect to `url` with status success
    window.location.href = url;
  }

  // TODO [Issue #1325] for standard error screen or initiate problem report flow.
  cancel() {
    // redirect to `redirectURL` with status error
    window.location.href = '/';
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

  // TODO delay logic to be removed once WACI ack is available.
  done(delay) {
    let h = this.handler;
    if (delay) {
      wait(delay).then(() => h.done());
    } else {
      h.done();
    }
  }

  cancel() {
    return this.handler.cancel();
  }
}
