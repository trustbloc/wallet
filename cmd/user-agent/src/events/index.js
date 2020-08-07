/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

export const PRE_STATE = "pre_state"
export const POST_STATE = "post_state"

const defaultTimeout = 10000
const defaultTimeoutError = "time out while waiting for event"
const defaultTopic = 'all'

export function waitForEvent(options = {}) {
    if (!options.timeout) {
        options.timeout = defaultTimeout
    }

    if (!options.timeoutError) {
        options.timeoutError = defaultTimeoutError
    }

    if (!options.topic) {
        options.topic = defaultTopic
    }

    return new Promise((resolve, reject) => {
        setTimeout(() => reject(new Error(options.timeoutError)), options.timeout)
        const stop = window.$aries.startNotifier(event => {
            try {
                let payload = event.payload;

                if (options.connectionID && payload.Properties &&
                    payload.Properties.connectionID !== options.connectionID) {
                    return
                }

                if (options.stateID && payload.StateID !== options.stateID) {
                    return
                }

                if (options.type && payload.Type !== options.type) {
                    return
                }

                stop()
                resolve(payload)
            } catch (e) {
                stop()
                reject(e)
            }
        }, [options.topic])
    })
}
