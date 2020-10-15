/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/


/**
 * Messenger provides messaging service for sending/replying secured messages our didcomm
 * @param agent instance
 * @class
 */
export class Messenger {
    constructor(agent) {
        this.agent = agent
    }

    /**
     * registers message service with given name, purpose & type
     */
    async register(name, purpose, type) {
        await this.agent.messaging.registerService({name, purpose, type})
    }

    /**
     * services returns list of all message services regstered
     */
    async services() {
        return await this.agent.messaging.services()
    }

    /**
     * send sends given message over connection
     *
     * @param {string} connectionID : connection over which message to be sent
     * @param {object} msg : reply message to be sent
     */
    send(connectionID, msg) {
        this.agent.messaging.send({"connection_ID": connectionID, "message_body": msg})
    }

    /**
     * reply sends given message as a reply to given message ID
     *
     * @param {string} msgID : message ID to which reply to be sent
     * @param {object} msg : reply message to be sent
     */
    reply(msgID, msg) {
        this.agent.messaging.reply({"message_ID": msgID, "message_body": msg})
    }

    /**
     * sendAndWaitForReply sends message to given connection and waits for reply for
     * given topic under same thread
     *
     * @param {string} connectionID : connection over which message to be sent
     * @param {object} msg : message to be sent
     * @param {string} replyTopic (optional) : topic on reply is expected. Will be 'all' if none passed.
     * @param {int} replyTimeout (optional) : time in millisecond to wait for reply. Will be '15000ms' if none passed.
     * @returns {Promise<Object>} containing replied message payload.
     */
    async sendAndWaitForReply(connectionID, msg, replyTopic, replyTimeout) {
        this.send(connectionID, msg)

        const msgID = msg['@id']
        const topic = replyTopic ? replyTopic : 'all'
        const timeout = replyTimeout ? replyTimeout : 15000

        const incomingMsg = await new Promise((resolve, reject) => {
            setTimeout(() => reject(new Error("time out waiting for reply")), timeout)
            const stop = this.agent.startNotifier(msg => {
                let thID = msg.payload.message['~thread'] ? msg.payload.message['~thread'].thid : ''
                if (thID != msgID) {
                    return
                }

                stop()
                resolve(msg.payload.message)
            }, [topic])
        })

        return incomingMsg
    }
}


