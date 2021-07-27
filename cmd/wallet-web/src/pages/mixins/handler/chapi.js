/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

const responses = {
  DONE: 'response',
  PRESENT: 'VerifiablePresentation',
  CANCEL: 'error',
};

const enabledTypes = ['VerifiablePresentation', 'VerifiableCredential'];

/**
 * CHAPIEventHandler handles credential events.
 *
 * @param credEvent instance
 * @class
 */
export class CHAPIEventHandler {
  constructor(credEvent) {
    this.credEvent = credEvent;
  }

  // returns credential for store event & query for get event.
  getEventData() {
    if (this.credEvent.credential) {
      return this.credEvent.credential;
    }

    let { query, challenge, domain, invitation, credentials } =
      this.credEvent.credentialRequestOptions.web.VerifiablePresentation;
    return { query, challenge, domain, invitation, credentials };
  }

  getRequestor() {
    let { credentialRequestOrigin } = this.credEvent;

    return credentialRequestOrigin;
  }

  // call credential handler call back with basic success response.
  done() {
    this.credEvent.respondWith(
      new Promise(function (resolve) {
        return resolve({ dataType: responses.DONE, data: 'success' });
      })
    );
  }

  // call credential handler call back with verifiable presentation response.
  present(data) {
    this.credEvent.respondWith(
      new Promise(function (resolve) {
        return resolve({ dataType: responses.PRESENT, data });
      })
    );
  }

  // credential handler call back with cancel response.
  //TODO reject promise
  cancel() {
    this.credEvent.respondWith(
      new Promise(function (resolve) {
        return resolve({ dataType: responses.CANCEL, data: 'operation cancelled' });
      })
    );
  }

  sendError(e) {
    this.credEvent.respondWith(
      new Promise(function (resolve, reject) {
        return reject(e);
      })
    );
  }
}

/**
 * CHAPIHandler handles CHAPI life cycles like install, uninstall.
 *
 * @param polyfill & webcredential handler instance
 * @class
 */
export class CHAPIHandler {
  constructor(polyfill, webCredHandler, mediator) {
    this.polyfill = polyfill;
    this.webCredHandler = webCredHandler;
    this.mediator = mediator;
  }

  // install credential handler polyfill handlers.
  async install(name) {
    try {
      await this.polyfill.loadOnce(this.mediator);
    } catch (e) {
      console.error('Error in loadOnce:', e);
      throw 'failed to register wallet, please try again later';
    }

    const registration = await this.webCredHandler.installHandler({
      url: `${__webpack_public_path__}worker`,
    });

    await registration.credentialManager.hints.set('edge', { name, enabledTypes });
  }

  // uninstall credential handler polyfill handlers.
  async uninstall() {
    try {
      await this.polyfill.loadOnce(this.mediator);
    } catch (e) {
      console.error('Error in loadOnce:', e);
      return;
    }

    await this.webCredHandler.uninstallHandler({ url: `${__webpack_public_path__}worker` });
  }
}
