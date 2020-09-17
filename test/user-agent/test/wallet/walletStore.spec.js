/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import Vue from 'vue'
import {expect} from 'chai'
import {shallowMount} from '@vue/test-utils'
import Store from '../../../../cmd/user-agent/src/pages/chapi/Store.vue'
import {loadFrameworks, localVue, mockStore, promiseWhen, wcredHandler} from '../common.js'
import * as polyfill from 'credential-handler-polyfill'
import * as trustblocAgent from "@trustbloc/trustbloc-agent"
import {studentCardAndDegreeToStore, studentCardToStore} from './testdata.js'

var uuid = require('uuid/v4')

const storeCredentialFriendlyName = `StudentCard_Mr.Foo_${uuid()}`

function mountStore(wch, done) {
    return function (frameworks) {
        toBeDestroyed.push(frameworks.aries)
        done(shallowMount(Store, {
            localVue,
            store: mockStore(frameworks.aries),
            mocks: {
                $polyfill: polyfill,
                $webCredentialHandler: wch,
                $trustblocAgent: trustblocAgent
            }
        }))
    }
}

let toBeDestroyed = []
after(async () => {
    toBeDestroyed.forEach((obj) => obj.destroy())
})

describe('store a credential in wallet', () => {
    // create web credential handler
    let wch = new wcredHandler()
    // add a credential event
    let credResponse = wch.addEventToQueue({
        type: "credentialstore",
        credentialRequestOrigin: "https://issuer.example.dev",
        credential: {type: 'web', dataType: "VerifiablePresentation", data: studentCardToStore},
    })

    // wait for aries to load to mount component
    let wrapper
    before(async function () {
        return loadFrameworks({loadTrustbloc: true}).then(mountStore(wch, wr => wrapper = wr)
        ).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })

    });

    it('store credential wizard is loaded in wallet', async () => {
        await promiseWhen(() => !wrapper.vm.sendButton)
    })

    it('all credential store metadata are pre-populated in wallet', async () => {
        expect(wrapper.vm.subject).to.equal("StudentCard")
        expect(wrapper.vm.issuer).to.equal("did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw")
        expect(wrapper.vm.issuance).to.deep.equal(new Date("2020-05-27T20:36:05.301Z"))
        expect(wrapper.vm.friendlyName).to.include(wrapper.vm.subject)
        expect(wrapper.vm.credData).to.equal(studentCardToStore)
    })

    it('friendly name is mandatory while storing credential', async () => {
        wrapper.setData({friendlyName: ""})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        expect(wrapper.vm.errors.length).to.equal(1)
        expect(wrapper.vm.errors).to.include('friendly name required.')
    })

    it('stored credential in wallet successfully', async () => {
        wrapper.setData({friendlyName: storeCredentialFriendlyName})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        expect(resp.dataType).to.equal("Response")
        expect(resp.data).to.equal("success")
    })

})

describe('store a credential in wallet with existing friendly name', () => {
    // create web credential handler
    let wch = new wcredHandler()
    // add a credential event
    let credResponse = wch.addEventToQueue({
        type: "credentialstore",
        credentialRequestOrigin: "https://issuer.example.dev",
        credential: {type: 'web', dataType: "VerifiablePresentation", data: studentCardToStore},
    })

    // wait for aries to load to mount component
    let wrapper
    before(function () {
        return loadFrameworks({loadTrustbloc: true}).then(mountStore(wch, (wr) => {
            wrapper = wr
        })).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    });

    it('store credential wizard is loaded in wallet', async () => {
        await promiseWhen(() => !wrapper.vm.sendButton)
    })

    it('all credential store metadata are pre-populated in wallet', async () => {
        expect(wrapper.vm.subject).to.equal("StudentCard")
        expect(wrapper.vm.issuer).to.equal("did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw")
        expect(wrapper.vm.issuance).to.deep.equal(new Date("2020-05-27T20:36:05.301Z"))
        expect(wrapper.vm.friendlyName).to.include(wrapper.vm.subject)
        expect(wrapper.vm.credData).to.equal(studentCardToStore)
    })

    it('stored credential expects "credential name already exists" error', async () => {
        wrapper.setData({friendlyName: storeCredentialFriendlyName})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        expect(resp.dataType).to.equal("Response")
        expect(resp.data).to.include("credential name already exists")
    })

})

describe('store multiple credentials in wallet', () => {
    // create web credential handler
    let wch = new wcredHandler()
    // add a credential event
    let credResponse = wch.addEventToQueue({
        type: "credentialstore",
        credentialRequestOrigin: "https://issuer.example.dev",
        credential: {type: 'web', dataType: "VerifiablePresentation", data: studentCardAndDegreeToStore},
    })

    // wait for aries to load to mount component
    let wrapper
    before(function () {
        return loadFrameworks({loadTrustbloc: true}).then(mountStore(wch, wr => wrapper = wr)
        ).catch(err => {
            console.error('error starting aries framework : errMsg=', err)
        })
    });

    it('store credential wizard is loaded in wallet', async () => {
        await promiseWhen(() => !wrapper.vm.sendButton)
    })

    it('all combined credential store metadata are pre-populated in wallet', async () => {
        expect(wrapper.vm.subject).to.equal("StudentCard,UniversityDegreeCredential")
        expect(wrapper.vm.issuer).to.equal("did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw,did:trustbloc:testnet.trustbloc.dev:EiC_G_44Xq0hj_JmxLScbtMBjOouSgBNI_HuqPm40-t_Uw")
        expect(wrapper.vm.issuance).to.deep.equal(new Date("2020-05-28T21:16:57.780923246Z"))
        expect(wrapper.vm.friendlyName).to.include(wrapper.vm.subject)
        expect(wrapper.vm.credData).to.equal(studentCardAndDegreeToStore)
    })

    it('friendly name is mandatory while storing credential', async () => {
        wrapper.setData({friendlyName: ""})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        expect(wrapper.vm.errors.length).to.equal(1)
        expect(wrapper.vm.errors).to.include('friendly name required.')
    })

    it('stored credentials in wallet successfully', async () => {
        wrapper.setData({friendlyName: `StudentCard_Mr.Foo_${uuid()}`})
        wrapper.find("#storeVCBtn").trigger('click')
        await Vue.nextTick()

        const resp = await credResponse
        expect(resp.dataType).to.equal("Response")
        expect(resp.data).to.equal("success")
    })

})
