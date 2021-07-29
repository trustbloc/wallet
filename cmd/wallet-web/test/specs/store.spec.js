/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue';
import { shallowMount } from '@vue/test-utils';
import Store from '../../src/pages/Store.vue';
import {
  testConfig,
  getTestData,
  Setup,
  loadFrameworks,
  MockCredentialHandler,
  localVue,
  promiseWhen,
} from '../common';
import { expect } from 'chai';

const SAVE_CRED_USER = 'save_cred_user';

let setup = new Setup({ user: SAVE_CRED_USER });

before(async function () {
  await setup.loadAgent();
  await setup.createProfile();
  await setup.unlockWallet();
});

after(function () {
  setup.destroyAgent();
});

describe('saving a credential into wallet', function () {
  //test data
  let samplePresentation = getTestData('sample-presentation-1.json');

  // credential handler
  let credHandler = new MockCredentialHandler();
  let response = credHandler.setRequestEvent({
    type: 'credentialstore',
    credentialRequestOrigin: 'https://issuer.example.dev',
    credential: { type: 'web', dataType: 'VerifiablePresentation', data: samplePresentation },
  });

  // mount vue component
  let wrapper;
  before(function () {
    wrapper = shallowMount(Store, {
      localVue,
      store: setup.getStateStore(),
      mocks: {
        $webCredentialHandler: credHandler,
      },
    });
  });

  it('save credential wizard is loaded in wallet', async () => {
    await promiseWhen(() => !wrapper.vm.loading);
  });

  it('credential to be saved are presented to user', async () => {
    expect(wrapper.vm.records).to.have.lengthOf(1);
  });

  it('saved credential into wallet successfully !', async () => {
    wrapper.find('#storeVCBtn').trigger('click');
    await Vue.nextTick();

    const resp = await response;
    expect(resp.dataType).to.equal('response');
    expect(resp.data).to.equal('success');
  });
});

describe('saving multiple credentials into wallet', function () {
  //test data
  let samplePresentation = getTestData('sample-presentation-2.json');

  // credential handler
  let credHandler = new MockCredentialHandler();
  let response = credHandler.setRequestEvent({
    type: 'credentialstore',
    credentialRequestOrigin: 'https://issuer.example.dev',
    credential: { type: 'web', dataType: 'VerifiablePresentation', data: samplePresentation },
  });

  // mount vue component
  let wrapper;
  before(function () {
    wrapper = shallowMount(Store, {
      localVue,
      store: setup.getStateStore(),
      mocks: {
        $webCredentialHandler: credHandler,
      },
    });
  });

  it('save credential wizard is loaded in wallet', async () => {
    await promiseWhen(() => !wrapper.vm.sendButton);
  });

  it('credentials to be saved are presented to user', async () => {
    expect(wrapper.vm.records).to.have.lengthOf(2);
  });

  it('saved credentials in wallet successfully', async () => {
    wrapper.find('#storeVCBtn').trigger('click');
    await Vue.nextTick();

    const resp = await response;
    expect(resp.dataType).to.equal('response');
    expect(resp.data).to.equal('success');
  });
});
