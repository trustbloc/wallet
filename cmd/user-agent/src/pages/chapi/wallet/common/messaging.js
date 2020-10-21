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
     * registers message service with given name, purpose & type.
     *
     * @param {string} name : name on which message service to ve registered.
     * @param {string[]} purpose : purpose of the message.
     * @param {string} type : message type
     */
    async register(name, purpose, type) {
        await this.agent.messaging.registerService({name, purpose, type})
    }

    /**
     * services returns list of all message services registered.
     */
    async services() {
        return await this.agent.messaging.services()
    }

    /**
     * send sends given message over connection.
     *
     * @param {string} connectionID : connection over which message to be sent.
     * @param {object} msg : reply message to be sent.
     *
     * @param {replyTopic, timeout}
     *
     * replyTopic: (optional) topic on which response for this message is expected : is provided, this function
     * will wait for incoming reply on this topic.
     *
     * timeout: (optional) timeout for reply.
     */
    async send(connectionID, msg, { replyTopic = '', timeout = undefined} = {}) {
        this.agent.messaging.send({"connection_ID": connectionID, "message_body": msg})

        if (replyTopic) {
            return await this.waitForReply(msg['@id'], replyTopic, timeout)
        }
    }

    /**
     * reply sends given message as a reply to given message ID.
     *
     * @param {string} msgID : message ID to which reply to be sent.
     * @param {object} msg : reply message to be sent.
     * @param {replyTopic, startNewThread, timeout}
     *
     * replyTopic : (optional) topic on which response for this message is expected : is provided, this function
     * will wait for incoming reply on this topic.
     *
     * startNewThread: (optional) if provided, then will messenger will start new thread by keeping original
     * thread as a parent thread.
     *
     * timeout: (optional) timeout for reply.
     */
    async reply(msgID, msg, { replyTopic = '', startNewThread = false, timeout = undefined} = {}) {
        this.agent.messaging.reply({"message_ID": msgID, "message_body": msg, "start_new_thread": startNewThread})

        if (replyTopic) {
            return await this.waitForReply(msg['@id'], replyTopic, timeout)
        }
    }

    /**
     * wait for incoming message for given message thread id and topic.
     *
     * @param {string} msgID : message ID to which reply to be sent.
     * @param {string} topic : topic on which reply is expected.
     * @param {int} timeout (optional) : timeout for reply.
     *
     */
    async waitForReply(msgID, topic, timeout) {
        timeout = timeout ? timeout : 30000

        const incomingMsg = await new Promise((resolve, reject) => {
            setTimeout(() => reject(new Error(`time out waiting reply for topic '${topic}'`)), timeout)
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


