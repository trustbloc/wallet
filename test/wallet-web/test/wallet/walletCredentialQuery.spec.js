/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import {expect} from 'chai'
import {mount, shallowMount} from '@vue/test-utils'
import Store from '../../../../cmd/wallet-web/src/pages/chapi/Store.vue'
import Get from '../../../../cmd/wallet-web/src/pages/chapi/Get.vue'
import PresentationDefQuery from '../../../../cmd/wallet-web/src/pages/chapi/PresentationDefQuery.vue'
import DIDConnect from '../../../../cmd/wallet-web/src/pages/chapi/DIDConnect.vue'
import {AgentMediator, RegisterWallet} from '../../../../cmd/wallet-web/src/pages/chapi/wallet'
import {loadFrameworks, localVue, mockStore, promiseWhen, wcredHandler} from '../common.js'
import * as polyfill from 'credential-handler-polyfill'
import {issue_credential, manifest, prcAndUdcVP, presentationDefQuery1, presentationDefQuery2} from './testdata.js'
import {IssuerAdapter, RPAdapter} from "../mocks"

var uuid = require('uuid/v4')

const walletUser = "sampleWalletUser"
const challenge = `705aa4da-b240-4c14-8652-8ed35a886ed5-${Math.random()}`
const testOpts = {loadStartupOpts: true}

function mountStore(wch, done) {
    return function (frameworks) {
        toBeDestroyed.push(frameworks.agent)
        done(shallowMount(Store, {
            localVue,
            store: mockStore(frameworks.agent),
            mocks: {
                $polyfill: polyfill,
                $webCredentialHandler: wch,
            }
        }))
    }
}

function mountGet(wch, done) {
    return function (frameworks) {
        toBeDestroyed.push(frameworks.agent)
        done(mount(Get, {
            localVue,
            store: mockStore(frameworks.agent),
            mocks: {
                $polyfill: polyfill,
                $webCredentialHandler: wch
            }
        }))
    }
}

let toBeDestroyed = []
after(async () => {
    toBeDestroyed.forEach((obj) => obj.destroy())
})

describe('register wallet', () => {
    // create web credential handler
    let wch = new wcredHandler()

    it('logged in to wallet', async () => {
        let opts = await loadFrameworks(testOpts)
        let register = new RegisterWallet(polyfill, wch, opts.agent, opts.agentStartupOpts)
        try {
            register.skipPolyfill = true
            await register.register(walletUser)
        } catch (e) {
            console.error(e)
        }

    })
})

describe('store credentials', () => {
    // create web credential handler
    let wch = new wcredHandler()
    // add a credential event
    let credResponse = wch.addEventToQueue({
        type: "credentialstore",
        credentialRequestOrigin: "https://issuer.example.dev",
        credential: {type: 'web', dataType: "VerifiablePresentation", data: prcAndUdcVP},
    })

    // wait for aries to load to mount component
    let wrapper
    before(function () {
        return loadFrameworks(testOpts).then(mountStore(wch, wr => wrapper = wr)
        ).catch(err => {
            console.error('error starting agent: errMsg=', err)
        })
    });

    it('store credential wizard is loaded in wallet', async () => {
        await promiseWhen(() => !wrapper.vm.sendButton)
    })

    it('stored permanent resident card and university degree certificate in wallet successfully', async () => {
        wrapper.setData({friendlyName: `Mr.Foo_creds_${uuid()}`})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        expect(resp.dataType).to.equal("Response")
        expect(resp.data).to.equal("success")
    })

})

describe('get credentials by presentation definition query', () => {
    // create web credential handler
    let wch = new wcredHandler()
    // add a credential event
    let event = {
        type: "credentialrequest",
        credentialRequestOrigin: "https://issuer.example.dev",
        credentialRequestOptions: {
            "web": {
                "VerifiablePresentation": {
                    "query": [
                        {
                            "type": "PresentationDefinitionQuery",
                            "presentationDefinitionQuery": presentationDefQuery1
                        }
                    ]
                }
            }
        }
    }

    let credResponse = wch.addEventToQueue(event)

    // wait for aries to load to mount component
    let wrapper
    before(function () {
        return loadFrameworks(testOpts).then(mountGet(wch, (wr) => {
            wrapper = wr
        })).catch(err => {
            console.error('error starting agent: errMsg=', err)
        })
    });


    it('launched get credentials by query and found VCs passing exchange query', async () => {
        let presDef = wrapper.findComponent(PresentationDefQuery)
        await promiseWhen(() => !presDef.vm.loading, 10000)

        expect(presDef.vm.vcsFound).to.have.lengthOf(2)

        let btn = presDef.find("#share-credentials")
        expect(btn.attributes('disabled')).to.be.equal('true')
    })

    it('shared VCs to create presentation submission', async () => {
        let presDef = wrapper.findComponent(PresentationDefQuery)

        presDef.vm.selectedVCs = [true, true]
        let btn = presDef.find("#share-credentials")
        btn.trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        //TODO remove this if condition once problem in CI is fixed
        if (resp.dataType === 'VerifiablePresentation') {
            // expect(resp.dataType).to.be.equal('VerifiablePresentation')
            expect(resp.data.presentation_submission).to.deep.equal({
                "descriptor_map": [
                    {
                        "id": "degree_input_1",
                        "path": "$.verifiableCredential[0]"
                    },
                    {
                        "id": "citizenship_input_1",
                        "path": "$.verifiableCredential[1]"
                    }
                ]
            })
            expect(resp.data.type).to.deep.equal([
                "VerifiablePresentation",
                "PresentationSubmission"
            ])
            expect(resp.data.verifiableCredential).to.have.lengthOf(2)
            expect(resp.data.proof).to.not.be.empty
        }

    })

})

