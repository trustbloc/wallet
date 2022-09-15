/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import { nextTick } from 'vue';
import { mount } from '@vue/test-utils';
import Get from '@/layouts/GetLayout.vue';
import MultipleQuery from '@/pages/MultipleQueryPage.vue';
import { getTestData, Setup, MockCredentialHandler, promiseWhen } from '../common';
import { expect } from 'chai';
import i18n from '@/plugins/i18n';
import { v4 as uuidv4 } from 'uuid';

const GET_CRED_USER = 'get_cred_user';

var setup = new Setup({ user: GET_CRED_USER });

before(async function () {
  await setup.loadAgent();
  await setup.createProfile();
  await setup.unlockWallet();
  await setup.createPreference();

  let prcVC = getTestData('prc-vc.json');
  let udcVC = getTestData('udc-vc.json');
  let manifest = getTestData('allvcs-cred-manifest.json');

  await setup.saveCredentials([prcVC, udcVC], {
    manifest,
    descriptorMap: [
      {
        id: 'prc_output',
        format: 'ldp_vc',
        path: '$[0]',
      },
      {
        id: 'udc_output',
        format: 'ldp_vc',
        path: '$[1]',
      },
    ],
  });
});

after(function () {
  setup.destroyAgent();
});

// TODO: revisit to fix this test
describe.skip('sharing a credential from wallet - QueryByExample', function () {
  // credential handler
  const credHandler = new MockCredentialHandler();
  let response;

  // mount vue component
  let wrapper;
  before(function () {
    response = credHandler.setRequestEvent({
      type: 'credentialrequest',
      credentialRequestOrigin: 'https://verifier.example.dev',
      credentialRequestOptions: {
        web: {
          VerifiablePresentation: {
            query: [
              {
                type: 'QueryByExample',
                credentialQuery: {
                  reason: 'Please present a credential for JaneDoe.',
                  example: {
                    '@context': [
                      'https://www.w3.org/2018/credentials/v1',
                      'https://www.w3.org/2018/credentials/examples/v1',
                    ],
                    type: ['UniversityDegreeCredential'],
                  },
                },
              },
            ],
            challenge: uuidv4(),
            domain: 'example.com',
          },
        },
      },
    });
    const store = setup.getStateStore();
    wrapper = mount(Get, {
      global: {
        plugins: [store, i18n],
        mocks: {
          webCredentialHandler: credHandler,
          t: () => '',
        },
      },
    });
  });

  it('share credential screen is presented to user', async function () {
    let query = wrapper.findComponent(MultipleQuery);
    await promiseWhen(() => !query.vm.loading, 10000);
  });

  it('found matching result in wallet', function () {
    let query = wrapper.findComponent(MultipleQuery);
    expect(query.vm.processedCredentials).to.have.lengthOf(1);
  });

  it('user shares credential successfully !', async function () {
    let query = wrapper.findComponent(MultipleQuery);

    let btn = query.find('#share-credentials');
    btn.trigger('click');

    await nextTick();

    const result = await response;
    console.log('response from wallet', result);
    expect(result.dataType).to.equal('VerifiablePresentation');
    expect(result.data.proof).to.not.empty;
    expect(result.data.verifiableCredential).to.have.lengthOf(1);
  });
});

// TODO: revisit to fix this test
describe.skip('sharing a credential from wallet - PresentationExchange', function () {
  // credential handler
  const credHandler = new MockCredentialHandler();
  let response;

  // mount vue component
  let wrapper;
  before(function () {
    response = credHandler.setRequestEvent({
      type: 'credentialrequest',
      credentialRequestOrigin: 'https://verifier.example.dev',
      credentialRequestOptions: {
        web: {
          VerifiablePresentation: {
            query: [
              {
                type: 'PresentationExchange',
                credentialQuery: [
                  {
                    id: '22c77155-edf2-4ec5-8d44-b393b4e4fa38',
                    input_descriptors: [
                      {
                        id: '20b073bb-cede-4912-9e9d-334e5702077b',
                        schema: [
                          { uri: 'https://www.w3.org/2018/credentials#VerifiableCredential' },
                        ],
                        constraints: { fields: [{ path: ['$.credentialSubject.familyName'] }] },
                      },
                    ],
                  },
                ],
              },
            ],
            challenge: uuidv4(),
            domain: 'example.com',
          },
        },
      },
    });
    const store = setup.getStateStore();
    wrapper = mount(Get, {
      global: {
        plugins: [store, i18n],
        mocks: {
          webCredentialHandler: credHandler,
          t: () => '',
        },
      },
    });
  });

  it('share credential screen is presented to user', async function () {
    let query = wrapper.findComponent(MultipleQuery);
    await promiseWhen(() => !query.vm.loading, 10000);
  });

  it('found matching result in wallet', function () {
    let query = wrapper.findComponent(MultipleQuery);
    expect(query.vm.processedCredentials).to.have.lengthOf(1);
  });

  it('user shares credential successfully !', async function () {
    let query = wrapper.findComponent(MultipleQuery);

    let btn = query.find('#share-credentials');
    btn.trigger('click');

    await nextTick();

    const result = await response;
    console.log('response from wallet', result);
    expect(result.dataType).to.equal('VerifiablePresentation');
    expect(result.data.proof).to.not.empty;
    expect(result.data.verifiableCredential).to.have.lengthOf(1);
  });
});

