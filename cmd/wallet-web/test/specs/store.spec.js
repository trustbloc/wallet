/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { nextTick } from 'vue';
import { shallowMount } from '@vue/test-utils';
import Store from '@/pages/StorePage.vue';
import { getTestData, Setup, MockCredentialHandler, promiseWhen } from '../common';
import { expect } from 'chai';
import i18n from '@/plugins/i18n';

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
  // credential handler
  const credHandler = new MockCredentialHandler();
  let response;

  // mount vue component
  let wrapper;
  before(function () {
    //test data
    const samplePresentation = getTestData('sample-presentation-1.json');
    response = credHandler.setRequestEvent({
      type: 'credentialstore',
      credentialRequestOrigin: 'https://issuer.example.dev',
      credential: { type: 'web', dataType: 'VerifiablePresentation', data: samplePresentation },
    });
    const store = setup.getStateStore();
    wrapper = shallowMount(Store, {
      global: {
        plugins: [store, i18n],
        mocks: {
          webCredentialHandler: credHandler,
          t: () => '',
        },
      },
    });
  });

  it('save credential wizard is loaded in wallet', async function () {
    await promiseWhen(() => !wrapper.vm.loading);
  });

  it('credential to be saved are presented to user', function () {
    expect(wrapper.vm.processedCredentials).to.have.lengthOf(1);
  });

  it('saved credential into wallet successfully!', async function () {
    wrapper.find('button.btn-primary').trigger('click');
    await nextTick();

    const resp = await response;
    expect(resp.dataType).to.equal('response');
    expect(resp.data).to.equal('success');
  });
});

describe('saving multiple credentials into wallet', function () {
  // credential handler
  const credHandler = new MockCredentialHandler();
  let response;

  // mount vue component
  let wrapper;
  before(function () {
    //test data
    let samplePresentation = getTestData('sample-presentation-2.json');
    response = credHandler.setRequestEvent({
      type: 'credentialstore',
      credentialRequestOrigin: 'https://issuer.example.dev',
      credential: { type: 'web', dataType: 'VerifiablePresentation', data: samplePresentation },
    });
    const store = setup.getStateStore();
    wrapper = shallowMount(Store, {
      global: {
        plugins: [store, i18n],
        mocks: {
          webCredentialHandler: credHandler,
          t: () => '',
        },
      },
    });
  });

  it('save credential wizard is loaded in wallet', async function () {
    await promiseWhen(() => !wrapper.vm.sendButton);
  });

  it('credentials to be saved are presented to user', function () {
    expect(wrapper.vm.processedCredentials).to.have.lengthOf(2);
  });

  it('saved credentials in wallet successfully', async function () {
    wrapper.find('button.btn-primary').trigger('click');
    await nextTick();

    const resp = await response;
    expect(resp.dataType).to.equal('response');
    expect(resp.data).to.equal('success');
  });
});
