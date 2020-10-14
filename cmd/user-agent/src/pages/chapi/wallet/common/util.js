/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export const toLower = text => text.toString().toLowerCase()

export const getCredentialType = (types) => {
    let result = types.filter(type => !isCredentialType(type))
    return result.length > 0 ? result[0] : ""
}

export const isCredentialType = (type) => isVCType(type) || isVPType(type)

export const isVCType = (type) => toLower(type) == 'verifiablecredential'

export const isVPType = (type) => toLower(type) == 'verifiablepresentation'

export function getCredentialMetadata(data, dataType) {
    if (isVCType(dataType)) {
        return getVCMetadata(data)
    }

    if (!data.verifiableCredential) {
        return
    }

    const allCreds = Array.isArray(data.verifiableCredential) ? data.verifiableCredential
        : [data.verifiableCredential];

    let res = {issuer: "", subject: ""}
    let allSubjects = []
    allCreds.forEach((vc) => {
        const {issuance, issuer, subject} = getVCMetadata(vc)

        if (!res.issuance) {
            res.issuance = issuance
        }

        if (issuer && issuer.length > 0) {
            res.issuer = res.issuer == "" ? issuer : `${res.issuer},${issuer}`
        }

        if (subject && subject.length > 0 && !allSubjects.includes(subject)) {
            res.subject = res.subject == "" ? subject : `${res.subject},${subject}`
            allSubjects.push(subject)
        }
    })

    return res
}

export function getDomainAndChallenge(credEvent) {
    if (!credEvent.credentialRequestOptions.web.VerifiablePresentation) {
        return {}
    }

    const verifiable = credEvent.credentialRequestOptions.web.VerifiablePresentation

    let {challenge, domain, query} = verifiable;

    if (query && query.challenge) {
        challenge = query.challenge;
    }

    if (query && query.domain) {
        domain = query.domain;
    }

    if (!domain && credEvent.credentialRequestOrigin) {
        domain = credEvent.credentialRequestOrigin.split('//').pop()
    }

    return {domain, challenge}
}

export const searchByTypeAndHolder = (items, term, key) => {
    if (key) {
        items = items.filter(item => item.holder == key)
    }

    if (term) {
        return items.filter(item => toLower(item.type).includes(toLower(term)))
    }

    return items
}


function getVCMetadata(vc) {
    // issuance date, issuer & subject
    let issuance = (vc.issuanceDate) ? new Date(vc.issuanceDate) : new Date()
    let issuer = (vc.issuer && vc.issuer.id) ? vc.issuer.id : vc.issuer
    let subject = (vc.type && Array.isArray(vc.type)) ? getCredentialType(vc.type) : ''

    return {issuance: issuance, issuer: issuer, subject: subject}
}

export async function waitForNotification(agent, topics, eventType, callback, timeout) {
    return new Promise((resolve, reject) => {
        const stop = agent.startNotifier(notice => {
            if (eventType && notice.payload.Type !== eventType) {
                return
            }

            stop()

            if (callback) {
                try {
                    callback().then(() => {
                        resolve()
                    })
                } catch (err) {
                    reject(err)
                }
            } else {
                resolve(notice.payload)
            }

        }, topics)

        setTimeout(() => {
            stop()
            reject(new Error("timeout while waiting for notification"))
        }, timeout ? timeout : 10000)
    })
}

export const filterCredentialsByType = (creds, types, include) =>
    creds.filter(c => include ? types.includes(getCredentialType(c.type)) : !types.includes(getCredentialType(c.type)))
