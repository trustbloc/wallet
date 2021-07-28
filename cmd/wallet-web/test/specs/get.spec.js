/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import {mount} from '@vue/test-utils'
import Get from '../../src/pages/Get.vue'
import MultipleQuery from '../../src/pages/MultipleQuery.vue'
import {testConfig, getTestData, Setup, loadFrameworks, MockCredentialHandler, localVue, promiseWhen} from "../common";
import {expect} from "chai";


const GET_CRED_USER = 'get_cred_user'

var uuid = require('uuid/v4')
var setup = new Setup({user: GET_CRED_USER})

before(async function () {
    await setup.loadAgent()
    await setup.createProfile()
    await setup.unlockWallet()
    await setup.createPreference()

    let prcVC = getTestData('prc-vc.json')
    let udcVC = getTestData('udc-vc.json')

    await setup.saveCredentials(prcVC, udcVC)
});

after(function () {
    setup.destroyAgent()
});

describe('sharing a credential from wallet - QueryByExample', function () {
    // credential handler
    let credHandler = new MockCredentialHandler()
    let response = credHandler.setRequestEvent({
        type: "credentialrequest",
        credentialRequestOrigin: "https://verifier.example.dev",
        credentialRequestOptions: {
            web: {
                VerifiablePresentation: {
                    query: [
                        {
                            type: "QueryByExample",
                            credentialQuery: {
                                reason: "Please present a credential for JaneDoe.",
                                example: {
                                    "@context": [
                                        "https://www.w3.org/2018/credentials/v1",
                                        "https://www.w3.org/2018/credentials/examples/v1"
                                    ],
                                    type: ["UniversityDegreeCredential"]
                                }
                            }
                        }
                    ],
                    challenge: uuid(),
                    domain: "example.com"
                }
            }
        }
    })



    // mount vue component
    let wrapper
    before(async function () {
        wrapper = mount(Get, {
            localVue,
            store: setup.getStateStore(),
            mocks: {
                $webCredentialHandler: credHandler
            }
        })
    });

    it('share credential screen is presented to user', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        await promiseWhen(() => !query.vm.loading, 10000)
    })

    it('found matching result in wallet', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        expect(query.vm.records).to.have.lengthOf(1)
    })

    it('user shares credential successfully !', async () => {
        let query = wrapper.findComponent(MultipleQuery)

        let btn = query.find("#share-credentials")
        btn.trigger('click')

        await Vue.nextTick()

        const result = await response
        console.log('response from wallet', result)
        expect(result.dataType).to.equal('VerifiablePresentation')
        expect(result.data.proof).to.not.empty
        expect(result.data.verifiableCredential).to.have.lengthOf(1)
    })
})

describe('sharing a credential from wallet - PresentationExchange', function () {
    // credential handler
    let credHandler = new MockCredentialHandler()
    let response = credHandler.setRequestEvent({
        type: "credentialrequest",
        credentialRequestOrigin: "https://verifier.example.dev",
        credentialRequestOptions: {
            web: {
                VerifiablePresentation: {
                    query: [
                        {
                            "type": "PresentationExchange",
                            "credentialQuery": [{
                                "id": "22c77155-edf2-4ec5-8d44-b393b4e4fa38",
                                "input_descriptors": [{
                                    "id": "20b073bb-cede-4912-9e9d-334e5702077b",
                                    "schema": [{"uri": "https://www.w3.org/2018/credentials/v1#VerifiableCredential"}],
                                    "constraints": {"fields": [{"path": ["$.credentialSubject.familyName"]}]}
                                }]
                            }]
                        }
                    ],
                    challenge: uuid(),
                    domain: "example.com"
                }
            }
        }
    })

    // mount vue component
    let wrapper
    before(async function () {
        wrapper = mount(Get, {
            localVue,
            store: setup.getStateStore(),
            mocks: {
                $webCredentialHandler: credHandler
            }
        })
    });

    it('share credential screen is presented to user', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        await promiseWhen(() => !query.vm.loading, 10000)
    })

    it('found matching result in wallet', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        expect(query.vm.records).to.have.lengthOf(1)
    })

    it('user shares credential successfully !', async () => {
        let query = wrapper.findComponent(MultipleQuery)

        let btn = query.find("#share-credentials")
        btn.trigger('click')

        await Vue.nextTick()

        const result = await response
        console.log('response from wallet', result)
        expect(result.dataType).to.equal('VerifiablePresentation')
        expect(result.data.proof).to.not.empty
        expect(result.data.verifiableCredential).to.have.lengthOf(1)
    })
})

describe('sharing multiple credentials from wallet - MultiQuery (QueryByExample, QueryByFrame)', function () {
    // credential handler
    let credHandler = new MockCredentialHandler()
    let response = credHandler.setRequestEvent({
        type: "credentialrequest",
        credentialRequestOrigin: "https://verifier.example.dev",
        credentialRequestOptions: {
            web: {
                VerifiablePresentation: {
                    query: [
                        {
                            "type": "QueryByFrame",
                            "credentialQuery": [{
                                "reason": "Please provide your Passport details.",
                                "frame": {
                                    "@context": ["https://www.w3.org/2018/credentials/v1", "https://w3id.org/citizenship/v1", "https://w3id.org/security/bbs/v1"],
                                    "type": ["VerifiableCredential", "PermanentResidentCard"],
                                    "@explicit": true,
                                    "identifier": {},
                                    "issuer": {},
                                    "issuanceDate": {},
                                    "credentialSubject": {"@explicit": true, "name": {}, "spouse": {}}
                                },
                                "trustedIssuer": [{"issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f", "required": true}],
                                "required": true
                            }]
                        }, {
                            "type": "QueryByExample",
                            "credentialQuery": [{
                                "reason": "Please present your valid degree certificate.",
                                "example": {
                                    "@context": ["https://www.w3.org/2018/credentials/v1", "https://www.w3.org/2018/credentials/examples/v1"],
                                    "type": ["UniversityDegreeCredential"],
                                    "trustedIssuer": [
                                        {"issuer": "urn:some:required:issuer"},
                                        {
                                            "required": true,
                                            "issuer": "did:example:76e12ec712ebc6f1c221ebfeb1f"
                                        }
                                    ],
                                    "credentialSubject": {"id": "did:example:ebfeb1f712ebc6f1c276e12ec21"}
                                }
                            }]
                        }
                    ],
                    challenge: uuid(),
                    domain: "example.com"
                }
            }
        }
    })

    let udcBBSVC = getTestData('udc-bbs-vc.json')

    // mount vue component
    let wrapper
    before(async function () {
        await setup.saveCredentials(udcBBSVC)
        wrapper = mount(Get, {
            localVue,
            store: setup.getStateStore(),
            mocks: {
                $webCredentialHandler: credHandler
            }
        })
    });

    it('share credential screen is presented to user', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        await promiseWhen(() => !query.vm.loading, 10000)
    })

    it('found matching result in wallet', async () => {
        let query = wrapper.findComponent(MultipleQuery)
        expect(query.vm.records).to.have.lengthOf(3)
    })

    it('user shares credential successfully !', async () => {
        let query = wrapper.findComponent(MultipleQuery)

        let btn = query.find("#share-credentials")
        btn.trigger('click')

        await Vue.nextTick()

        const result = await response
        console.log('response from wallet', result)
        expect(result.dataType).to.equal('VerifiablePresentation')
        expect(result.data.proof).to.not.empty
        expect(result.data.verifiableCredential).to.have.lengthOf(3)
    })
})