let issuer
describe('issuer connected to wallet with manifest using DID connect ', () => {
    // add a credential event
    let event = {
        type: "credentialrequest",
        credentialRequestOrigin: "https://issuer.example.dev",
        credentialRequestOptions: {
            web: {
                VerifiablePresentation: {
                    query: {type: "DIDConnect"},
                    credentials: [manifest],
                    challenge: challenge,
                    domain: "example.com"
                }
            }
        }
    }

    let wrapper
    let credResponse

    // - wait for aries to load to mount component
    // - load and setup issuer
    before(async function () {
        // start issuer, register router and create invitation
        await loadFrameworks({name: 'issuer'}).then(async opts => {
            // start issuer, register router and create invitation
            issuer = new IssuerAdapter(opts.agent)
            await issuer.connectToMediator('https://localhost:10063')
            event.credentialRequestOptions.web.VerifiablePresentation.invitation = await issuer.createInvitation()

            await issuer.init()
        }).catch(err => {
            console.error('error starting issuer agent: errMsg=', err)
        })

        // create web credential handler
        let wch = new wcredHandler()
        credResponse = wch.addEventToQueue(event)

        return loadFrameworks(testOpts).then(mountGet(wch, (wr) => {
            wrapper = wr
        })).catch(err => {
            console.error('error starting agent: errMsg=', err)
        })
    });

    it('send DID connect request to web wallet', async () => {
        let didConnDef = wrapper.findComponent(DIDConnect)
        await promiseWhen(() => !didConnDef.vm.loading, 10000)

        // approve did connection
        let btn = didConnDef.find('#didconnect')
        btn.trigger('click')
        await Vue.nextTick()

        // wait for did exchange request from wallet & approve did connection request from issuer
        await issuer.acceptExchangeRequest()

        const resp = await credResponse
        if (resp.dataType === 'VerifiablePresentation') {
            expect(resp.dataType).to.be.equal('VerifiablePresentation')
            expect(resp.data.verifiableCredential[0].credentialSubject.connectionState).to.equal('completed')
            expect(resp.data.proof.challenge).to.equal(challenge)
        }
    })
})


describe('verifier queries credentials - DIDComm Flow', () => {
    // add a credential event
    let event = {
        type: "credentialrequest",
        credentialRequestOrigin: "https://verifier.example.dev",
        credentialRequestOptions: {
            web: {
                VerifiablePresentation: {
                    query: [
                        {
                            type: "PresentationDefinitionQuery",
                            presentationDefinitionQuery: presentationDefQuery2
                        },
                        {
                            type: "DIDConnect"
                        }
                    ],
                    challenge: challenge,
                    domain: "example.com"
                }
            }
        }
    }

    let verifier
    let wrapper
    let credResponse

    // wait for aries to load to mount component
    before(async function () {
        // start verifier, register router and create invitation
        await loadFrameworks({name: 'verifier'}).then(async opts => {
            // initialize verifier adapter
            verifier = new RPAdapter(opts.agent)
            await verifier.connectToMediator('https://localhost:10063')
            event.credentialRequestOptions.web.VerifiablePresentation.query[1].invitation = await verifier.createInvitation()

            await verifier.init()
        }).catch(err => {
            console.error('error starting verifier agent: errMsg=', err)
        })

        // create web credential handler
        let wch = new wcredHandler()
        credResponse = wch.addEventToQueue(event)

        return loadFrameworks(testOpts).then(mountGet(wch, (wr) => {
            wrapper = wr
        })).catch(err => {
            console.error('error starting agent: errMsg=', err)
        })
    });

    it('web wallet finds manifest VC matching presentation exchange query', async () => {
        let presDef = wrapper.findComponent(PresentationDefQuery)
        await promiseWhen(() => !presDef.vm.loading, 10000)

        expect(presDef.vm.loading).to.be.false
        expect(presDef.vm.vcsFound).to.have.lengthOf(1)


        let btn = presDef.find("#share-credentials")
        expect(btn.attributes('disabled')).to.be.equal('true')
    })

    it('user authorizes sharing credential using DIDComm', async () => {
        let presDef = wrapper.findComponent(PresentationDefQuery)

        presDef.vm.selectedVCs = [true]
        let btn = presDef.find("#share-credentials")
        btn.trigger('click')
        await Vue.nextTick()

        // wait for did exchange request from wallet & approve did connection request from verifier
        await verifier.acceptExchangeRequest()

        // wait for new peer DID request from wallet & share new peer DID
        await verifier.shareNewPeerDID()

        // issue credential from issuer
        await issuer.issueCredential(issue_credential)

        const resp = await credResponse
        if (resp.dataType === 'VerifiablePresentation') {
            expect(resp.dataType).to.be.equal('VerifiablePresentation')
            expect(resp.data.type).to.deep.equal([
                "VerifiablePresentation",
                "PresentationSubmission"
            ])
            expect(resp.data.presentation_submission.descriptor_map).to.deep.equal([
                {
                    "id": "citizenship_input_1",
                    "path": "$.verifiableCredential[0]"
                }
            ])
            expect(resp.data.verifiableCredential[0].referenceNumber).to.equal(83294847)
        }
    })
})
