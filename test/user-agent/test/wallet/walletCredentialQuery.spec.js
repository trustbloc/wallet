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
import {RegisterWallet} from '../../../../cmd/user-agent/src/pages/chapi/wallet'
import {loadAries, trustBlocStartupOpts, wcredHandler, promiseWhen} from '../common.js'
import * as polyfill from 'credential-handler-polyfill'
import * as trustblocAgent from "@trustbloc/trustbloc-agent"
import {prcAndUdcVP, presentationDefQuery1} from './testdata.js'


const walletUser = "sample-user"

function mountStore(wch, done) {
    return function (aries) {
        toBeDestroyed.push(aries)
        done(shallowMount(Store, {
            mocks: {
                $polyfill: polyfill,
                $webCredentialHandler: wch,
                $arieslib: aries
            }
        }))
    }
}

function mountGet(wch, done) {
    return function (aries) {
        toBeDestroyed.push(aries)
        done(mount(Get, {
            mocks: {
                $polyfill: polyfill,
                $webCredentialHandler: wch,
                $arieslib: aries
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
        let aries = await loadAries()
        let register = new RegisterWallet(polyfill, wch, aries, trustblocAgent, trustBlocStartupOpts)
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
        return loadAries().then(mountStore(wch, wr => wrapper = wr)
        ).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    });

    it('stored permanent resident card and university degree certificate in wallet successfully', async () => {
        wrapper.setData({friendlyName: 'Foo'})
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
        return loadAries().then(mountGet(wch, (wr) => {
            wrapper = wr
        })).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    });


    it('launched get credentials by query and found VCs passing exchange query', async () => {
        let presDef =  wrapper.findComponent(PresentationDefQuery)
        await promiseWhen(() => presDef.vm.vcsFound.length > 0, 10000)

        expect(presDef.vm.vcsFound).to.have.lengthOf(2)
        expect(presDef.vm.loading).to.be.false

        let btn = presDef.find("#share-credentials")
        expect(btn.attributes('disabled')).to.be.equal('true')
    })

    it('shared VCs to create presentation submission', async () => {
        let presDef =  wrapper.findComponent(PresentationDefQuery)

        presDef.vm.selectedVCs = [true, true]
        let btn = presDef.find("#share-credentials")
        btn.trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        //TODO remove this if condition once problem in CI is fixed
        if (resp.dataType == 'VerifiablePresentation') {
            // expect(resp.dataType).to.be.equal('VerifiablePresentation')
            expect(resp.data.presentation_submission).to.deep.equal({
                "descriptor_map": [
                    {
                        "id": "degree_input_1",
                        "path": "$.verifiableCredential.[0]"
                    },
                    {
                        "id": "citizenship_input_1",
                        "path": "$.verifiableCredential.[1]"
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
