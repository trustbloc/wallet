/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

import {DIDAuth} from "./didAuth";
import {getCredentialType} from "..";

const responseType = "VerifiablePresentation"

/**
 * WalletGetByQuery provides CHAPI get vp features
 * @param aries instance & credential event
 * @class
 */
export class WalletGetByQuery extends DIDAuth {
    constructor(aries, credEvent) {
        super(aries, credEvent);
        this.setRequestedContext()
    }

    setRequestedContext() {
        const verifiable = this.credEvent.credentialRequestOptions.web.VerifiablePresentation

        let query = {}
        if (verifiable.query) {
            // supporting only one query for now
            query = Array.isArray(verifiable.query) ? verifiable.query[0] : verifiable.query;
        }

        if (!query.presentationDefinitionQuery) {
            throw 'query type not supported';
        }

        let q=query.presentationDefinitionQuery[0]

        let scopes = [];
        let context=[];


        q.submission_requirements.forEach(function(entry){
            if (!entry.rule.type=="all"){
                throw 'rule type not supported';
            }
            entry.rule.from.forEach(element => scopes.push(element));
        })


        q.input_descriptors.forEach(function(entry){
            let addURI=false
            if (scopes.length==0){
                addURI=true
            }else{
                entry.group.forEach(function(group){
                    scopes.forEach(function(scope) {
                        if (scope==group){
                            addURI=true
                        }
                    });
                });
            }

            if (addURI) {
                context.push(entry.schema.uri)
            }
        })

         this.requestedContext=context

    }

    async getCredentialRecords(did) {
        let vcs = []

        await this.aries.verifiable.getCredentials().then(
            resp => {
                if (resp.result) {
                    resp.result.forEach((item, id) => {
                    if (item.subjectId == did) {
                        this.requestedContext.forEach(function(value, index, object){
                            if (item.context.includes(value)) {
                                vcs.push({
                                    id: id,
                                    name: item.name,
                                    key: item.id,
                                    type: getCredentialType(item.type),
                                    holder: item.subjectId,
                                })
                                object.splice(index, 1);
                            }
                        })
                    }
                    })

                    this.requestedContext.forEach(function(value){
                        // TODO Create Consent Credential for all credentials that not found locally
                        console.log(value)
                    })
                }
            })
            .catch(err => {
                throw err
            })

        return vcs
    }


    async createAndSendPresentation(did, selections) {
        try {
            let vcs = []
            for (let selectedVC of selections) {
                const resp = await this.aries.verifiable.getCredential({
                    id: selectedVC.key
                })
                vcs.push(JSON.parse(resp.verifiableCredential))
            }

            let didMetadata = await this.didStore.getDIDMetadata(did)

            let data
            await this.aries.verifiable.generatePresentation({
                verifiableCredential: vcs,
                did: did,
                domain: this.domain,
                challenge: this.challenge,
                skipVerify: true,
                signatureType: didMetadata.signatureType,
                verificationMethod: didMetadata.keyID
            }).then(resp => {
                    data = resp.verifiablePresentation
                }
            )

            this.sendResponse(responseType, data)
        } catch (e) {
            this.sendResponse("error", e)
        }

    }
}