// TODO: revisit to fix this test
describe.skip('sharing multiple credentials from wallet - MultiQuery (QueryByExample, QueryByFrame)', function () {
  // credential handler
  const credHandler = new MockCredentialHandler();
  let response;

  // mount vue component
  let wrapper;
  before(async function () {
    response = credHandler.setRequestEvent({
      type: 'credentialrequest',
      credentialRequestOrigin: 'https://verifier.example.dev',
      credentialRequestOptions: {
        web: {
          VerifiablePresentation: {
            query: [
              {
                type: 'QueryByFrame',
                credentialQuery: [
                  {
                    reason: 'Please provide your Passport details.',
                    frame: {
                      '@context': [
                        'https://www.w3.org/2018/credentials/v1',
                        'https://w3id.org/citizenship/v1',
                        'https://w3id.org/security/bbs/v1',
                      ],
                      type: ['VerifiableCredential', 'PermanentResidentCard'],
                      '@explicit': true,
                      identifier: {},
                      issuer: {},
                      issuanceDate: {},
                      credentialSubject: { '@explicit': true, name: {}, spouse: {} },
                    },
                    trustedIssuer: [
                      { issuer: 'did:example:76e12ec712ebc6f1c221ebfeb1f', required: true },
                    ],
                    required: true,
                  },
                ],
              },
              {
                type: 'QueryByExample',
                credentialQuery: [
                  {
                    reason: 'Please present your valid degree certificate.',
                    example: {
                      '@context': [
                        'https://www.w3.org/2018/credentials/v1',
                        'https://www.w3.org/2018/credentials/examples/v1',
                      ],
                      type: ['UniversityDegreeCredential'],
                      trustedIssuer: [
                        { issuer: 'urn:some:required:issuer' },
                        {
                          required: true,
                          issuer: 'did:example:76e12ec712ebc6f1c221ebfeb1f',
                        },
                      ],
                      credentialSubject: { id: 'did:example:ebfeb1f712ebc6f1c276e12ec21' },
                    },
                  },
                ],
              },
            ],
            challenge: uuidv4(),
            domain: 'example.com',
          },
        },
      },
    });
    const udcBBSVC = getTestData('udc-bbs-vc.json');
    const store = setup.getStateStore();
    // prepare manifest
    let manifest = getTestData('allvcs-cred-manifest.json');
    manifest.id = uuidv4();

    await setup.saveCredentials([udcBBSVC], {
      manifest,
      descriptorMap: [
        {
          id: 'udc_output',
          format: 'ldp_vc',
          path: '$[0]',
        },
      ],
    });

    wrapper = mount(Get, {
      global: {
        plugins: [store, i18n],
        mocks: {
          webCredentialHandler: credHandler,
          t: () => '',
        },
      },
    });
  });

  it('share credential screen is presented to user', async function () {
    let query = wrapper.findComponent(MultipleQuery);
    await promiseWhen(() => !query.vm.loading, 10000);
  });

  it('found matching result in wallet', function () {
    let query = wrapper.findComponent(MultipleQuery);
    expect(query.vm.processedCredentials).to.have.lengthOf(3);
  });

  it('user shares credential successfully !', async function () {
    let query = wrapper.findComponent(MultipleQuery);

    let btn = query.find('#share-credentials');
    btn.trigger('click');

    await nextTick();

    const result = await response;
    console.log('response from wallet', result);
    expect(result.dataType).to.equal('VerifiablePresentation');
    expect(result.data.proof).to.not.empty;
    expect(result.data.verifiableCredential).to.have.lengthOf(3);
  });
});
