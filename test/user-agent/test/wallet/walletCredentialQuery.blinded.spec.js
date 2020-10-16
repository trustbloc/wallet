/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import {expect} from 'chai'
import {mount, shallowMount} from '@vue/test-utils'
import Store from '../../../../cmd/user-agent/src/pages/chapi/Store.vue'
import Get from '../../../../cmd/user-agent/src/pages/chapi/Get.vue'
import PresentationDefQuery from '../../../../cmd/user-agent/src/pages/chapi/PresentationDefQuery.vue'
import DIDConnect from '../../../../cmd/user-agent/src/pages/chapi/DIDConnect.vue'
import {AgentMediator, RegisterWallet} from '../../../../cmd/user-agent/src/pages/chapi/wallet'
import {loadFrameworks, localVue, mockStore, promiseWhen, wcredHandler} from '../common.js'
import * as polyfill from 'credential-handler-polyfill'
import {issue_credential, manifest, prcAndUdcVP, presentationDefQuery1, presentationDefQuery2} from './testdata.js'
import {waitForEvent} from "../../../../cmd/user-agent/src/events";
import {getMediatorConnections} from "../../../../cmd/user-agent/src/pages/chapi/wallet/didcomm/mediator.js"

var uuid = require('uuid/v4')

const walletUser = "sampleWalletUser"
const challenge = `705aa4da-b240-4c14-8652-8ed35a886ed5-${Math.random()}`
const testOpts = {loadStartupOpts: true, blinded: true}

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

describe('register wallet with blinded routing enabled', () => {
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

let issuer
describe('issuer with manifest connected to wallet over blinded routing', () => {
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

    // - wait for aries to load to mount component
    // - load and setup issuer
    before(async function () {
        await loadFrameworks({name: 'issuer-blinded'}).then(async opts => {
            let mediator = new AgentMediator(opts.agent)
            await mediator.connect('https://localhost:10063').then(ur => {
                console.log("issuer mediator registered successfully")
            }).catch(err => {
                console.error('failed to register mediator for issuer agent: errMsg=', err)
            })
            event.credentialRequestOptions.web.VerifiablePresentation.invitation = await mediator.createInvitation()

            issuer = opts.agent

            await issuer.messaging.registerService({
                name: 'request-for-diddoc',
                type: 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-req'
            })
            await issuer.messaging.registerService({
                name: 'share-diddoc-req',
                type: 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-req'
            })

        }).catch(err => {
            console.error('error starting issuer agent: errMsg=', err)
        })
    })

    // create web credential handler
    let wch = new wcredHandler()
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

    it('send DID connect request to web wallet', async () => {
        let didConnDef = wrapper.findComponent(DIDConnect)
        await promiseWhen(() => !didConnDef.vm.loading, 10000)

        // approve did connection
        let btn = didConnDef.find('#didconnect')
        btn.trigger('click')
        await Vue.nextTick()

        // wait for did exchange request from wallet
        let res = await waitForEvent(issuer, {topic: 'didexchange_actions'})

        // approve did connection request from issuer
        await issuer.didexchange.acceptExchangeRequest({
            id: res.Properties.connectionID,
            router_connections: await getMediatorConnections(issuer, true),
        })

        // wait for request from wallet for peer DID
        let reqForIssuerDID = await waitForEvent(issuer, {topic: 'request-for-diddoc'})

        // send any sample peer DID to wallet
        let sampleRes = await issuer.vdri.resolveDID({id: reqForIssuerDID.mydid})
        issuer.messaging.reply({
            "message_ID": reqForIssuerDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-res',
                data: {didDoc: sampleRes.did},
            }
        })

        // wait for did shared by wallet
        let sharedDID = await waitForEvent(issuer, {topic: 'share-diddoc-req'})

        // send acknowledgement to wallet
        issuer.messaging.reply({
            "message_ID": sharedDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-res'
            }
        })

        expect(sharedDID.message.data).to.not.be.null
        expect(sharedDID.message.data.didDoc).to.not.be.null
        expect(sharedDID.message.data.didDoc['@context']).to.deep.equal(["https://w3id.org/did/v1"])

        const resp = await credResponse
        if (resp.dataType === 'VerifiablePresentation') {
            expect(resp.dataType).to.be.equal('VerifiablePresentation')
            expect(resp.data.verifiableCredential[0].credentialSubject.connectionState).to.equal('completed')
            expect(resp.data.proof.challenge).to.equal(challenge)
        }
    })
})


describe('verifier queries credentials - DIDComm Flow using blinded routing', () => {
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
            let mediator = new AgentMediator(opts.agent)

            await mediator.connect('https://localhost:10063').then(ur => {
                console.log("verifier mediator registered successfully")
            }).catch(err => {
                console.error('failed to register mediator for verifier agent: errMsg=', err)
            })

            event.credentialRequestOptions.web.VerifiablePresentation.query[1].invitation = await mediator.createInvitation()
            verifier = opts.agent

            await verifier.messaging.registerService({
                name: 'request-for-diddoc',
                type: 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-req'
            })
            await verifier.messaging.registerService({
                name: 'share-diddoc-req',
                type: 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-req'
            })

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

        // approve did connection request from verifier
        let res = await waitForEvent(verifier, {topic: 'didexchange_actions'})
        await verifier.didexchange.acceptExchangeRequest({
            id: res.Properties.connectionID,
            router_connections: await getMediatorConnections(verifier, true),
        })

        // wait for request from wallet for peer DID
        let reqForRPDID = await waitForEvent(verifier, {topic: 'request-for-diddoc'})

        // send any sample peer DID to wallet
        let sampleRes = await verifier.vdri.resolveDID({id: reqForRPDID.mydid})
        verifier.messaging.reply({
            "message_ID": reqForRPDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/diddoc-res',
                data: {didDoc: sampleRes.did},
            }
        })

        // wait for did shared by wallet
        let sharedDID = await waitForEvent(verifier, {topic: 'share-diddoc-req'})

        // send acknowledgement to wallet
        verifier.messaging.reply({
            "message_ID": sharedDID.message['@id'],
            "message_body": {
                "@id": uuid(),
                "@type": 'https://trustbloc.github.io/blinded-routing/1.0/share-diddoc-res'
            }
        })

        expect(sharedDID.message.data).to.not.be.null
        expect(sharedDID.message.data.didDoc).to.not.be.null
        expect(sharedDID.message.data.didDoc['@context']).to.deep.equal(["https://w3id.org/did/v1"])

        // issue credential from issuer
        res = await waitForEvent(issuer, {topic: 'issue-credential_actions'})
        await issuer.issuecredential.acceptRequest({
            piid: res.Properties.piid,
            issue_credential
        })

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
