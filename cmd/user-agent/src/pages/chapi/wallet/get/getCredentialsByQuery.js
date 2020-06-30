/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {WalletGet} from "./getCredentials";
import jp from 'jsonpath';
import {PresentationExchange} from '../common/presentationExchange'
import {DIDExchange} from '../common/didExchange'

const responseType = "VerifiablePresentation"

/**
 * WalletGetByQuery provides CHAPI get vp features
 * @param aries instance & credential event
 * @class
 */
export class WalletGetByQuery extends WalletGet {
    constructor(aries, credEvent) {
        super(aries, credEvent);

        // validate query and init Presentation Exchange
        let query = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="PresentationDefinitionQuery")]');
        if (query.length == 0) {
            throw "invalid request, incorrect query type"
        }

        this.exchange = new PresentationExchange(query[0].presentationDefinitionQuery)

        this.invitation = jp.query(credEvent, '$..credentialRequestOptions.web.VerifiablePresentation.query[?(@.type=="DIDConnect")].invitation');
    }

    requirementDetails() {
        return this.exchange.requirementDetails()
    }

    async getPresentationSubmission() {
        let credentials = await super.getCredentialRecords()

        let vcs = []
        for (let credential of credentials) {
            const resp = await this.aries.verifiable.getCredential({
                id: credential.key
            })
            vcs.push(JSON.parse(resp.verifiableCredential))
        }

        let submission =  this.exchange.createPresentationSubmission(vcs, this.invitation.length > 0)
        if (Object.keys(submission.presentation_location.descriptor_map).length > 0) {
            let exchange = new DIDExchange(this.aries)
            let connection = await exchange.connect(this.invitation[0])
            console.log(`got connection with ${connection.result.State} status`, connection)
        }

        return submission
    }

    async createAndSendPresentation(walletUser, presentationSubmission, selectedIndexes) {

        if (selectedIndexes && selectedIndexes.length > 0) {
            presentationSubmission = retainOnlySelected(presentationSubmission, selectedIndexes)
        }

        try {
            let data
            await this.aries.verifiable.generatePresentation({
                presentation: presentationSubmission,
                did: walletUser.did,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                signatureType: walletUser.signatureType
            }).then(resp => {
                    data = resp.verifiablePresentation
                }
            )

            this.sendResponse(responseType, data)
        } catch (e) {
            this.sendResponse("error", e)
        }

    }

    cancel() {
        this.sendResponse("error", "wallet declined credential share")
    }

    sendNoCredntials() {
        this.sendResponse("error", "no credentials found for given presentation exchange request")
    }
}

// retainOnlySelected retain only selected VCs and their respective descriptors
function retainOnlySelected(presentationSubmission, selectedIndexes){
    let descriptors = []
    let vcs = []

    let vcCount = 0
    selectedIndexes.forEach(function (selected, index) {
        presentationSubmission.verifiableCredential.forEach(function (vc, vcIndex) {
            if (selected && index == vcIndex) {
                vcs.push(vc)

                let vcDescrs = jp.query(presentationSubmission, `$.presentation_submission.descriptor_map[?(@.path=="$.verifiableCredential.[${vcIndex}]")].id`)
                vcDescrs.forEach(function (id) {
                    descriptors.push({id, path:`$.verifiableCredential.[${vcCount}]`})
                })

                vcCount++
            }
        })
    })

    presentationSubmission.verifiableCredential = vcs
    presentationSubmission.presentation_submission.descriptor_map = descriptors

    return presentationSubmission
}
