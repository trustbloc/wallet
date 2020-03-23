/*
Copyright SecureKey Technologies Inc. All Rights Reserved.

SPDX-License-Identifier: Apache-2.0
*/

'use strict'

const {loadWorker} = require("worker_loader")

// registers messages in pending and posts them to the worker
async function invoke(w, pending, pkg, fn, arg, msgTimeout) {
    return new Promise((resolve, reject) => {
        const timer = setTimeout(_ => reject(new Error(msgTimeout)), 10000)
        let payload = arg
        if (typeof arg === "string") {
            payload = JSON.parse(arg)
        }
        const msg = newMsg(pkg, fn, payload)
        pending.set(msg.id, result => {
            clearTimeout(timer)
            if (result.isErr) {
                reject(new Error(result.errMsg))
            }
            resolve(result.payload)
        })
        w.postMessage(msg)
    })
}

function newMsg(pkg, fn, payload) {
    return {
        // TODO there are several approaches to generate random strings:
        // - which should we implement? do we need cryptographic-grade randomness for this?
        // - alternatively, should the generator be provided by the client?
        id: Math.random().toString(36).slice(2),
        pkg: pkg,
        fn: fn,
        payload: payload
    }
}

/**
 * UserAgent framework class provides UserAgent SSI-agent features.
 *
 * `opts` is an object with the framework's initialization options:
 *
 * {
 *      assetsPath: "/path/serving/the/framework/assets",
 *      blocDomain: "domain"
 * }
 *
 * @param opts framework initialization options.
 * @class
 */
export const Framework = class {
    constructor(opts) {
        return UserAgent(opts)
    }
};


/**
 * UserAgent provides UserAgent SSI-agent functions.
 * @param opts initialization options.
 * @constructor
 */
const UserAgent = function (opts) {
    if (!opts) {
        throw new Error("useragent: missing options")
    }

    if (!opts.assetsPath) {
        throw new Error("useragent: missing assets path")
    }

    // TODO synchronized access
    const notifications = new Map()
    const pending = new Map()

    const instance = {
        /**
         * Test methods.
         * TODO - remove. Used for testing.
         * @type {{_echo: (function(*=): Promise<String>)}}
         * @private
         */
        _test: {
            /**
             * Returns the input text prepended with "echo: ".
             * TODO - remove.
             * @param text
             * @returns {Promise<Object>}
             * @private
             */
            _echo: async function (text) {
                return new Promise((resolve, reject) => {
                    invoke(aw, pending, "test", "echo", {"echo": text}, "_echo() timed out").then(
                        resp => resolve(resp.echo),
                        err => reject(new Error("user agent: _echo() failed. error: " + err.message))
                    )
                })
            }

        },
        destroy: async function () {
            var response = await invoke(aw, pending, "useragent", "Stop", "{}", "timeout while stopping user agent")
            aw.terminate()
            aw = null
            return response
        },
        /**
         * DIDClient methods
         */
        didclient: {
            pkgname: "didclient",

            /**
             * Creates a DID.
             *
             * @param req - json document
             * @returns {Promise<Object>}
             */
            createDID: async function (req) {
                return invoke(aw, pending, this.pkgname, "CreateDID", req, "timeout while creating invitation")
            },
        },
    }

    // start useragent worker
    var aw = loadWorker(
        pending,
        notifications,
        {
            dir: opts.assetsPath,
            wasm: opts.assetsPath + "/user-agent-js-worker.wasm",
            wasmJS: opts.assetsPath + "/wasm_exec.js"
        }
    )


    // return promise which waits for worker to load and user agent to start.
    return new Promise((resolve, reject) => {
        const timer = setTimeout(_ => reject(new Error("timout waiting for user agent to initialize")), 10000)
        notifications.set("asset-ready", async (result) => {
            clearTimeout(timer)
            invoke(aw, pending, "useragent", "Start", opts, "timeout while starting user agent").then(
                resp => resolve(instance),
                err => reject(new Error(err.message))
            )
        })
    })
}
